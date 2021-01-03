package channelstructs

import (
  "../agent"
)

type SenderMessage struct {
  Header string `json:"header"`
  Body string `json:"body"`
  SendTime int64 `json:"stime"`
}
type SenderIntake struct {
  Message SenderMessage
  AgentsToSend []agent.Agent // could be one or multiple
}
