package main

/*
rabbitmqctl stop_app
rabbitmqctl reset
rabbitmqctl start_app
*/
import (
	"log"
	"os"

	//"os"
	"strconv"
	"time"

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
	args := os.Args
	log.Printf("[*] Aquesta és l'abella: %s \n", args[1])
	//Errors
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//Persmis per produir
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"permisProduir", // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	// Pot Mel
	ch1, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	pot, err := ch1.QueueDeclare(
		"potMel", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)

	// Despertador Os
	ch2, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	despertadorOs, err := ch2.QueueDeclare(
		"DespertaOs", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)

	var fin = false
	for true {
		//log.Printf("*************************")
		for d := range msgs {
			if string(d.Body) == "-1" { //l'os indica que vol acabar
				fin = true
				d.Ack(false)
				break
			}

			time.Sleep(2 * time.Second)
			//Produiem mel
			body := string(d.Body)
			err = ch1.Publish(
				"",       // exchange
				pot.Name, // routing key
				false,    // mandatory
				false,
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "text/plain",
					Body:         []byte(body),
				})
			failOnError(err, "Failed to publish a message")
			log.Printf(" L'abella %s produeix mel %s", args[1], body)

			if string(string(d.Body)) == "9" { //SI ES LA 10 DESPIERTA EL OSO
				//despertamOs
				body := args[1]
				err = ch2.Publish(
					"",                 // exchange
					despertadorOs.Name, // routing key
					false,              // mandatory
					false,
					amqp.Publishing{
						DeliveryMode: amqp.Persistent,
						ContentType:  "text/plain",
						Body:         []byte(body),
					})
				failOnError(err, "Failed to publish a message")
				log.Printf("Despert l'os \n")
			}
			d.Ack(false)
		}

		if fin {
			//log.Printf("----------------- \n")
			break
		}

	}

}
