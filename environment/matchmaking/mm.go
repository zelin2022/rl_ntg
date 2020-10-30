package matchmaking

import (
  "../channelstructs"
  "log"
  "../agent"
  "../match"
)

type MM struct {
  onlineAgents []agent.Agent
  inGameAgents []agent.Agent
  waitingAgents []agent.Agent
  matches []match.Match
}

func (mm *MM) Run (chanAgentStatusIntake <-chan channelstructs.ServerIn) {
  var err error
  for {
    select {
    case msg := chanAgentStatusIntake:
      log.Printf(myutli.TimeStamp() + " MM Receive: src:AgentStatusIntake, header:" + msg.Header + ", agent:" + msg.AgentID)
      err = mm.updateAgents(msg)
      myutli.FailOnError(err, "MM.updateAgent failed, msg:" + msg)
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
  pos1, pos2 := mm.mmStrategy0()
  if pos1 == nil || pos2 == nil {
    return
  }

  var newMatch match.Match
  



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
  found, pos := findAgent(mm.onlineAgents, myAgent)
  if found {
    return error.New("agentSignIn, but agent is already online")
  }
  mm.onlineAgents = append (mm.onlineAgents, myAgent)
  return nil
}

func (mm *MM) agentSignOut(myAgent agent.Agent){
  found, pos := findAgent(mm.onlineAgents, myAgent)
  if !found {
    return error.New("agentSignOut, but agent is not online")
  }
  mm.onlineAgents = deleteAgent(mm.onlineAgents, pos)
  return nil
}

func (mm *MM) agentWaiting(myAgent agent.Agent){
  found, pos := findAgent(mm.waitingAgents, myAgent)
  if !found {
    mm.waitingAgents = append (mm.waitingAgents, myAgent)
  } else {
    mm.waitingAgents[pos] = myAgent // update anyway
  }
  return nil
}
