package channelstructs

import (
  "encoding/json"
  "../myutil"
  "fmt"
)

type MatchRecord struct {
  Players []string `json="players"`
  StartTime int64 `json="start_time"`
  EndTime int64 `json="end_time"`
  Winner string `json="winner"`
}

func (mr *MatchRecord) ToString()string{
  str, err := json.Marshal(mr)
  myutil.FailOnError(err, "json parsing failed, struct: " + fmt.Sprintf("%v", mr))
  return string(str)
}
