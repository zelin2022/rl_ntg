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
  pActiveMatches *match.ActiveMatches
  Channels ChannelBundle
  minimumWaitTimeForAnotherMatchMaking float64
}

type ChannelBundle struct{
  ChanLS2MM chan channelstructs.ListenerOutput
  ChanMS2RK chan channelstructs.MatchRecord
  ChanMS2SE chan channelstructs.SenderIntake
  ChanMM2SE chan channelstructs.SenderIntake
  ChanMS2MM chan string
}

func (mm *MM) run () {
  var err error
  var minDiffSelfUpdateTime int64 = int64(mm.minimumWaitTimeForAnotherMatchMaking * 1000)
  nextSelfUpdateTime := myutil.GetCurrentEpochMilli() + minDiffSelfUpdateTime


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
        // remove finished match w/ mutex

      default:
        currentTime := myutil.GetCurrentEpochMilli()
        timeDiff := nextSelfUpdateTime - currentTime
        if timeDiff <= 0  {
          err = mm.selfUpdate()
          myutil.FailOnError(err, "MatchMaking, selfUpdate")
          // if mm was successful, we chain it, by not incrementing nextSelfUpdateTime
          if err != nil { // else we induce sleep
            nextSelfUpdateTime = myutil.GetCurrentEpochMilli() + minDiffSelfUpdateTime
          }

        }else{
          myutil.Sleep("matchmaking", timeDiff)
        }

    }
  }
}

func (mm *MM)selfUpdate()error{
    // first drop afks before matchmaking
    var err error
    err = dropAFK()
    myutil.FailOnError(err, "MM.dropAFK failed")
    err = mm.attemptMatchMaking()
    myutil.FailOnError(err, "MM.attemptMatchMaking failed")
    return err
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
  mm.pActiveMatches.Mutex.Lock()
  mm.pActiveMatches.Matches = append(mm.pActiveMatches.Matches, newMatch)
  mm.pActiveMatches.Mutex.Unlock()

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
  case MSG_AGENT_SIGN_IN:
    err = mm.agentSignIn(theAgent)
  case MSG_AGENT_WAITING:
    err = mm.agentSignOut(theAgent)
  case MSG_AGENT_SIGN_OUT:
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
