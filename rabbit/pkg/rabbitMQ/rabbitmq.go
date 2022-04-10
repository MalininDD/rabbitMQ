package rabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"proj/config"
	"proj/internal"
)

type rabbitManager struct {
	config *config.Config
}

func NewRabbitManager(config *config.Config) *rabbitManager {
	return &rabbitManager{config: config}
}

func (r *rabbitManager) CreateConsumer(queryName string, amqpChannel *amqp.Channel) error {
	queue, err := amqpChannel.QueueDeclare(queryName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = amqpChannel.Qos(1, 0, false)
	if err != nil {
		return err
	}
	messageChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range messageChannel {

			log.Printf("Received a message: %s", d.Body)

			fmt.Println(internal.Grabber(string(d.Body)))

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}
		}
	}()
	return nil
}
func (r *rabbitManager) CreateProducer(queryName string, body []byte) error {
	conn, err := amqp.Dial(r.config.RabbitMQ.URL)
	if err != nil {
		return err
	}
	defer conn.Close()
	amqpChannel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer amqpChannel.Close()
	queue, err := amqpChannel.QueueDeclare(queryName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
	if err != nil {
		return err
	}
	log.Printf("AddTask")

	return nil
}

//func (r * rabbitManager) PublishProducer(amqpChannel *amqp.Channel) {
//	err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
//		DeliveryMode: amqp.Persistent,
//		ContentType:  "text/plain",
//		Body:         body,
//	})
//	if err != nil {
//		return err
//	}
//}
