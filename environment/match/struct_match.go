package match
import (
  "fmt"
  "time"
  "errors"
  "../game"
  "../agent"
  "../channelstructs"
  "log"
  "../myutil"
  "strconv"
)


type Match struct {
  ID string
  Channels ChannelBundle

  Players []agent.Agent
  TheGame game.Game
  StartTime time.Time
  roundStartTime int64
  moveHistory []string // for record keeper
  winReason string
}

func (m *Match) run () {
  m.matchStart()
  m.matchUnderway()
  /*

  How to determine end of match?
  * if one side resigns
  * if both sides gives up moving

  Should server actively determine end of match?
  No i guess

  */
  m.matchEnd()
}

/*
\========================================
*/

func (m *Match) matchStart(){
  m.StartTime = time.Now()
  m.roundStartTime = time.Now().Unix()
  m.TheGame = game.NewGame(agent.GetAllAgentIDs(m.Players))
  // send match start to all players
  m.broadcastStartToAllPlayers()
}

/*
\========================================
*/

func (m *Match) matchUnderway(){

  // game layer also keep track of its own version of counter
  // these are tracked by the match layer
  var moveNum uint8 = 0
  var expectedPlayer uint8 = 0

  /*
  get move from a player, then broadcast the move to all other players
  assume malicious intents from the sender though... will not process anything not from sender
  */
  for {
    select{
    case moveReceived := <- m.Channels.ChansLS2MS:
      if moveReceived.AgentID != m.Players[expectedPlayer].ID {
        log.Print("Expecting message from " + m.Players[expectedPlayer].ID + " but received a message from " + moveReceived.AgentID + " instead")
        continue
      }
      info, err := ToMatchMoveInfo(moveReceived.Body)
      if err != nil {
        myutil.FailOnError(err, "Error parsing MatchMoveInfo(json):\n" + moveReceived.Body)
        continue
      }
      err = m.doMove(moveNum, info)
      if err != nil{
        // send a response to agent
        // or end game/match directly

        // in this case do nothing and skip to next
        myutil.FailOnError(err, "doMove() failed, will continue")
        continue
      }
      // send move message to all players
      // since these games are perfect information, we can just forward a players' move to all players
      m.broadcastMoveToAllPlayers(moveReceived.Body)

      // check for a win
      if m.TheGame.CheckWinCondition(){ // if a game is over, broadcast to all agents then return
        if m.TheGame.IsResigned(){ // if  resigned, then one player forfeited, otherwise, one player reached win condition
          m.winReason = fmt.Sprintf("%s won because other player resigned", m.TheGame.GetWinner())
        }else{
          m.winReason = fmt.Sprintf("%s won by reaching winning condition", m.TheGame.GetWinner())
        }
        m.broadcastEndToAllPlayers()
        return
      }

      moveNum += 1
      expectedPlayer = (expectedPlayer + 1) % uint8(len(m.Players))
      m.newRoundTimeStamp()
    default:
      if m.timeoutCheck(){ // if a time out happens
        m.TheGame.PlayerResign() // resign the current player
        m.winReason = fmt.Sprintf("%s won because other player timed out", m.TheGame.GetWinner())
        m.broadcastEndToAllPlayers() // tell everyone game has ended
        return
      }
    }
  }
}

func (m *Match) doMove(serverMoveNum uint8, info MatchMoveInfo) error {
  /*
  * send the 2 variables ( move, statehash) to game
  * if no error occur, that means the move has been accepted
  * we can send this to other players
  */

  if info.MoveNum != serverMoveNum {
    return errors.New("Error server think move number is " + strconv.Itoa(int(serverMoveNum)) + " but received move number is " + strconv.Itoa(int(info.MoveNum)))
  }

  // forward the move to game
  moveErr := m.TheGame.TryMove(info.Move, info.StateHash)
  if moveErr != nil{
    // this is different than previous errors in this section
    // previous errors would be system errors
    // therefore more critical
    // this is soft error & game specific errors like illegal moves
    // should get a response to agent ?

    // for now if we get a invalid move, we'll not progress the game
    // implement a timeout for receiving valid move
    // and use timeout to penalize the player
    return errors.New("Error move failed. Player has made an INVALID move.")
  }
  // add move to move history for record keeping
  m.moveHistory = append(m.moveHistory, info.Move)
  return nil
}

