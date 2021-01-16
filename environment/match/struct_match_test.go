package match

import (
  "testing"
  "../agent"
  "../game"
  "time"
  "../channelstructs"
)

func TestDoMove(t *testing.T){} // game specific

func TestBroadcastStartToAllPlayers(t *testing.T){
  test_channel := make(chan channelstructs.SenderIntake)
  test_match := Match{
    Channels: ChannelBundle{
      ChanMS2SE: test_channel,
    },
    TheGame: game.NewGame([]string{"s player?"}),
  }
  go func(){
    test_match.broadcastStartToAllPlayers()
  }()
  time.Sleep(20 * time.Millisecond)
  select{
  case out := <- test_channel:
    if out.Message.Header != p_HEADER_SERVER_GAME_START {
      t.Error("broadcastStartToAllPlayers receive but wrong header")
    }
  default:
    t.Error("broadcastStartToAllPlayers failed to send through channel")
  }
}

func TestBroadcastMoveToAllPlayers(t *testing.T){
  test_channel := make(chan channelstructs.SenderIntake)
  test_match := Match{
    Channels: ChannelBundle{
      ChanMS2SE: test_channel,
    },
    TheGame: game.NewGame([]string{"s player?"}),
  }
  go func(){
    test_match.broadcastMoveToAllPlayers("a move")
  }()
  time.Sleep(20 * time.Millisecond)
  select{
  case out := <- test_channel:
    if out.Message.Header != p_HEADER_SERVER_MOVE || out.Message.Body != "a move" {
      t.Error("broadcastMoveToAllPlayers receive but wrong header")
    }
  default:
    t.Error("broadcastMoveToAllPlayers failed to send through channel")
  }
}

func TestBroadcastEndToAllPlayers(t *testing.T){
  test_channel := make(chan channelstructs.SenderIntake)
  test_match := Match{
    Channels: ChannelBundle{
      ChanMS2SE: test_channel,
    },
    TheGame: game.NewGame([]string{"s player?"}),
  }
  go func(){
    test_match.broadcastEndToAllPlayers()
  }()
  time.Sleep(20 * time.Millisecond)
  select{
  case out := <- test_channel:
    if out.Message.Header != p_HEADER_SERVER_GAME_END {
      t.Error("broadcastEndToAllPlayers receive but wrong header")
    }
  default:
    t.Error("broadcastEndToAllPlayers failed to send through channel")
  }
}

func TestSendMatchToRecordKeeper(t *testing.T){
  current_time := time.Now()
  test_channel := make(chan channelstructs.MatchRecord)
  test_match := Match{
    Channels: ChannelBundle{
      ChanMS2RK: test_channel,
    },
    Players: []agent.Agent{
      agent.Agent{
        ID: "alice",
      },
      agent.Agent{
        ID: "bob",
      },
      agent.Agent{
        ID: "zjq",
      },
    },
    TheGame: game.NewGame([]string{"game player 1", "game player 2", "game player 3"}),
    StartTime: current_time,
    winReason: "simply better than you",
    moveHistory: []string{
      "move0",
      "move2",
      "aje.2;",
    },
  }

  go func(){
    test_match.sendMatchToRecordKeeper()
  }()
  time.Sleep(20 * time.Millisecond)
  select{
  case <- test_channel:
    // not gonna compare
  default:
    t.Error("sendMatchToRecordKeeper failed to receive on the other end")
  }
}

func TestSignalEndToMM(t *testing.T){
  test_channel := make(chan string)
  test_match := Match{
    ID: "tis me",
    Channels: ChannelBundle{
      ChanMS2MM: test_channel,
    },
  }
  // test message is sent through channel
  go func(){
    test_match.signalEndToMM()
  }()
  time.Sleep(20 * time.Millisecond)
  select{
  case out := <- test_channel:
    if out != "tis me" {
      t.Error("signalEndToMM channels receivs but is wrong message" + out)
    }
  default:
    t.Error("signalEndToMM message not received by other side of channel")
  }
}

