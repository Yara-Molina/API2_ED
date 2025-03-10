package application

import (
	"fmt"
	"log"
	"reflect"
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
		PublishServiceAPI2: publishServiceAPI2,
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

			for key, value := range event {
				log.Printf("Campo: %s, Tipo: %T, Valor: %+v\n", key, value, value)
			}

			if _, exists := event["Status"]; !exists {
				log.Println("Clave 'Status' no encontrada en el evento")
			}
			if _, exists := event["Title"]; !exists {
				log.Println("Clave 'Title' no encontrada en el evento")
			}

			loanIDRaw, idExists := event["ID"]
			if !idExists || loanIDRaw == nil {
				log.Println("Error: 'ID' no encontrado o es nil")
			} else {
				log.Printf("'ID' encontrado: %v (Tipo: %v)", loanIDRaw, reflect.TypeOf(loanIDRaw))
			}

			var loanID int32
			if floatVal, ok := loanIDRaw.(float64); ok {
				loanID = int32(floatVal)
			} else {
				log.Println("Error: No se pudo convertir 'ID' a int32, valor recibido:", loanIDRaw)
			}

			log.Printf("LoanID convertido: %d", loanID)

			titleRaw, ok := event["Title"]
			if !ok || titleRaw == nil {
				log.Println("Error: 'Title' no encontrado o es nil")
				titleRaw = "Desconocido"
			} else {
				log.Printf("'Title' encontrado: %v (Tipo: %T)", titleRaw, titleRaw)
			}

			statusRaw, ok := event["Status"]
			if !ok || statusRaw == nil {
				log.Println("Error: 'Status' no encontrado o es nil")
				statusRaw = "Desconocido"
			} else {
				log.Printf("'Status' encontrado: %v (Tipo: %T)", statusRaw, statusRaw)
			}

			title := fmt.Sprintf("%v", titleRaw)
			status := fmt.Sprintf("%v", statusRaw)

			log.Printf("Title convertido: %s", title)
			log.Printf("Status convertido: %s", status)

			message := fmt.Sprintf("El préstamo '%s' ha cambiado su estado a: %s", title, status)
			notification := domain.Notification{
				LoanID:    loanID,
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

			err = uc.PublishServiceAPI2.PublishToAPI2(&notification)
			if err != nil {
				log.Println("Error publicando en la cola de API2:", err)
			} else {
				log.Println("Notificación publicada en la cola de API2:", message)
			}
		}
	}
}
