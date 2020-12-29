package game

type Game struct{
  Players []string
  state *currentState
}

func NewGame( players []string)Game{
  // create a new blank game
  newGame := Game{
    Players: players,
    state: &currentState{
      currentPlayer: 0,
      currentMoveCount: 0,
      board: GAME_INITIAL_BOARD_NIM,
      lastValidMove: "",
      maxPlayer: len(players),
      isResigned: false,
      resignedPlayer: 0,
      winner: 0,      // initial 0 should be okay... this value should only be relevant when CheckWin() returns anyway
    },
  }
  return newGame
}

func (g *Game)TryMove(move string, hash string) error{ // hash is to verify board state
  // SPECIAL CASE: if move is "resign", then set up flags and leave checkWinCondition() to pick it up
  // note: counter do not increment if resign, make sure agent hashes accordingly
  if move == GAME_MOVE_RESIGN{
    g.PlayerResign()
    return nil
  }

  /*
    first create a backup, in case currentState fails move
  */
  backup := *g.state
  err := g.cs.doMove(move)
  if err != nil{
    g.state = &backup
    return err
  }
  if g.getHash() != hash {
    g.state = &backup
    return err
  }
  return nil
}

func (g *Game)CheckWinCondition()bool{
  return g.state.checkWinCondition()
}

func (g *Game)GetWinner() string{
  return g.Players[g.state.winner]
}

func (g *Game)GetMatchEndInfo(endGameMove string) MatchEndInfo{
  output := MatchEndInfo{
    MoveNum: g.state.currentMoveCount,
    Move: g.state.lastValidMove,
    AfterEndHash: g.getHash(),
    Winner g.GetWinner(),   // potentially multiple winners
  }
  return output
}

func (g *Game)PlayerResign(){
  g.cs.playerResign()
}
