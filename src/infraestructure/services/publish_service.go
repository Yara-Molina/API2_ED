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
	Queue   string
}

func NewRabbitMQPublishService(queueName string) (*RabbitMQPublishService, error) {
	// Conectar a RabbitMQ
	conn, err := amqp.Dial("amqp://yara:noobmaster69@54.161.81.210:5672/")
	if err != nil {
		log.Println("Error al conectar con RabbitMQ:", err)
		return nil, err
	}

	// Crear un canal de comunicación
	channel, err := conn.Channel()
	if err != nil {
		log.Println("Error al crear el canal de RabbitMQ:", err)
		conn.Close()
		return nil, err
	}

	// Declarar la cola una sola vez al iniciar el servicio
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
		channel.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQPublishService{
		Conn:    conn,
		Channel: channel,
		Queue:   queueName,
	}, nil
}

func (p *RabbitMQPublishService) PublishToAPI2(notification *domain.Notification) error {
	if p.Channel == nil {
		log.Println("No hay conexión con RabbitMQ")
		return nil
	}

	log.Printf("Evento a enviar a RabbitMQ: %+v\n", notification)

	body, err := json.Marshal(notification)
	if err != nil {
		log.Println("Error al convertir la notificación a JSON:", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.Channel.PublishWithContext(ctx,
		"", p.Queue, false, false,
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

func (p *RabbitMQPublishService) Close() {
	if p.Channel != nil {
		p.Channel.Close()
	}
	if p.Conn != nil {
		p.Conn.Close()
	}
	log.Println("Conexión con RabbitMQ cerrada")
}
