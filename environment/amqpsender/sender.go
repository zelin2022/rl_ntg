package amqpsender

import(
  "../channelstructs"
  "encoding/json"
  "../myutil"
  "github.com/streadway/amqp"
)

type ChannelBundle struct{
  ChanMS2SE chan channelstructs.SenderIntake
  ChanMM2SE chan channelstructs.SenderIntake
}

type AMQPSender struct{
  Channels ChannelBundle
  AMQP amqp.Channel
}


func (se *AMQPSender)Run(){
  for{
    select{
    case msg := <- se.Channels.ChanMS2SE:
      se.send_SenderIntake(msg)
    case msg := <- se.Channels.ChanMM2SE:
      se.send_SenderIntake(msg)
    }
  }
}


func (se *AMQPSender)send_SenderIntake(toSend channelstructs.SenderIntake){
  jsonString, err := json.Marshal(toSend.Message)
  myutil.FailOnError(err, "Failed to JSON Marshal a struct")
  for i := 0; i < len(toSend.AgentsToSend); i++{
    err = se.sendString(string(jsonString), toSend.AgentsToSend[i].Queue)
    myutil.FailOnError(err, "Fail to send to agent: " + toSend.AgentsToSend[i].ID +
      "\nqueue: " + toSend.AgentsToSend[i].Queue +
      "\nmessage: \n" + string(jsonString))
  }
}

func (se *AMQPSender)sendString(msg string, targetQueue string) error{
  err := se.AMQP.Publish(
  "",     // exchange
  targetQueue, // routing key
  false,  // mandatory
  false,  // immediate
  amqp.Publishing {
    ContentType: "text/plain",
    Body:        []byte(msg),
  })
  return err
}
