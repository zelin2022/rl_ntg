package matchmaking

import (
  "../channelstructs"
  "../match"
)

type ChannelBundle struct{
  ChanLS2MM chan channelstructs.ListenerOutput
  ChanMM2LS chan []match.Match
  ChanMS2RK chan string
  ChanMS2SE chan channelstructs.SenderIntake
  ChanMM2SE chan channelstructs.SenderIntake
  ChanMS2MM chan string
}

func Create(channels ChannelBundle){
  var mm MM
  mm.Channels = channels
  go mm.run()
}
