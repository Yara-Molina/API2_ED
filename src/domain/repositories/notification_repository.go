package repositories

import domain "notifications/src/domain/entities"

type NotificationRepository interface {
	Create(notification *domain.Notification) error
	GetAll() ([]domain.Notification, error)
}
