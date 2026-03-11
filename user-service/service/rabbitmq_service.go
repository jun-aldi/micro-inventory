package service

import (
	"context"
	"encoding/json"
	"fmt"
	"micro-inventory/user-service/configs"

	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"
)

type EmailPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Type     string `json:"type"`
	UserID   uint   `json:"user_id"`
	Name     string `json:"name"`
}

type RabbitMQServiceInterface interface {
	PublishMail(ctx context.Context, payload EmailPayload) error
	Close() error
}

type rabbitMQService struct {
	connection *amqp.Connection
	ch         *amqp.Channel
	config     configs.Config
}

// Close implements RabbitMQServiceInterface.
func (r *rabbitMQService) Close() error {
	if r.ch != nil {
		return r.ch.Close()
	}

	if r.connection != nil {
		return r.connection.Close()
	}
	return nil
}

// PublishMail implements RabbitMQServiceInterface.
func (r *rabbitMQService) PublishMail(ctx context.Context, payload EmailPayload) error {
	// Convert Payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("[RabbitMQService] PublishMail -1: %v", err)
		return err
	}

	// Declare queue if not exists
	_, err = r.ch.QueueDeclare(
		"email_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Errorf("[RabbitMQService] PublishMail -2: %v", err)
		return err
	}

	// Publish to queue
	err = r.ch.Publish(
		"",
		"email_queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonPayload,
		},
	)
	if err != nil {
		log.Errorf("[RabbitMQService] PublishMail -3: %v", err)
		return err
	}

	return nil
}

func NewRabbitMQService(config configs.Config) (RabbitMQServiceInterface, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", config.RabbitMQ.Username, config.RabbitMQ.Password, config.RabbitMQ.Host, config.RabbitMQ.Port))
	if err != nil {
		log.Errorf("[RabbitMQService] NewRabbitMQService -1: %v", err)
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("[RabbitMQService] NewRabbitMQService -2: %v", err)
		return nil, err
	}

	return &rabbitMQService{
		connection: conn,
		ch:         ch,
		config:     config,
	}, nil

}
