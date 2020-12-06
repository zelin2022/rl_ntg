package amqplistener

import (
  "../myutil"
  "github.com/streadway/amqp"
  "log"
  "../channelstructs"
  "encoding/json"
  "../match"
  "errors"
)

type ChannelBundle struct{
  ChanLS2MM chan channelstructs.ListenerOutput
  ChanMM2LS chan []match.Match
  ChanAMQP <-chan amqp.Delivery
}

type AMQPListener struct{
  activeMatches []match.Match
  Channels ChannelBundle
}

func (ls *AMQPListener)Run() {
  for {
    select { // I think for this one, blocking is fine
    case msg := <- ls.Channels.ChanAMQP:
      log.Printf("Received a message: %s", msg.Body)
      err := ls.processMessage(msg.Body, myutil.TimeStamp())
      myutil.FailOnError(err, "Failed to processMessage" + string(msg.Body))

    case newMatches := <- ls.Channels.ChanMM2LS:
      log.Print("Listener here, just got an update of active matches from MM")
      ls.activeMatches = newMatches // update matches

    }
  }
}


func (ls *AMQPListener)processMessage(body []byte, recvTime string) error {
  // // 4 cases:
  // // Agent sign in
  // // Agent idle
  // // Agent sign off
  // // Agent move
  var serverIn channelstructs.ListenerOutput
  err := json.Unmarshal([]byte(body), &serverIn)
  if err != nil {
    return err
  }

  serverIn.RecvTime = recvTime

  switch serverIn.Header {
  case "sign in": // send to match-making
    ls.Channels.ChanLS2MM <- serverIn
  case "waiting": // send to match-making
    ls.Channels.ChanLS2MM <- serverIn
  case "sign out": // send to match-making
    ls.Channels.ChanLS2MM <- serverIn
  case "move": // send to match
    match_pos := match.FindMatchByAgentID(ls.activeMatches, serverIn.AgentID)
    if match_pos < 0 {
      return errors.New("received move this agent ID: " + serverIn.AgentID + "\nbut this agent is not in match")
    }
    ls.activeMatches[match_pos].Channels.ChansLS2MS <- serverIn
  default:

  }

  return err
}
