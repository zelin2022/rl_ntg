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
  return time.Now().Format("Jan _2 15:04:05.000")
}
