package match
import (
  "github.com/google/uuid"
  "../channelstructs"
  "../agent"
  "time"
  "log"
)

type ChannelBundle struct {
  ChanMS2RK chan string   // send to record keeper
  ChanMS2SE chan channelstructs.SenderIntake // sender
  ChanMS2MM chan string // back to matchmaking
  ChansLS2MS chan channelstructs.ListenerOutput
}


func Create(playersToPlay []agent.Agent, channels ChannelBundle)(Match) {
  log.Printf("Creating a new match, players: %v ", agent.GetAllAgentIDs(playersToPlay))
  // create a new match
  newMatch := Match{
    ID: uuid.New().String(),
    Channels: channels,
    Players: playersToPlay,
    StartTime: time.Now(),
    allowedNextMovePlayersPositions: make([]int,0),
  }
  newMatch.ID = uuid.New().String()
  newMatch.Players = playersToPlay
  newMatch.StartTime = time.Now()


  go newMatch.run()
  // start the match\
  return newMatch
}
