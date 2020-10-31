package channelstructs

import (
  "../agent"
)
type ServerOut struct {
  MoveOwnerID string
  Move string
  AgentsToSend []agent.Agent
}
