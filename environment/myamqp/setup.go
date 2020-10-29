package myamqp

import (
  "../myerror"
  "github.com/streadway/amqp"

  "log"

)

var QUEUE_NAME = "server_in"

func Init(processFunc func(body []byte) error) (close func()) {
  conn, err := amqp.Dial("amqp://test:test@localhost:5672/")
  myerror.FailOnError(err, "Failed to connect to RabbitMQ")
  ch, err := conn.Channel()
  myerror.FailOnError(err, "Failed to open a channel")

  q, err := ch.QueueDeclare(
    QUEUE_NAME, // name
    false,   // durable
    true,   // delete when unused
    true,   // exclusive
    false,   // no-wait
    nil,     // arguments
  )
  myerror.FailOnError(err, "Failed to declare a queue")

  msgs, err := ch.Consume(
    q.Name, // queue
    "",     // consumer
    true,   // auto-ack
    true,  // exclusive
    false,  // no-local
    false,  // no-wait
    nil,    // args
  )
  myerror.FailOnError(err, "Failed to register a consumer")

  go onConsume(msgs, processFunc)

  log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

  return func (){
    ch.Close()
    conn.Close()
  }
}

func onConsume (delivery <-chan amqp.Delivery, processFunc func(body []byte) error){
  for d := range delivery {
    log.Printf("Received a message: %s", d.Body)
    err := processFunc(d.Body)
    myerror.FailOnError(err, "Failed to process message")
  }
}
