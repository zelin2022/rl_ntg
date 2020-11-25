package myutil

import (
  "log"
  "time"
)

func FailOnError(err error, msg string) {
  if err != nil {
    log.Print("\n" + msg + "\n" + err.Error())
  }
}

func TimeStamp() string {
  return time.Now().Format("2006-01-02 03:04:05.000")
}

func Sleep(caller string, sec float64, ){
  log.Printf("%s will sleep for %f seconds.", caller, sec)
  time.Sleep(time.Duration(float64(time.Second) * sec)) //https://stackoverflow.com/a/42606191 how dumb
  log.Printf("%s has awakened after %fs of sleep.", caller, sec)
}
