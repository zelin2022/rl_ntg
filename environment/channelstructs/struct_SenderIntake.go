package channelstructs

import (
  "../agent"
)
type SenderIntake struct {
  Header string
  GamePlayers []string // only at the start of a game
  MoveOwnerID string
  Move string
  AgentsToSend []agent.Agent // could be one or multiple
}
