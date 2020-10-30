package channelstructs

type ServerIn struct {
  Header string // purpose
  AgentID string  // name of agent
  AgentQueue string  // queue of agent
  Move string // if there is move
  SendTime string // send time
  RecvTime string
}
