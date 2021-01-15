package channelstructs

import (
  "testing"
)

func TestMatchRecordToString(t *testing.T){
  mr := MatchRecord{
    Players: []string{"tom", "jerry", "Princess Consuela Banana-Hammock"},
    StartTime: int64(999999),
    EndTime: int64(47),
    Winner: "TSM, LOL, no.",
    WinReason: "I want a Alaskan Malamute",
    Moves: []string{"wasd", "knight to a7", "like jagger"},
  }
  expected_outcome :=
`{
    "players": [
        "tom",
        "jerry",
        "Princess Consuela Banana-Hammock"
    ],
    "start_time": 999999,
    "end_time": 47,
    "winner": "TSM, LOL, no.",
    "win_reason": "I want a Alaskan Malamute",
    "moves": [
        "wasd",
        "knight to a7",
        "like jagger"
    ]
}`
  if expected_outcome != mr.ToString(){
    t.Error("MatchRecordToString failed to produce expected string")
  }
}