func (m *Match) broadcastStartToAllPlayers(){
  startInfo := MatchStartInfo{
    GamePlayers: agent.GetAllAgentIDs(m.Players),
    TimePerMove: p_PLAYER_TIME_PER_MOVE,
  }
  senderMessage := channelstructs.SenderMessage {
    Header: p_HEADER_SERVER_GAME_START,
    Body: startInfo.ToString(),
  }
  senderIntake := channelstructs.SenderIntake {
    Message: senderMessage,
    AgentsToSend: m.Players,
  }
  m.Channels.ChanMS2SE <- senderIntake
}

func (m *Match) broadcastMoveToAllPlayers(body string){
  sendPackage := channelstructs.SenderIntake{
    Message: channelstructs.SenderMessage{
      Header: p_HEADER_SERVER_MOVE,
      Body: body,
    },
    AgentsToSend: m.Players,
  }
  m.Channels.ChanMS2SE <- sendPackage
}

func (m *Match) broadcastEndToAllPlayers(){
  sendPackage := channelstructs.SenderIntake{
    Message: channelstructs.SenderMessage{
      Header: p_HEADER_SERVER_GAME_END,
      Body: (&MatchEndInfo{
        Winner: m.TheGame.GetWinner(),   // potentially multiple winners
      }).ToString(),
      SendTime: time.Now().Unix(),
    },
    AgentsToSend: m.Players,
  }
  m.Channels.ChanMS2SE <- sendPackage
}




/*
\========================================
*/

func (m *Match) matchEnd(){
  //close channel with listener so listener doesn't hang trying to send
  close(m.Channels.ChansLS2MS)
  // send to record keeper and match
  m.sendMatchToRecordKeeper()
  m.signalEndToMM()
}



func (m *Match) sendMatchToRecordKeeper(){
  record := channelstructs.MatchRecord{
    Players: agent.GetAllAgentIDs(m.Players),
    StartTime: m.StartTime.Unix(),
    EndTime: time.Now().Unix(),
    Winner: m.TheGame.GetWinner(),
    WinReason: m.winReason,
    Moves: m.moveHistory,
  }

  m.Channels.ChanMS2RK <- record
}

func (m *Match) signalEndToMM(){
  m.Channels.ChanMS2MM <- m.ID // this is a many to one channel, MM will use ID to identify which match is over
}

/*
\========================================
*/

func (m *Match) newRoundTimeStamp(){
  m.roundStartTime = time.Now().Unix()
}

func (m *Match) timeoutCheck()bool{
  return (time.Now().Unix() - m.roundStartTime) > int64(p_PLAYER_TIME_PER_MOVE + p_SERVER_SIDE_TIMEOUT_GRACE)
}











func FindMatchByMatchID (matches []Match, id string)(int, error){
  for i:= range matches{
    if matches[i].ID == id{
      return i, nil
    }
  }
  return -1, errors.New("Failed to find Match by match ID:  " + id)
}

func DeleteMatchByMatchID(matches []Match, id string)([]Match, error){
  for i := range matches{
    if matches[i].ID == id {
      // swap and return
      matches[i] = matches[ len(matches)-1 ]
      return matches[ :len(matches)-1 ], nil
    }
  }
  return matches, errors.New("Failed to delete match by match id" + id)
}

// HELPR METHOD
func FindMatchByAgentID(matches []Match, agentID string)(int, error){
  for  i := 0; i < len(matches); i++{
    for j:= 0; j < len(matches[i].Players); j++{  // interate over matches and players
      if matches[i].Players[j].ID == agentID {
        return i, nil
      }
    }
  }

  return -1, errors.New("Failed to find Match by agent ID:  " + agentID)
}
