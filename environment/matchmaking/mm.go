package matchmaking

import (
  "../channelstructs"
  "log"
  "../agent"
  "../match"
  "../myutil"
  "errors"
  "strconv"
  "sort"
  "time"
)

type MM struct {
  onlineAgents []agent.Agent
  inGameAgents []agent.Agent
  waitingAgents []agent.Agent
  pActiveMatches *match.ActiveMatches
  Channels ChannelBundle
  agentLastWaitingStatusTime map[string]int64
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
  var minDiffSelfUpdateTime int64 = int64(p_MinimumWaitTimeForAnotherMatchMaking * 1000)
  nextSelfUpdateTime := myutil.GetCurrentEpochMilli() + minDiffSelfUpdateTime


  for {
    log.Print("================MM================")
    select {
    case msg := <- mm.Channels.ChanMS2MM:
      log.Printf("Match finished: " + msg )
      mm.closeMatch(msg)
      log.Printf("Match closed: " + msg )
      break
    default:
      break
    }

    select {
    case msg := <- mm.Channels.ChanLS2MM: // from amqplistener, tells you some agent satus
      log.Printf(myutil.TimeStamp() + " MM Receive: src:ChanLS2MM, header:" + msg.Header + ", agent:" + msg.AgentID)
      err = mm.updateAgents(msg)
      myutil.FailOnError(err, "MM.updateAgent failed\nmsg.Header: " + msg.Header +
        "\nmsg.AgentID: " + msg.AgentID +
        "\nmsg.Body: " + msg.Body +
        "\nmsg.SendTime: " + msg.SendTime +
        "\nmsg.RecvTime: " + msg.RecvTime + "\n")
      break
    default:
      currentTime := myutil.GetCurrentEpochMilli()
      timeDiff := nextSelfUpdateTime - currentTime
      if timeDiff <= 0  {
        err = mm.selfUpdate()
        myutil.FailOnError(err, "mm.selfUpdate() failed")
        // if mm was successful, we chain it, by not incrementing nextSelfUpdateTime
        if err != nil { // else we induce sleep
          nextSelfUpdateTime = myutil.GetCurrentEpochMilli() + minDiffSelfUpdateTime
        }
        }else{
          myutil.Sleep("matchmaking", timeDiff)
        }
        break
      }
  }
}

func (mm *MM)selfUpdate()error{
    // first drop afks before matchmaking
    var err error
    err = mm.dropAFK()
    myutil.FailOnError(err, "MM.dropAFK failed")
    err = mm.attemptMatchMaking()
    myutil.FailOnError(err, "MM.attemptMatchMaking failed")
    return nil
}

// MATCH MAKING +++++++++++++++++++++++++++++++++++

func (mm *MM)attemptMatchMaking()(error){
  playersPositionsInWait := mm.mmStrategy0()
  if playersPositionsInWait == nil {
    return errors.New("Match making failed (possibly not enough players) (current players in waiting:" +strconv.Itoa(len(mm.waitingAgents)) + ")")
  }
  log.Printf("Match found")
  mm.createMatch(playersPositionsInWait)
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
  theAgent.Queue = serverIn.Body   // for status message, body is straight up agent_queue?
  theAgent.RenewActive()
  switch serverIn.Header {
  case p_HEADER_AGENT_SIGN_IN:
    err = mm.agentSignIn(theAgent)
  case p_HEADER_AGENT_WAITING:
    err = mm.agentWaiting(theAgent)
  case p_HEADER_AGENT_SIGN_OUT:
    err = mm.agentSignOut(theAgent)
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
  mm.agentLastWaitingStatusTime[myAgent.ID] = time.Now().Unix()
  return nil
}

// ==========================================

func (mm *MM)createMatch(playersPositionsInWait []int){
  // sort by descending order
  sort.Slice(playersPositionsInWait, func(i,j int) bool{
    return playersPositionsInWait[i] > playersPositionsInWait[j]
  })

  var playersToPlay []agent.Agent

  for i := range playersPositionsInWait{
    player := mm.waitingAgents[playersPositionsInWait[i]]
    // fetch players to array from waiting players
    playersToPlay = append(playersToPlay, player)
    // delete players from waiting
    // this is safe and won't mess with position because we sorted in descending order
    mm.waitingAgents = agent.DeleteAgent(mm.waitingAgents, playersPositionsInWait[i])
    // also add the players to in-game
    mm.inGameAgents = append(mm.inGameAgents, player)
  }

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
}

func (mm *MM)closeMatch(id string){
  //delete match
  var err error
  mm.pActiveMatches.Mutex.Lock()
  players := mm.pActiveMatches.Matches[match.FindMatchByMatchID(mm.pActiveMatches.Matches, id)].Players
  mm.pActiveMatches.Matches, err = match.DeleteMatchByMatchID(mm.pActiveMatches.Matches, id)
  mm.pActiveMatches.Mutex.Unlock()
  if err != nil {
    panic("Deleting match, but match not found in pActiveMatches")
  }
  // release players back to online
  for i := range players{
    err = mm.agentWaiting(players[i])
    myutil.FailOnError(err, "Failed to put agent back into waiting")
  }
}

func (mm *MM)dropAFK() (error){
  var cutOffTime int64 = time.Now().Unix() - p_WaitingTimeoutSecondsAgo
  for i := range mm.waitingAgents{
    if mm.agentLastWaitingStatusTime[mm.waitingAgents[i].ID] < cutOffTime{ // if last waiting is before...
      log.Printf("Dropping Agent %s for being AFK, last waiting was %d and cutoff is %d", mm.waitingAgents[i].ID, mm.agentLastWaitingStatusTime[mm.waitingAgents[i].ID], cutOffTime)
      // drop agent
      mm.waitingAgents = agent.DeleteAgent(mm.waitingAgents, i)
      // potentially delete agent from onlineAgents too
      // however sicne we don't use onlineAgents for anything, we don't need to
    }
  }
  return nil
}
