package recordkeeper

import (
  "os"
  "myutil"
  "../channelstructs"
)


/*
right now we just save match result, in the future we can have more things potentially
*/

type RecordKeeper struct {
  file *File
  ChanMS2RK channelstructs.MatchRecord
}



func (rk *RecordKeeper)run (){
  var err error
  rk.file, err = os.Create(RECORDKEEPER_PATH_TO_RECORD + myutil.TimeStamp_RC())
  defer os.close(rk.file)
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
