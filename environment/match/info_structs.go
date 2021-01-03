package match

import (
  "encoding/json"
  "../myutil"
  "fmt"
)

type MatchStartInfo struct {
  GamePlayers []string `json:"players"`
  TimePerMove uint8 `json:"mtime"`
}

type MatchMoveInfo struct {
  MoveNum uint8 `json:"move"`
  Move string `json:"num"`
  AfterMoveHash string `json:"hash"`
}

type MatchEndInfo struct {
  Winner string  `json:"winner"` // potentially multiple winners
}

func (m *MatchStartInfo) ToString() (string){
  str, err := json.Marshal(m)
  myutil.FailOnError(err, "json parsing failed, struct: " + fmt.Sprintf("%v", m))
  return string(str)
}

func (m *MatchMoveInfo) ToString() (string){
  str, err := json.Marshal(m)
  myutil.FailOnError(err, "json parsing failed, struct: " + fmt.Sprintf("%v", m))
  return string(str)
}

func (m *MatchEndInfo) ToString() (string){
  str, err := json.Marshal(m)
  myutil.FailOnError(err, "json parsing failed, struct: " + fmt.Sprintf("%v", m))
  return string(str)
}

func ToMatchMoveInfo(jsonString string)(MatchMoveInfo, error){
  var outputStruct MatchMoveInfo
  err := json.Unmarshal([]byte(jsonString), &outputStruct)
  return outputStruct, err
}
