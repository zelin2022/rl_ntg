package main

import (
  // "log"
  "./amqpmaster"
  "./channelstructs"
  "./match"
  "./matchmaking"
  // "encoding/json"
)

func main() {
  // create all channels
  // listener => matchmaking
  chanLS2MM := make(chan channelstructs.ListenerOutput)
  // matchmaking => listener
  chanMM2LS := make(chan []match.Match)
  // matches => recordKeeping
  chanMS2RK := make(chan string)
  // matches => sender
  chanMS2SE := make(chan channelstructs.SenderIntake)



 

  mmChannels := matchmaking.ChannelBundle{
    ChanLS2MM: chanLS2MM,
    ChanMM2LS: chanMM2LS,
    ChanMS2RK: chanMS2RK,
    ChanMS2SE: chanMS2SE,
  }


  amqpChannels := amqpmaster.ChannelBundle{
    ChanLS2MM: chanLS2MM,
    ChanMM2LS: chanMM2LS,
    ChanMS2SE: chanMS2SE,
  }




  // run the modules




  close := amqpmaster.Create(amqpChannels)
  defer close() // ideally...but doesn't work for ctrl+C

  matchmaking.Create(mmChannels)


  forever := make(chan bool)
  <-forever
}
