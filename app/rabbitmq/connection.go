package rabbitmq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type Connection struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func MakeConnection() (*Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") // Создаем подключение к RabbitMQ
	if err != nil {
		log.Fatalf("unable to open connect to RabbitMQ server. Error: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel. Error: %s", err)
	}

	return &Connection{
		Connection: conn,
		Channel:    ch,
	}, nil
}

func (c *Connection) CloseAll() {
	_ = c.Connection.Close()
	_ = c.Channel.Close()
}

type Producer struct {
	conn  *Connection
	Queue string
}

func MakeProducer(queue string) (*Producer, error) {

	c, err := MakeConnection()
	if err != nil {
		return nil, err
	}

	// DECLARE
	_, err = c.Channel.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	return &Producer{
		conn:  c,
		Queue: queue,
	}, nil
}

func (p *Producer) Close() {
	p.conn.CloseAll()
}

func (p *Producer) SendMessage(message string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.conn.Channel.PublishWithContext(ctx,
		"",      // exchange
		p.Queue, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		return err
	}

	return nil
}

type Consumer struct {
	conn  *Connection
	queue string
}

func MakeConsumer(queue string) (*Consumer, error) {

	c, err := MakeConnection()
	if err != nil {
		return nil, err
	}

	// DECLARE
	_, err = c.Channel.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn:  c,
		queue: queue,
	}, nil
}

func (c *Consumer) Close() {
	c.conn.CloseAll()
}

func (c *Consumer) RunWithHandler(handler func(message amqp.Delivery)) (messages chan amqp.Delivery, err error) {
	delivery, err := c.conn.Channel.Consume(
		c.queue, // Queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return nil, err
	}

	// wrapper
	go func() {
		for message := range delivery {
			if handler != nil {
				handler(message)
			}
		}
	}()

	return messages, nil
}
