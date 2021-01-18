package myutil
import (
  "log"
)
func FailOnError(err error, msg string) {
  if err != nil {
    log.Printf(msg + "\n" + err.Error())
  }
}

func PanicOnError(err error, msg string) {
  if err != nil {
    panic(err.Error() + " " + msg)
  }
}
