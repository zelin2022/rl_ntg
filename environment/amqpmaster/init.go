package amqpmaster


import (
  "../myutli"
  "../amqplistener"
  "github.com/streadway/amqp"

  "log"

)

var QUEUE_NAME = "server_in"

func Init() (close func()) {
  conn, err := amqp.Dial("amqp://test:test@localhost:5672/")
  myutli.FailOnError(err, "Failed to connect to RabbitMQ")
  ch, err := conn.Channel()
  myutli.FailOnError(err, "Failed to open a channel")

  q, err := ch.QueueDeclare(
    QUEUE_NAME, // name
    false,   // durable
    true,   // delete when unused
    true,   // exclusive
    false,   // no-wait
    nil,     // arguments
  )
  myutli.FailOnError(err, "Failed to declare a queue")

  chanQueueIntake, err := ch.Consume(
    q.Name, // queue
    "",     // consumer
    true,   // auto-ack
    true,  // exclusive
    false,  // no-local
    false,  // no-wait
    nil,    // args
  )
  myutli.FailOnError(err, "Failed to register a consumer")

  go amqplistener.Run(chanQueueIntake)

  log.Printf(myutli.TimeStamp() + " [*] Waiting for messages. To exit press CTRL+C")

  return func (){
    ch.Close()
    conn.Close()
  }
}
