package amqpmaster


import (
  "../myutil"
  "../amqplistener"
  "../channelstructs"
  "../match"
  "github.com/streadway/amqp"

  "log"

)

var QUEUE_NAME = "server_in"


type ChannelBundle struct{
  ChanLS2MM chan channelstructs.ListenerOutput
  ChanMM2LS chan []match.Match
  ChanMS2SE chan channelstructs.SenderIntake
  ChanMM2SE chan channelstructs.SenderIntake
}

func Create(channels ChannelBundle) (close func()) {
  conn, err := amqp.Dial("amqp://test:test@localhost:5672/")
  myutil.FailOnError(err, "Failed to connect to RabbitMQ")
  ch, err := conn.Channel()
  myutil.FailOnError(err, "Failed to open a channel")

  q, err := ch.QueueDeclare(
    QUEUE_NAME, // name
    false,   // durable
    true,   // delete when unused
    true,   // exclusive
    false,   // no-wait
    nil,     // arguments
  )
  myutil.FailOnError(err, "Failed to declare a queue")

  ChanConsumeCallback, err := ch.Consume(
    q.Name, // queue
    "",     // consumer
    true,   // auto-ack
    true,  // exclusive
    false,  // no-local
    false,  // no-wait
    nil,    // args
  )
  myutil.FailOnError(err, "Failed to register a consumer")

  LSChannels := amqplistener.ChannelBundle{
    ChanLS2MM: channels.ChanLS2MM,
    ChanMM2LS: channels.ChanMM2LS,
    ChanAMQP: ChanConsumeCallback,
  }
  go amqplistener.Run(chanQueueIntake)

  log.Printf(myutil.TimeStamp() + " [*] Waiting for messages. To exit press CTRL+C")

  return func (){
    ch.Close()
    conn.Close()
  }
}
