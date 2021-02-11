
### when agent sends message to server:
```
{
  "header":
  "body":
  "aid":   #agentID
  "stime": #send time
}
```


### when server sends message to agent:
```
{
  "header":
  "body":
  "stime":  # send time
}
```


### body:

#### status
for waiting message only
```
{
  "mmc": #match making code
  "queue": # AMQP queue to reach agent
}
```

#### game start
```
{
  "players":
  "mtime": #time per move
}
```

#### game move
```
{
  "move":
  "movenum":
  'hash':
}
```
 ##### hash
 hash is used to synchronize the state of the game
 hash of a move contains the following state information AFTER a move is made:
*   player to play the next move
*   next move num
*   board after the move

 then combine them with `,` as separator, the resulting string is then hashed using:
https://stackoverflow.com/a/28094882/9520921 (OP's python and answer's go)

#### game end
```
{
  "winner":
}
```
