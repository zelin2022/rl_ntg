package game

import (
  "crypto/sha1"
  "encoding/hex"
  "fmt"
)








func (cs *currentState)getStateString() string {
 //  hash takes a string of 4 things:
 // *   player that made the move
 // *   move num
 // *   move
 // *   board after the move
 //
 //  then combine them with `,` as separator
  return fmt.Sprintf("%s,%d,%d", cs.players[cs.currentPlayer], cs.currentMoveCount, cs.board)
}

func (g *Game) getHash()string{ //https://stackoverflow.com/a/28094882/9520921 (OP's python and answer's go)
  data := []byte(g.state.getStateString())
  byteHash := sha1.Sum(data)
  s := hex.EncodeToString(byteHash[:])
  return s
}
