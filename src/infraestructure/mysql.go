package infraestructure

import (
	"database/sql"
	domain "notifications/src/domain/entities"
	"notifications/src/domain/repositories"
)

type MySQLNotificationRepository struct {
	db *sql.DB
}

func NewMySQLNotificationRepository(db *sql.DB) repositories.NotificationRepository {
	return &MySQLNotificationRepository{db: db}
}

// Crear una nueva notificación en la base de datos
func (repo *MySQLNotificationRepository) Create(notification *domain.Notification) error {
	query := "INSERT INTO notifications (loan_id, title, status, message, timestamp) VALUES (?, ?, ?, ?, ?)"
	_, err := repo.db.Exec(query, notification.LoanID, notification.Title, notification.Status, notification.Message, notification.Timestamp)
	return err
}

// Obtener todas las notificaciones de la base de datos
func (repo *MySQLNotificationRepository) GetAll() ([]domain.Notification, error) {
	query := "SELECT loan_id, title, status, message, timestamp FROM notifications ORDER BY timestamp DESC"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []domain.Notification
	for rows.Next() {
		var notification domain.Notification
		if err := rows.Scan(&notification.LoanID, &notification.Title, &notification.Status, &notification.Message, &notification.Timestamp); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}
