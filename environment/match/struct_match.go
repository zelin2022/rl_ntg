package match
import (
  "time"
  "errors"
  "../game"
  "../agent"
  "../channelstructs"

)


type Match struct {
  ID string
  Channels ChannelBundle

  Players []agent.Agent
  TheGame game.Game
  StartTime time.Time
  allowedNextMovePlayersPositions []int // which player(s) can make the next move?
}

func (m *Match) run () {
  for {
    select {
    case msg := <- m.Channels.ChansLS2MS: // every match has a unique channel to receive from listener
    m.doMove(msg)

    }
  }
}

func (m *Match) doMove(msg channelstructs.ListenerOutput) error {
  found, playerPosition := agent.FindAgent(m.Players, msg.AgentID)
  if !found {
    return errors.New("Received message from Agent:" + msg.AgentID + "\nbut agent is not a player")
  }
  // forward the move to game
  m.TheGame.TryMove(msg.AgentID, msg.Move)

  // respond to other players
  // since tictactoe and go are perfect information games
  // we dont need game to come up with a response
  m.sendToOtherPlayers(msg.AgentID, msg.Move, playerPosition)
  return nil
}

func (m *Match) sendToOtherPlayers( moveAgentID string, move string, moveAgentPosition int ) {
  var sendPackage channelstructs.SenderIntake
  sendPackage.MoveOwnerID = moveAgentID
  sendPackage.Move = move
  sendPackage.AgentsToSend = agent.DeleteAgent(m.Players, moveAgentPosition)
  m.Channels.ChanMS2SE <- sendPackage
}


















// HELPR METHOD
func FindMatchByAgentID(matches []Match, agentID string)(int){
  for  i := 0; i < len(matches); i++{
    for j:= 0; j < len(matches[i].Players); j++{  // interate over matches and players
      if matches[i].Players[j].ID == agentID {
        return i
      }
    }
  }

  return -1
}
