package application

import (
	"fmt"
	"log"
	"time"

	domain "notifications/src/domain/entities"
	"notifications/src/domain/repositories"
	"notifications/src/infraestructure/services"
)

type ProcessLoanUseCase struct {
	Repo               repositories.NotificationRepository
	Rabbit             *services.RabbitMQService
	PublishServiceAPI2 *services.RabbitMQPublishService
}

func NewProcessLoanUseCase(repo repositories.NotificationRepository, rabbit *services.RabbitMQService, publishServiceAPI2 *services.RabbitMQPublishService) *ProcessLoanUseCase {
	return &ProcessLoanUseCase{
		Repo:               repo,
		Rabbit:             rabbit,
		PublishServiceAPI2: publishServiceAPI2, // No necesitas crear una nueva instancia aquí
	}
}

func (uc *ProcessLoanUseCase) StartProcessingLoans() {
	for {
		time.Sleep(5 * time.Second)

		loanEvents, err := uc.Rabbit.FetchAlerts()
		if err != nil {
			log.Println("Error al obtener eventos de préstamos desde RabbitMQ:", err)
			continue
		}

		for _, event := range loanEvents {
			log.Printf("Evento recibido: %+v\n", event)

			loanIDInterface, ok := event["loan_id"]
			if !ok {
				log.Println("Error: loan_id no encontrado en el evento:", event)
				continue
			}

			loanID, ok := loanIDInterface.(float64)
			if !ok {
				log.Println("Error: loan_id no es del tipo esperado:", loanIDInterface)
				continue
			}
			loanIDInt := int32(loanID)

			title, ok := event["title"].(string)
			if !ok {
				log.Println("Error: 'title' no es del tipo esperado:", event["title"])
				continue
			}

			status, ok := event["status"].(string)
			if !ok {
				log.Println("Error: 'status' no es del tipo esperado:", event["status"])
				continue
			}

			message := fmt.Sprintf("El préstamo '%s' ha cambiado su estado a: %s", title, status)
			notification := domain.Notification{
				LoanID:    loanIDInt,
				Title:     title,
				Status:    status,
				Message:   message,
				Timestamp: time.Now().Format(time.RFC3339),
			}

			log.Printf("Notificación generada: %+v\n", notification)

			if err := uc.Repo.Create(&notification); err != nil {
				log.Println("Error guardando notificación en MySQL:", err)
			} else {
				log.Println("Notificación guardada correctamente en MySQL:", message)
			}

			// Publicar en la cola de API2 usando el método correcto
			err = uc.PublishServiceAPI2.PublishToAPI2(&notification) // Cambié esto a PublishToAPI2
			if err != nil {
				log.Println("Error publicando en la cola de API2:", err)
			} else {
				log.Println("Notificación publicada en la cola de API2:", message)
			}
		}
	}
}
