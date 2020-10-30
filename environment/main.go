package main

import (
  // "log"
  "./amqpmaster"
  // "encoding/json"
)

func main() {
    close := amqpmaster.Init()
    defer close() // ideally...but doesn't work for ctrl+C
    forever := make(chan bool)
    <-forever
}
