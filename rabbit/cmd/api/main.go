package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"proj/config"
	"proj/internal"
	"time"
)

func main() {
	log.Println("Starting server")
	cfgFile, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
		//return
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
		//return
	}
	log.Println("Config loaded")

	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	amqpChannel, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer amqpChannel.Close()
	queue, err := amqpChannel.QueueDeclare(cfg.RabbitMQ.QueryName, true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
	}
	flag := true
	for i := 0; i < cfg.RabbitMQ.CountConsumers; i++ {

		queue, err = amqpChannel.QueueDeclare(cfg.RabbitMQ.QueryName, true, false, false, false, nil)
		if err != nil {
			fmt.Println(err)
		}

		err = amqpChannel.Qos(1, 0, false)
		if err != nil {
			fmt.Println(err)
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
			fmt.Println(err)
		}

		go func() {
			log.Printf("Consumer ready, id: %d", i)
			for d := range messageChannel {
				flag = false
				//log.Printf("Received a message: %s", d.Body)

				fmt.Println(internal.Grabber(string(d.Body)))

				if err := d.Ack(false); err != nil {
					log.Printf("Error acknowledging message : %s", err)
				} else {
					log.Printf("Acknowledged message")
				}
				flag = true
			}
		}()
	}

	for {
		if flag {
			time.Sleep(200*time.Millisecond)
			var link string
			fmt.Print("Введите ссылку на сайт -> ")
			fmt.Fscan(os.Stdin, &link)

			err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(link),
			})
			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(200*time.Millisecond)
		}

	}

}
