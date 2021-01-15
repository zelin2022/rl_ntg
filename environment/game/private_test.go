package game

import (
  "testing"
)

func TestGetStateString(t *testing.T){ // game specific
  cs := currentState{
    players: []string{"protoss", "john cena", "karen"},
    currentMoveCount: 33,
    board: 121,
  }
  if cs.getStateString() != "protoss,33,121" {
    t.Error("getStateString Failed to gete the correct state string")
  }
}

func TestGetHash(t *testing.T){ // game specific
  cs := currentState{
    players: []string{"protoss", "john cena", "karen"},
    currentMoveCount: 33,
    board: 121,
  }
  
  if cs.getStateString() != "protoss,33,121" {
    t.Error("getStateString Failed to gete the correct state string")
  }
}
