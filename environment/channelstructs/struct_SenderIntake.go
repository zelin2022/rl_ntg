package channelstructs

import (
  "../agent"
)

type SenderMessage struct {
  Header string
  Body string

}
type SenderIntake struct {
  Message SenderMessage
  AgentsToSend []agent.Agent // could be one or multiple
}
