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
  "math"
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
  var timeoutBase int64 = 2
  var timeoutBaseExp int64 = 0
  nextSelfUpdateTime := myutil.GetCurrentEpochMilli() + minDiffSelfUpdateTime


  for {
    // log.Print("================MM================")

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
      myutil.PanicOnError(err, "MM.updateAgent failed\nmsg.Header: " + msg.Header +
        "\nmsg.AgentID: " + msg.AgentID +
        "\nmsg.Body: " + msg.Body +
        "\nmsg.SendTime: " + msg.SendTime +
        "\nmsg.RecvTime: " + msg.RecvTime + "\n")
      timeoutBaseExp = 0
      break
    default:
      currentTime := myutil.GetCurrentEpochMilli()
      timeDiff := nextSelfUpdateTime - currentTime
      if timeDiff <= 0  {
        err = mm.selfUpdate()
        myutil.PanicOnError(err, "selfUpdate() error")
        nextSelfUpdateTime = myutil.GetCurrentEpochMilli() + minDiffSelfUpdateTime
        timeoutBaseExp = 0
      }else
      {
        to_sleep := int64(math.Pow(float64(timeoutBase), float64(timeoutBaseExp))) - 1
        myutil.Sleep("matchmaking", to_sleep)
        timeoutBaseExp += 1
      }
      break
    }
  }
}

func (mm *MM)selfUpdate()(error){
    // first drop afks before matchmaking
    var err error
    err = mm.dropAFK()
    myutil.FailOnError(err, "MM.dropAFK failed")
    err = mm.attemptMatchMaking()
    myutil.FailOnError(err, "MM.attemptMatchMaking failed")
    return err
}


func (mm *MM)dropAFK() (error){
  var err error
  var cutOffTime int64 = time.Now().Unix() - p_WaitingTimeoutSecondsAgo
  for i := range mm.waitingAgents{
    if mm.agentLastWaitingStatusTime[mm.waitingAgents[i].ID] < cutOffTime{ // if last waiting is before...
      log.Printf("Dropping Agent %s for being AFK, last waiting was %d and cutoff is %d", mm.waitingAgents[i].ID, mm.agentLastWaitingStatusTime[mm.waitingAgents[i].ID], cutOffTime)
      // drop agent
      mm.waitingAgents, err = agent.DeleteAgent(mm.waitingAgents, i)
      myutil.FailOnError(err, "Error while dropping agents from waitingAgents")
      // potentially delete agent from onlineAgents too
      // however sicne we don't use onlineAgents for anything, we don't need to
    }
  }
  return nil
}

// MATCH MAKING +++++++++++++++++++++++++++++++++++

func (mm *MM)attemptMatchMaking()(error){
  playersPositionsInWait2D, err := mm.mmStrategy0()
  if err != nil {
    myutil.FailOnError(err, "Matchmaking strategy failed (current players in waiting:" +strconv.Itoa(len(mm.waitingAgents)) + ")")
    return nil
  }
  log.Printf("Matches found: total %d", len(playersPositionsInWait2D))
  for i := range playersPositionsInWait2D{
    err = mm.createMatch(playersPositionsInWait2D[i])
    if err != nil {
      return err
    }
  }
  err = mm.movePlayersAfterCreateMatches(playersPositionsInWait2D)
  return err
}

func (mm *MM)mmStrategy0() ([][]int, error){
  var out2D [][]int
  if len(mm.waitingAgents) >= 2 {
    for i := 0; i < (len(mm.waitingAgents) / 2); i++ {
        out2D = append(out2D, []int{ 2 * i , 2 * i + 1 })
    }
    return out2D, nil
  }
  return out2D, errors.New("mmStrategy0 failed Not enough players")
}

// UPDATE AGENTS ===================================

func (mm *MM)updateAgents(serverIn channelstructs.ListenerOutput) (error) {
  var err error = nil
  var theAgent agent.Agent
  theAgent.ID = serverIn.AgentID
  switch serverIn.Header {
  case p_HEADER_AGENT_SIGN_IN:
    err = mm.agentSignIn(theAgent)
  case p_HEADER_AGENT_WAITING:
    extractStatusMessageBody(&theAgent, serverIn.Body)
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
  var err error
  found, pos := agent.FindAgent(mm.onlineAgents, myAgent.ID)
  if !found {
    return errors.New("agentSignOut, but agent is not online")
  }
  mm.onlineAgents, err = agent.DeleteAgent(mm.onlineAgents, pos)
  return err
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

func extractStatusMessageBody(theAgent *agent.Agent, body string){
  message, err := toStatusBody(body)
  myutil.PanicOnError(err, "extract status message body")
  theAgent.Queue = message.Queue
  theAgent.MMC = message.MMC
}

// ==========================================

func (mm *MM)createMatch(playersPositionsInWait []int) error {
  // sort by descending order
  sort.Slice(playersPositionsInWait, func(i,j int) bool{
    return playersPositionsInWait[i] > playersPositionsInWait[j]
  })

  var playersToPlay []agent.Agent

  for i := range playersPositionsInWait{
    player := mm.waitingAgents[playersPositionsInWait[i]]
    // fetch players to array from waiting players
    playersToPlay = append(playersToPlay, player)
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
  return nil
}

func (mm *MM)movePlayersAfterCreateMatches(players2D [][]int)error{
  var players []int
  for i := range players2D{
    for j := range players2D[i]{
      players = append(players, players2D[i][j])
    }
  }
  // add the players to in-game
  for k := range players{
    mm.inGameAgents = append(mm.inGameAgents, mm.waitingAgents[players[k]])
  }
  // delete players from waiting
  // this is safe and won't mess with position because we sorted in descending order
  var err error
  mm.waitingAgents, err = agent.DeleteAgents(mm.waitingAgents, players)
  if err != nil {
    myutil.FailOnError(err, "Error while deleting from waitingAgents")
    return err
  }
  return nil
}

func (mm *MM)closeMatch(id string){
  //delete match
  var err error
  mm.pActiveMatches.Mutex.Lock()
  pos, err := match.FindMatchByMatchID(mm.pActiveMatches.Matches, id)
  if err != nil {
    mm.pActiveMatches.Mutex.Unlock()
    myutil.PanicOnError(err, "Error when finding match by match ID")
  }
  players := mm.pActiveMatches.Matches[pos].Players
  mm.pActiveMatches.Matches, err = match.DeleteMatchByMatchID(mm.pActiveMatches.Matches, id)
  mm.pActiveMatches.Mutex.Unlock()
  if err != nil {
    myutil.PanicOnError(err, "Error when deleting match by match ID")
  }
  // release players back to online
  // no need to put them to waiting, because agent should send waiting after game ends
  for i := range players{
    mm.inGameAgents, err = agent.DeleteAgentByID(mm.inGameAgents, players[i].ID)
    if err != nil {
      myutil.FailOnError(err, "Failed to delete agent from inGameAgents")
    }
  }
}
