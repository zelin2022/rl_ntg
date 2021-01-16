package match

import (
  "testing"
)

func TestMatchStartInfoToString(t *testing.T){
  msi := MatchStartInfo{
    GamePlayers: []string{"faker", "November", "dUCk", "/sd2."},
    TimePerMove: uint8(32),
  }
  if msi.ToString() != `{"players":["faker","November","dUCk","/sd2."],"mtime":32}`{
    t.Error("MatchStartInfo ToString failed to return correct string")
  }
}

func TestMatchMoveInfoToString(t *testing.T){
  mmi := MatchMoveInfo{
    Move: "take a shower, then @ eat some peanuts",
    MoveNum: uint8(32),
    StateHash: "LAM!1jdmcvoi2ml;.;v'",
  }
  if mmi.ToString() != `{"move":"take a shower, then @ eat some peanuts","movenum":32,"hash":"LAM!1jdmcvoi2ml;.;v'"}` {
    t.Error("MatchMoveInfo ToString failed to return correct string")
  }
}

func TestMatchEndInfoToString(t *testing.T){
  mei := MatchEndInfo{
    Winner: "the best !",
  }
  if mei.ToString() != `{"winner":"the best !"}`{
    t.Error("MatchEndInfo ToString failed to return correct string")
  }
}

func TestToMatchMoveInfo(t *testing.T){
  js := `{"move":"take a shower, then @ eat some peanuts","movenum":32,"hash":"LAM!1jdmcvoi2ml;.;v'"}`
  expectedStruct := MatchMoveInfo{
    Move: "take a shower, then @ eat some peanuts",
    MoveNum: uint8(32),
    StateHash: "LAM!1jdmcvoi2ml;.;v'",
  }
  out, err := ToMatchMoveInfo(js)
  if out != expectedStruct || err != nil {
    t.Error("ToMatchMoveInfo failed to unmarshal json string")
  }
}
