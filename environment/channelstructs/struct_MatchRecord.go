package channelstructs

import (
  "encoding/json"
  "../myutil"
  "fmt"
)

type MatchRecord struct {
  players []string
  winner string
}

func (mr *MatchRecord) ToString()string{
  str, err := json.Marshal(mr)
  myutil.FailOnError(err, "json parsing failed, struct: " + fmt.Sprintf("%v", mr))
  return string(str)
}
