package game

import (
  "testing"
)

// ========== game specific ==========
func TestDoMove(t *testing.T){
} // game specific

func TestCheckWinCondition(t *testing.T){} // game specific

// ========== game specific ==========
func TestPlayerResign(t *testing.T){
  cs := currentState {
    isResigned: false,
    resignedPlayer: 2,
    currentPlayer: 3,
    winner: 6,
    maxPlayer: 2,
  }
  cs.playerResign()
  if cs.isResigned != true || cs.resignedPlayer != 3 || cs.winner != 0 {
    t.Error("playerResign failed to modify value to expected values")
  }
}
