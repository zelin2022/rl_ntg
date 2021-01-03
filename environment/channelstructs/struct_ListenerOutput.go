package channelstructs

type ListenerOutput struct {
  Header string `json:"header"`
  Body string `json:"body"`
  AgentID string `json:"aid"`
  SendTime string `json:"stime"`
  RecvTime string `json:"-"`
}
