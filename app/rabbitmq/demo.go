package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func Demo() {

	// Receiving
	consumer, err := MakeConsumer("test")
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	_, err = consumer.RunWithHandler(handler)
	if err != nil {
		panic(err)
	}

	// Sender

	producer, err := MakeProducer("test")
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	go func() {
		num := 1
		for {
			err = producer.SendMessage(fmt.Sprintf("Message %d", num))
			if err != nil {
				fmt.Printf("Cannot send message to queue %s", producer.Queue)
			}
			num++
			time.Sleep(time.Second * time.Duration(1))
		}
	}()
	select {}
}

func handler(message amqp.Delivery) {
	fmt.Printf("[x] RECEIVED: %s\n", message.Body)
}
