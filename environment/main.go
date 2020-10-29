package main

import (
  "log"
  "./myamqp"
  // "encoding/json"
)

func main() {
    close := myamqp.Init(processMessage)
    defer close() // ideally...but doesn't work for ctrl+C
    forever := make(chan bool)
    <-forever
}


func processMessage(body []byte) error {
  // // 4 cases:
  // // Agent sign in
  // // Agent idle
  // // Agent sign off
  // // Agent move
  // var serverMsg ServerMsg
  // json.Unmarshal([]byte(body), &serverMsg)
  return nil
}
