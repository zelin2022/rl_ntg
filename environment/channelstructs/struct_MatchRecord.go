package channelstructs

import (
  "encoding/json"
  "../myutil"
  "fmt"
)

type MatchRecord struct {
  Players []string `json:"players"`
  StartTime int64 `json:"start_time"`
  EndTime int64 `json:"end_time"`
  Winner string `json:"winner"`
  WinReason string `json:"win_reason"`
  Moves []string `json:"moves"`
}

func (mr *MatchRecord) ToString()string{
  str, err := json.MarshalIndent(mr, "", "    ")
  myutil.FailOnError(err, "MatchRecord.ToString() error: json parsing failed, struct: " + fmt.Sprintf("%v", mr))
  return string(str)
}
