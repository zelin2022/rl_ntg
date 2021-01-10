package recordkeeper

import (
  "os"
  "../myutil"
  "../channelstructs"
)


/*
right now we just save match result, in the future we can have more things potentially
*/

type RecordKeeper struct {
  file *os.File
  ChanMS2RK chan channelstructs.MatchRecord
}



func (rk *RecordKeeper)Run (){
  _ = os.Mkdir(p_RECORDKEEPER_PATH_TO_RECORD, 0700)
  var err error
  rk.file, err = os.Create(p_RECORDKEEPER_PATH_TO_RECORD + myutil.TimeStamp_RK())
  defer rk.file.Close()
  if err != nil{
    panic(err) // critical error
  }
  for{
    mr := <- rk.ChanMS2RK
    // now we have a MatchRecord struct, we can do all sorts of things...

    // #1 let's store the struct to file
    _, err = rk.file.WriteString(mr.ToString()+"\n")

  }
}
