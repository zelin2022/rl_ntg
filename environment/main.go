package main

import (
  // "log"
  "./amqpmaster"
  "./channelstructs"
  "./match"
  "./matchmaking"
  "log"
  // "encoding/json"
)

var QUEUE_AGENT_2_SERVER string = "server_in_0"
var QUEUE_SERVER_2_AGENT string = "server_out_0"

func main() {
  // add milliseconds logger timestamp
  log.SetFlags(log.LstdFlags | log.Lmicroseconds)

  // create all channels
  // listener => matchmaking
  chanLS2MM := make(chan channelstructs.ListenerOutput)
  // matches => recordKeeping
  chanMS2RK := make(chan string)
  // matches => sender
  chanMS2SE := make(chan channelstructs.SenderIntake)
  // mcathcmkaing => sender
  chanMM2SE := make(chan channelstructs.SenderIntake)





  mmChannels := matchmaking.ChannelBundle{
    ChanLS2MM: chanLS2MM,
    ChanMS2RK: chanMS2RK,
    ChanMS2SE: chanMS2SE,
    ChanMM2SE: chanMM2SE,
  }


  amqpChannels := amqpmaster.ChannelBundle{
    ChanLS2MM: chanLS2MM,
    ChanMS2SE: chanMS2SE,
    ChanMM2SE: chanMM2SE,
  }




  // run the modules



  // active matches struct is thread safe match slice
  var activeMatches match.ActiveMatches

  close := amqpmaster.Create(amqpChannels, QUEUE_AGENT_2_SERVER, QUEUE_SERVER_2_AGENT, &activeMatches)
  defer close() // ideally...but doesn't work for ctrl+C

  matchmaking.Create(mmChannels, &activeMatches)


  forever := make(chan bool)
  <-forever
}
