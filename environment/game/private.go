package game

import (
  "crypto/sha1"
  "encoding/hex"
)










func (g *Game) getHash()string{ //https://stackoverflow.com/a/28094882/9520921 (OP's python and answer's go)
  data := []byte(g.currentState.toString())
  byteHash := sha1.Sum(data)
  s := hex.EncodeToString(byteHash[:])
  return s
}
