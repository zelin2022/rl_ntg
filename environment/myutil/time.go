package myutil

import (
  "log"
  "time"
)

func TimeStamp() string {
  return time.Now().Format("2006-01-02 03:04:05.000")
}

func TimeStamp_RK() string { //for recordkeeper
  return time.Now().Format("20060102030405")
}

func Sleep(caller string, ms int64, ){
  // if ms > 0 {
  if false {
    log.Printf("%s will sleep for %dms.", caller, ms)
    time.Sleep(time.Duration(ms * int64(time.Millisecond))) //https://stackoverflow.com/a/42606191 how dumb
    log.Printf("%s has awakened after %dms of sleep.", caller, ms)
  }
}

func GetCurrentEpochMilli() int64{
  return time.Now().UnixNano() / 1e6
}
