package game

/*
  this is the first iteration, in this game we play a simple Nim game
  we start with a certain number, and see who can reach 0 first
*/
type currentState struct {
  currentPlayer uint8
  currentMoveCount uint8

  board uint8      // since the board can be encoded by 1 number in this game

  // meta_data
  lastValidMove string
  maxPlayer uint8
  isResigned bool
  resignedPlayer uint8
  winner uint8 // potentially []uint8 if multiple winner
  // do not allow draw
}


func (cs *currentState)toString() string {
  // disregard meta-info
  return fmt.Sprintf("%d,%d,%d", cs.currentPlayer, cs.currentMoveCount, cs.board)
}

func (cs *currentState)doMove(move string)error{
  // for the purpsoe of sending move at end game, we record last valid move
  cs.lastValidMove = move
  // we can record any move cs receives, because if move is invalid, game will revert

  /*
  for this specific game, we just reduce board by move
  */



  // Step #1 Convert move

  // we just get uint8 value out of move string
  toReduce64, err := strconv.ParseUint(move, 10, 8)
  // 10 is the base, i.e. decimal
  // 8 is the intsize, for uint8
  // https://stackoverflow.com/a/35154920
  if err != nil {
    return err
  }
  toReduce8 = (uint8)toReduce64

  // Step #2 Validate move
  if toReduce8 > 2 || toReduce8 < 1{ // range is [1,2]
    return errors.New("Error toReduce8(move) should be [1,2] but is " + toReduce8)
  }

  // Step #3 Make move
  cs.board -= toReduce8

  // Step #4 update other things
  cs.currentMoveCount += 1
  cs.currentPlayer = (cs.currentPlayer + 1 ) % cs.maxPlayer
  return nil
}

func (cs *currentState)checkWinCondition()(bool){
  /*
  win conditions for this game (nim):
  if one player reaches 0(or lower) after its move
  if one player resigns
  */

  // case #1: resign, winner is next player
  if cs.isResigned {
    cs.winner = (cs.currentPlayer + 1) % cs.maxPlayer // the other player wins
    return true
  }

  // case #2: player who achieves 0 wins
  if cs.board <= 0 {
    // the rule is "whoever reaches 0 wins"
    // so if state is already 0, then the current player is lost, the other player wins
    cs.winner = (cs.currentPlayer + 1) % cs.maxPlayer
    return true
  }

  return false

}

func (cs *currentState)playerResign(){
  cs.isResigned = true
  cs.resignedPlayer = cs.currentPlayer
  cs.winner = (currentPlayer + 1) % cs.maxPlayer
}
