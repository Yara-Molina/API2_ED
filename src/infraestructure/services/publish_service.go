package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	domain "notifications/src/domain/entities"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublishService struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

// NewRabbitMQPublishService crea y devuelve un nuevo servicio de publicación en RabbitMQ
func NewRabbitMQPublishService(queueName string) (*RabbitMQPublishService, error) {
	// Conectar a RabbitMQ
	conn, err := amqp.Dial("amqp://yara:noobmaster69@54.161.81.210:5672/") // Cambia según tu configuración de RabbitMQ
	if err != nil {
		log.Println("Error al conectar con RabbitMQ:", err)
		return nil, err
	}

	// Crear un canal de comunicación
	channel, err := conn.Channel()
	if err != nil {
		log.Println("Error al crear el canal de RabbitMQ:", err)
		return nil, err
	}

	// Declarar la cola si no existe
	_, err = channel.QueueDeclare(
		queueName, // Nombre de la cola
		true,      // Durable (persistente)
		false,     // Auto-delete
		false,     // Exclusive
		false,     // No-wait
		nil,       // Args
	)
	if err != nil {
		log.Println("Error al declarar la cola:", err)
		return nil, err
	}

	// Devolver el servicio de publicación configurado
	return &RabbitMQPublishService{
		Conn:    conn,
		Channel: channel,
	}, nil
}

// Método para publicar notificación en la cola de API2
func (p *RabbitMQPublishService) PublishToAPI2(notification *domain.Notification) error {
	if p.Channel == nil {
		log.Println("No hay conexión con RabbitMQ")
		return nil
	}

	// Declarar la cola para API2 si no está declarada
	_, err := p.Channel.QueueDeclare(
		"notifications", // Nombre de la cola
		true,            // Durable (persistente)
		false,           // Auto-delete
		false,           // Exclusive
		false,           // No-wait
		nil,             // Args
	)
	if err != nil {
		log.Println("Error al declarar la cola notifications:", err)
		return err
	}

	body, _ := json.Marshal(notification)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.Channel.PublishWithContext(ctx,
		"", "notifications", false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Println("Error al enviar la notificación:", err)
	} else {
		log.Println("Notificación enviada correctamente:", string(body))
	}

	return err
}