func TestNewRoundTimeStamp(t *testing.T){
  test_match_0 := Match{}
  test_match_0.newRoundTimeStamp()
  curtime := time.Now().Unix()
  if curtime != test_match_0.roundStartTime && curtime != test_match_0.roundStartTime + 1 && curtime != test_match_0.roundStartTime - 1 {
    t.Error("newRoundTimeStamp failed to set roundStartTime to current time")
  }
}

func TestTimeoutCheck(t *testing.T){
  test_match_0 := Match{
    roundStartTime: time.Now().Unix() - 100,
  }
  // returns true when it is timed out
  ret0 := test_match_0.timeoutCheck()
  if ret0 != true{
    t.Error("timeoutCheck failed to return true when it is timed out")
  }

  test_match_1 := Match{
    roundStartTime: time.Now().Unix(),
  }
  // returns false when it is not timed out
  ret1 := test_match_1.timeoutCheck()
  if ret1 != false{
    t.Error("timeoutCheck failed to return false when it is not timed out")
  }
}

func TestFindMatchByMatchID(t *testing.T){
  test_matches := []Match{
    Match{
      ID: "match 1",
    },
    Match{
      ID: "1q@W",
    },
    Match{
      ID: "match 3",
    },
  }

  // finds correct match by string
  ret1, err1 := FindMatchByMatchID(test_matches, "1q@W")
  if err1 != nil || ret1 != 1 || test_matches[ret1].ID != "1q@W"{
    t.Error("FindMatchByMatchID failed to find the correct match")
  }

  // returns error when passing it invalid ID
  ret2, err2 := FindMatchByMatchID(test_matches, "bad ID")
  if err2 == nil || ret2 != -1 {
    t.Error("FindMatchByMatchID failed to return error when passing it invalid ID")
  }
}

func TestDeleteMatchByMatchID(t *testing.T){
  test_matches0 := []Match{
    Match{
      ID: "match 1",
    },
    Match{
      ID: "1q@W",
    },
    Match{
      ID: "match 3",
    },
  }
  // deletes match when passing it valid ID
  ret0, err0 := DeleteMatchByMatchID(test_matches0, "1q@W")
  if len(ret0) != 2 || ret0[0].ID != "match 1" || ret0[1].ID != "match 3" || err0 != nil {
    t.Error("Failed to delete match correctly when passing valid ID")
  }

  test_matches1 := []Match{
    Match{
      ID: "match 1",
    },
    Match{
      ID: "1q@W",
    },
    Match{
      ID: "match 3",
    },
  }
  // returns error when passing it invalid ID
  ret1, err1 := DeleteMatchByMatchID(test_matches1, "not an ID")
  if len(ret1) != 3 || ret1[0].ID != "match 1" || ret1[1].ID != "1q@W" || ret1[2].ID != "match 3" || err1 == nil {
    t.Error("Failed to return error when passing invalid ID")
  }
}

func TestFindMatchByAgentID(t *testing.T){
  test_matches := []Match{
    Match{
      ID: "matchtest1",
      Players: []agent.Agent{
        agent.Agent{
          ID: "z5@g",
        },
        agent.Agent{
          ID: "!)fs",
        },
      },
    },
    Match{
      ID: "matchtest2",
      Players: []agent.Agent{
        agent.Agent{
          ID: "q1P",
        },
        agent.Agent{
          ID: "@)z",
        },
        agent.Agent{
          ID: "jvj2Z",
        },
      },
    },
  }

  // find existing match by agent ID
  ret_match1, err1 := FindMatchByAgentID(test_matches, "jvj2Z")
  if err1 != nil || ret_match1 != 1 || test_matches[ret_match1].ID != "matchtest2" {
    t.Error("FindMatchByAgentID failed to find existing match by ID")
  }

  // returns error if agent ID is invalid
  ret_match2, err2 := FindMatchByAgentID(test_matches, "fake name")
  if err2 == nil || ret_match2 != -1 {
    t.Error("FindMatchByAgentID failed to fail when passing it invalid agent ID")
  }
}
