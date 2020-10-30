package match
import (
  "time"
  "../agent"
)

type Match struct {
  ID uint32
  ChanFinish chan uint32 // tells MM this match is over
  PlayerBlack agent.Agent
  PlayerWhite agent.Agent
  // TheGame game.Game
  StartTime time.Time
}
