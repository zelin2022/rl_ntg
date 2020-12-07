package myutil

import (
  "log"
  "time"
)

func FailOnError(err error, msg string) {
  if err != nil {
    log.Printf(msg + "\n" + err.Error())
  }
}

func TimeStamp() string {
  return time.Now().Format("2006-01-02 03:04:05.999")
}

func Sleep(caller string, ms int64, ){
  log.Printf("%s will sleep for %dms.", caller, ms)
  time.Sleep(time.Duration(ms * int64(time.Millisecond))) //https://stackoverflow.com/a/42606191 how dumb
  log.Printf("%s has awakened after %dms of sleep.", caller, ms)
}

func GetCurrentEpochMilli() int64{
  return time.Now().UnixNano() / 1e6
}
