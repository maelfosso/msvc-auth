package rabbitmq

import (
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

var (
	rabbitMQConn    *amqp.Connection
	rabbitMQChannel *amqp.Channel
)

func Connect() {
	rabbitMQConn, err := amqp.Dial(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ : %s", err)
	}

	rabbitMQChannel, err = rabbitMQConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel : %s ", err)
	}

	exchanges := strings.Split(os.Getenv("RABBITMQ_EXCHANGES"), ",")
	log.Println("Exchanges : %v", exchanges)
	for i := range exchanges {
		err = rabbitMQChannel.ExchangeDeclare(
			exchanges[i],
			"fanout",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Panicf("Failed to dechange the exchange [%s] : %s", exchanges[i], err)
		}
	}
}

func Close() {
	rabbitMQChannel.Close()
	rabbitMQConn.Close()
}

func Publish(body []byte, exchange string) {
	err := rabbitMQChannel.Publish(
		exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message in [%s] : %s", exchange, err)
	}
}
