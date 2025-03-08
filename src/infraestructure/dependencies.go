package infraestructure

import (
	"log"
	"notifications/src/application"
	"notifications/src/core"
	"notifications/src/infraestructure/services"
)

type Dependencies struct {
	ProcessLoanUseCase *application.ProcessLoanUseCase
}

func NewDependencies() (*Dependencies, error) {
	db, err := core.InitDb()
	if err != nil {
		return nil, err
	}

	err = core.InitRabbitMQ()
	if err != nil {
		log.Fatal("Error iniciando RabbitMQ:", err)
	}

	rabbitService := services.NewRabbitMQService()

	// Crear servicio de publicación en la cola "notifications"
	rabbitPublishServiceAPI2, err := services.NewRabbitMQPublishService("notifications")
	if err != nil {
		log.Fatal("Error iniciando RabbitMQ Publish Service para API2:", err)
	}

	mysqlRepo := NewMySQLNotificationRepository(db)

	processLoanUseCase := application.NewProcessLoanUseCase(mysqlRepo, rabbitService, rabbitPublishServiceAPI2)

	go processLoanUseCase.StartProcessingLoans()

	return &Dependencies{
		ProcessLoanUseCase: processLoanUseCase,
	}, nil
}
