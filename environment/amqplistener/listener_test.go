package amqplistener
import (
  "testing"
  "time"
  "../channelstructs"
)

func TestTrySendToPotentiallyClosedChannel (t *testing.T){ // quite pround of this one
  recvr := make(chan channelstructs.ListenerOutput)

  // recevs when channel is open
  go func(){
    TrySendToPotentiallyClosedChannel(recvr, channelstructs.ListenerOutput{
      Header: "i am me",
    })
  }()
  // ^ looks like some javascript ugly
  time.Sleep(20 * time.Millisecond) // sleep a bit
  select {
  case m := <- recvr:
    if m.Header != "i am me" {
      t.Error("TrySendToPotentiallyClosedChannel receives on the other side but wrong data")
    }
  default:
    t.Error("TrySendToPotentiallyClosedChannel failed to have other side receive when channel is open")
  }

  // does not block when channel is closed
  close(recvr)
  blocking := true
  go func(){
    TrySendToPotentiallyClosedChannel(recvr, channelstructs.ListenerOutput{
      Header: "i am me",
    })
    blocking = false
  }()
  time.Sleep(20 * time.Millisecond)
  if blocking != false {
    t.Error("TrySendToPotentiallyClosedChannel failed to not block when channel is closed")
  }
}
