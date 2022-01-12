package main

// Macià Salvà Salvà i Juan Francisco Hernandez Fernandez
import (
	"log"
	"os"
	"strconv"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func bodyFrom(args []string, permis int) string {

	s := strconv.Itoa(permis)
	return s
}

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	//Coa produeix permisos
	q, err := ch.QueueDeclare(
		"permisProduir", // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	//Coa desperten l'os
	conn, err1 := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//Despertador OS
	ch1, err1 := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q1, err1 := ch1.QueueDeclare(
		"DespertaOs", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err1, "Failed to declare a queue")

	msgs1, err1 := ch.Consume(
		q1.Name, // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	failOnError(err1, "Failed to register a consumer")

	for i := 0; i < 3; i++ {
		for i := 0; i < 10; i++ {
			body := bodyFrom(os.Args, i)
			err = ch.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "text/plain",
					Body:         []byte(body),
				})
			failOnError(err, "Failed to publish a message")
		}
		log.Printf(" L'os se'n va a dormir ")
		for d := range msgs1 {
			log.Printf("Me ha despertat l'abella %s \n", string(d.Body))
			break
		}
	}

	for i := 0; i < 10; i++ {
		body := bodyFrom(os.Args, -1)
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(body),
			})
		failOnError(err, "Failed to publish a message")
		//log.Printf(" [Os] Envia Permis %s", body)
	}

}
