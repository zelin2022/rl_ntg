package channelstructs

import (
  "../agent"
)
type SenderIntake struct {
  Header string
  MoveOwnerID string
  Move string
  AgentsToSend []agent.Agent // could be one or multiple
}
