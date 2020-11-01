package match
import (
  "time"
  "../agent"
  "../channelstructs"
)

type Match struct {
  ID string
  ChanListenerIntake chan channelstructs.ListenerOutput
  ChanSenderOutlet chan channelstructs.SenderIntake
  ChanFinish chan string // tells MM this match is over
  Players []agent.Agent
  // TheGame game.Game
  StartTime time.Time
  allowedNextMovePlayersPositions []int
}

func (m *Match) run () {
  for {
    select {
    case msg := ChanListenerIntake //

    }
  }
}

func (m *Match) doMove(msg channelstructs.ListenerOutput) error {
  found, playerPosition := agent.FindAgent(m.Players, msg.AgentID)
  if !found {
    return error.New("Received message from Agent:" + msg.AgentID + "\nbut agent is not a player")
  }
  // forward the move to game
  game.TryMove(msg.AgentID, msg.Move)

  // respond to other players
  // since tictactoe and go are perfect information games
  // we dont need game to come up with a response
  m.sendToOtherPlayers(msg.AgentID, msg.Move, playerPosition)

}

func (m *Match) sendToOtherPlayers( moveAgentID string, move string, moveAgentPosition int ) {
  var sendPackage channelstructs.SenderIntake
  sendPackage.MoveOwnerID = moveAgentID
  sendPackage.Move = move
  sendPackage.AgentsToSend = agent.DeleteAgent(m.Players, moveAgentPosition)
  ChanSenderOutlet <- sendPackage
}
















// HELPR METHOD
func findInt(list []int, item int) bool, int {
  found, pos := false, nil
  for i := 0; i < len(list); i++ {
    if item == list[i] {
      found, pos = true, i
    }
  }
  return found, pos
}
