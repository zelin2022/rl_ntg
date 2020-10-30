package amqplistener

import (
  "../myutli"
  "github.com/streadway/amqp"
  "log"
  "time"
  "../channelstructs"
)

func Run(queueIntake <-chan amqp.Delivery) {
  for {
    select {
    case msg := <- queueIntake:
      log.Printf(myutli.TimeStamp() + " Received a message: %s", msg.Body)
      err := processMessage(msg.Body. time.Now())
      myutli.FailOnError(err, "Failed to processMessage" + string(body))
    }
  }
}


func processMessage(body []byte, recvTime time.Time) error {
  // // 4 cases:
  // // Agent sign in
  // // Agent idle
  // // Agent sign off
  // // Agent move
  var serverIn channelstructs.ServerIn
  err := json.Unmarshal([]byte(body), &serverMsg)
  myutli.FailOnError(err, "Failed to unmarshal to json" + string(body))
  serverIn.RecvTime = time.String(recvTime)

  switch serverIn.Header {
  case "sign in": // send to match-making
  case "waiting": // send to match-making
  case "sign out": // send to match-making
  case "move": // send to match

  }

  return err
}
