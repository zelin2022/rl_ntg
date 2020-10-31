package matchmaking

import (
  "../channelstructs"
  "log"
  "../agent"
  "../match"
  "github.com/google/uuid"
)

type MM struct {
  onlineAgents []agent.Agent
  inGameAgents []agent.Agent
  waitingAgents []agent.Agent
  matches []match.Match
  chanCollectMatchFinish chan string
}

func (mm *MM) Run (chanAgentStatusIntake <-chan channelstructs.ServerIn) {
  var err error
  chanCollectMatchFinish = make(chan string)
  for {
    select {
      case msg := chanAgentStatusIntake: // from amqplistener, tells you some agent satus
      log.Printf(myutli.TimeStamp() + " MM Receive: src:AgentStatusIntake, header:" + msg.Header + ", agent:" + msg.AgentID)
      err = mm.updateAgents(msg)
      myutli.FailOnError(err, "MM.updateAgent failed, msg:" + msg)
    case msg := chanCollectMatchFinish:
      log.Printf("Match finished" + msg.ID )
      // remove finished match

      default: // do nothing, unblocking
    }

    err = attemptMatchMaking()
    myutli.FailOnError(err, "MM.attemptMatchMaking failed, msg:" + msg)
    err = dropAFK()
    myutli.FailOnError(err, "MM.dropAFK failed, msg:" + msg)
    // update matches channel

  }
}

// MATCH MAKING +++++++++++++++++++++++++++++++++++

func (mm *MM)attemptMatchMaking(){
  playersPositionsInWait := mm.mmStrategy0()
  if playersPositionsInWait == nil {
    return
  }
  var playersToPlay []agent.Agent

  for i := range playersPositionsInWait{
    playersToPlay = append(playersToPlay, mm.waitingAgents[i])
    // also move from waiting to ingame
    mm.inGameAgents = append(mm.inGameAgents, mm.waitingAgents[i])
    mm.waitingAgents = agent.DeleteAgent(mm.waitingAgents, i)
  }

  // create a new match
  var newMatch match.Match
  newMatch.ID = uuid.New().String()
  newMatch.ChanListenerIntake = make(chan channelstructs.ServerIn)
  newMatch.ChanFinish = mm.chanCollectMatchFinish
  newMatch.Players = playersToPlay
  newMatch.StartTime = time.Now()

  // add match to match list
  mm.matches = append(mm.matches, newMatch)

  // start the match
}

func (mm *MM)mmStrategy0() int, int{
  if len(mm.waitingAgents) > 2 {
    return 0, 1
  }
  return nil, nil
}

// UPDATE AGENTS ===================================

func (mm *MM) updateAgents (serverIn channelstructs.ServerIn) error {
  err := nil
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
    err = error.New("header is invalid")
  }
  return err
}

func (mm *MM) agentSignIn(myAgent agent.Agent) error{
  found, pos := FindAgent(mm.onlineAgents, myAgent)
  if found {
    return error.New("agentSignIn, but agent is already online")
  }
  mm.onlineAgents = append (mm.onlineAgents, myAgent)
  return nil
}

func (mm *MM) agentSignOut(myAgent agent.Agent){
  found, pos := FindAgent(mm.onlineAgents, myAgent)
  if !found {
    return error.New("agentSignOut, but agent is not online")
  }
  mm.onlineAgents = DeleteAgent(mm.onlineAgents, pos)
  return nil
}

func (mm *MM) agentWaiting(myAgent agent.Agent){
  found, pos := FindAgent(mm.waitingAgents, myAgent)
  if !found {
    mm.waitingAgents = append (mm.waitingAgents, myAgent)
    } else {
      mm.waitingAgents[pos] = myAgent // update anyway
    }
    return nil
  }
