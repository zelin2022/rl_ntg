package matchmaking

import (
  "../channelstructs"
  "log"
  "../agent"
  "../match"
  "../myutil"
  "errors"
  "strconv"

)

type MM struct {
  onlineAgents []agent.Agent
  inGameAgents []agent.Agent
  waitingAgents []agent.Agent
  matches []match.Match
  Channels ChannelBundle
}

type ChannelBundle struct{
  ChanLS2MM chan channelstructs.ListenerOutput
  ChanMM2LS chan []match.Match
  ChanMS2RK chan string
  ChanMS2SE chan channelstructs.SenderIntake
  ChanMM2SE chan channelstructs.SenderIntake
  ChanMS2MM chan string
}

func (mm *MM) run () {
  var err error
  for {
    log.Print("================MM================")
    select {
      case msg := <- mm.Channels.ChanLS2MM: // from amqplistener, tells you some agent satus
        log.Printf(myutil.TimeStamp() + " MM Receive: src:ChanLS2MM, header:" + msg.Header + ", agent:" + msg.AgentID)
        err = mm.updateAgents(msg)
        myutil.FailOnError(err, "MM.updateAgent failed\nmsg.Header: " + msg.Header +
          "\nmsg.AgentID: " + msg.AgentID +
          "\nmsg.Move: " + msg.Move +
          "\nmsg.SendTime: " + msg.SendTime +
          "\nmsg.RecvTime: " + msg.RecvTime + "\n")
      case msg := <- mm.Channels.ChanMS2MM:
        log.Printf("Match finished: " + msg )
        // remove finished match
        // update matches

      default: // if there is no news coming in, then lets match make
        // first drop afks before matchmaking
        err = dropAFK()
        myutil.FailOnError(err, "MM.dropAFK failed")
        err = mm.attemptMatchMaking()
        myutil.FailOnError(err, "MM.attemptMatchMaking failed")
        if err != nil { // sleep for a bit if matchmaking failed, if successful then chain it
          myutil.Sleep("MatchMaking", 2.5)
        }else{ // if matchmaking successful then update match info
          mm.Channels.ChanMM2LS <- mm.matches
        }

    }
  }
}

// MATCH MAKING +++++++++++++++++++++++++++++++++++

func (mm *MM)attemptMatchMaking()(error){
  playersPositionsInWait := mm.mmStrategy0()
  if playersPositionsInWait == nil {
    return errors.New("Match making failed (possibly not enough players) (current players in waiting:" +strconv.Itoa(len(mm.waitingAgents)) + ")")
  }
  log.Printf("Match found")
  var playersToPlay []agent.Agent

  for i := range playersPositionsInWait{
    playersToPlay = append(playersToPlay, mm.waitingAgents[i])
    // also add the players to in game
    mm.inGameAgents = append(mm.inGameAgents, mm.waitingAgents[i])
  }
  // remove the players from waiting
  mm.waitingAgents = agent.DeleteAgents(mm.waitingAgents, playersPositionsInWait)

  matchChannels := match.ChannelBundle{
    ChanMS2RK: mm.Channels.ChanMS2RK,
    ChanMS2SE: mm.Channels.ChanMS2SE,
    ChanMS2MM: mm.Channels.ChanMS2MM,
    ChansLS2MS: make(chan channelstructs.ListenerOutput),
  }

  newMatch := match.Create(playersToPlay, matchChannels)
  mm.matches = append(mm.matches, newMatch)

  return nil
}

func (mm *MM)mmStrategy0() ([]int){
  if len(mm.waitingAgents) >= 2 {
    return []int{0, 1}
  }
  return nil
}

// UPDATE AGENTS ===================================

func (mm *MM)updateAgents(serverIn channelstructs.ListenerOutput) (error) {
  var err error = nil
  var theAgent agent.Agent
  theAgent.ID = serverIn.AgentID
  theAgent.Queue = serverIn.AgentQueue
  theAgent.RenewActive()
  switch serverIn.Header {
  case "sign in":
    err = mm.agentSignIn(theAgent)
  case "sign out":
    err = mm.agentSignOut(theAgent)
  case "waiting":
    err = mm.agentWaiting(theAgent)
  default:
    err = errors.New("header is invalid")
  }
  return err
}

func (mm *MM)agentSignIn(myAgent agent.Agent) (error){
  found, _ := agent.FindAgent(mm.onlineAgents, myAgent.ID)
  if found {
    return errors.New("agentSignIn, but agent is already online")
  }
  mm.onlineAgents = append (mm.onlineAgents, myAgent)
  return nil
}

func (mm *MM)agentSignOut(myAgent agent.Agent)(error){
  found, pos := agent.FindAgent(mm.onlineAgents, myAgent.ID)
  if !found {
    return errors.New("agentSignOut, but agent is not online")
  }
  mm.onlineAgents = agent.DeleteAgent(mm.onlineAgents, pos)
  return nil
}

func (mm *MM)agentWaiting(myAgent agent.Agent)(error){
  found, pos := agent.FindAgent(mm.waitingAgents, myAgent.ID)
  if !found {
    mm.waitingAgents = append(mm.waitingAgents, myAgent)
  } else {
    mm.waitingAgents[pos] = myAgent // update anyway
  }
  return nil
}

func dropAFK() (error){
  //pass
  log.Print("dropAFK() holder")
  return nil
}
