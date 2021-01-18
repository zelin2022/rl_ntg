package amqplistener

import (
  "../myutil"
  "github.com/streadway/amqp"
  "log"
  "../channelstructs"
  "encoding/json"
  "../match"
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
  for msg := range ls.Channels.ChanAMQP{
    log.Printf("Received a message: %s", msg.Body)
    err := ls.processMessage(msg.Body, myutil.TimeStamp())
    myutil.FailOnError(err, "Failed to processMessage" + string(msg.Body))
  }
}


func (ls *AMQPListener)processMessage(body []byte, recvTime string) error {
  var serverIn channelstructs.ListenerOutput
  err := json.Unmarshal([]byte(body), &serverIn)
  if err != nil {
    myutil.PanicOnError(err, "Failed to process message: " + string(body))
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
    match_pos, err2 := match.FindMatchByAgentID(ls.PActiveMatches.Matches, serverIn.AgentID)
    ls.PActiveMatches.Mutex.Unlock()
    if err2 != nil {
      return err2
    }
    TrySendToPotentiallyClosedChannel(ls.PActiveMatches.Matches[match_pos].Channels.ChansLS2MS, serverIn)
  default:

  }

  return err
}

func TrySendToPotentiallyClosedChannel(channel chan<- channelstructs.ListenerOutput, msg channelstructs.ListenerOutput){
  defer func(){
    if r := recover(); r != nil{
      log.Printf("Listener recovered a panic: %v", r)
    }
  }()

  channel <- msg
}
