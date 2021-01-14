package matchmaking

import (
  "../match"
)



func Create(channels ChannelBundle, activeMatches *match.ActiveMatches){
  mm := MM{
    Channels: channels,
    pActiveMatches: activeMatches,
    agentLastWaitingStatusTime: make(map[string]int64),
  }

  go mm.run()
}
