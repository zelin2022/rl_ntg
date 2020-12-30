package match

import (
  "encoding/json"
  "../myutil"
  "fmt"
)

type MatchStartInfo struct {
  GamePlayers []string
  TimePerMove uint8
}

type MatchMoveInfo struct {
  MoveNum uint8
  Move string
  AfterMoveHash string
}

type MatchEndInfo struct {
  MoveNum uint8
  Move string
  AfterEndHash string
  Winner string   // potentially multiple winners
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
