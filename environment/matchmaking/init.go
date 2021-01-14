package matchmaking

import (
  "../match"
)



func Create(channels ChannelBundle, activeMatches *match.ActiveMatches){
  var mm MM
  mm.Channels = channels
  mm.pActiveMatches = activeMatches
  go mm.run()
}
