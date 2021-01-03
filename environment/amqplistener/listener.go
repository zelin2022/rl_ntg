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
  ChanAMQP <-chan amqp.Delivery
}

type AMQPListener struct{
  PActiveMatches *match.ActiveMatches
  Channels ChannelBundle
}

func (ls *AMQPListener)Run() {
  for {
    log.Printf("select in")
    select { // I think for this one, blocking is fine
    case msg := <- ls.Channels.ChanAMQP:
      log.Printf("Received a message: %s", msg.Body)
      err := ls.processMessage(msg.Body, myutil.TimeStamp())
      myutil.FailOnError(err, "Failed to processMessage" + string(msg.Body))
    }
    log.Printf("select out")
  }
}


func (ls *AMQPListener)processMessage(body []byte, recvTime string) error {

  var serverIn channelstructs.ListenerOutput
  err := json.Unmarshal([]byte(body), &serverIn)
  if err != nil {
    return err
  }

  serverIn.RecvTime = recvTime

  switch serverIn.Header {
  case p_HEADER_AGENT_SIGN_IN: // send to match-making
    ls.Channels.ChanLS2MM <- serverIn
  case p_HEADER_AGENT_WAITING: // send to match-making
    ls.Channels.ChanLS2MM <- serverIn
  case p_HEADER_AGENT_SIGN_OUT: // send to match-making
    ls.Channels.ChanLS2MM <- serverIn
  case p_HEADER_AGENT_MOVE: // send to match
    ls.PActiveMatches.Mutex.Lock()
    match_pos := match.FindMatchByAgentID(ls.PActiveMatches.Matches, serverIn.AgentID)
    ls.PActiveMatches.Mutex.Unlock()
    if match_pos < 0 {
      return errors.New("received move this agent ID: " + serverIn.AgentID + "\nbut this agent is not in match")
    }
    ls.PActiveMatches.Matches[match_pos].Channels.ChansLS2MS <- serverIn
  default:

  }

  return err
}
