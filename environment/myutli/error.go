package myutli

import (
  "log"
  "time"
)

func FailOnError(err error, msg string) {
  if err != nil {
    log.Fatalf(TimeStamp() + "%s: %s", msg, err)
  }
}

func TimeStamp() string {
  return time.Now().Format("2006-01-02 03:04:05.000")
}
