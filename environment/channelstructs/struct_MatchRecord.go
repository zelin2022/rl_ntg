package channelstructs

type MatchRecord struct {
  players []string
  winner string
}

func (mr *MatchRecord) ToString(){
  str, err = json.Marshal(mr)
  myutil.FailOnError("json parsing failed, struct: " + fmt.Sprintf("%v", mr))
  return string(str)
}
