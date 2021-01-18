package amqpmaster


import (
  "../myutil"
  "../amqplistener"
  "../amqpsender"
  "../channelstructs"
  "../match"
  "github.com/streadway/amqp"

  "log"

)


type ChannelBundle struct{
  ChanLS2MM chan channelstructs.ListenerOutput
  ChanMS2SE chan channelstructs.SenderIntake
  ChanMM2SE chan channelstructs.SenderIntake
}

func Create(channels ChannelBundle, listener_queue string, sender_queue string, activeMatches *match.ActiveMatches) (close func()) {
  conn, err := amqp.Dial("amqp://test:test@localhost:5672/")
  myutil.PanicOnError(err, "Failed to connect to RabbitMQ")
  ch, err := conn.Channel()
  myutil.PanicOnError(err, "Failed to open a channel")

  q, err := ch.QueueDeclare(
    listener_queue, // name
    false,   // durable
    true,   // delete when unused
    true,   // exclusive
    false,   // no-wait
    nil,     // arguments
  )
  myutil.PanicOnError(err, "Failed to declare a queue")

  purged, err := ch.QueuePurge(
    listener_queue,
    false, // no-wait
  )
  myutil.FailOnError(err, "Failed to purge queue")
  if err == nil {
    log.Printf("Purged %d number of messages from queue %s", purged, listener_queue)
  }

  ChanConsumeCallback, err := ch.Consume(
    q.Name, // queue
    "",     // consumer
    true,   // auto-ack
    true,  // exclusive
    false,  // no-local
    false,  // no-waitamqpsender
    nil,    // args
  )
  myutil.PanicOnError(err, "Failed to register a consumer")

  LSChannels := amqplistener.ChannelBundle{
    ChanLS2MM: channels.ChanLS2MM,
    ChanAMQP: ChanConsumeCallback,
  }

  listener := amqplistener.AMQPListener{
    Channels: LSChannels,
    PActiveMatches: activeMatches,
  }

  SEChannels := amqpsender.ChannelBundle{
    ChanMS2SE: channels.ChanMS2SE,
    ChanMM2SE: channels.ChanMM2SE,
  }

  sender := amqpsender.AMQPSender{
    Channels: SEChannels,
    AMQP: *ch,
  }

  go listener.Run()
  log.Printf("[*] Waiting for messages. To exit press CTRL+C")

  go sender.Run()

  return func (){
    ch.Close()
    conn.Close()
  }
}
