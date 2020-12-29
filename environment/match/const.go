package match


const(
  PLAYER_TIME_PER_MOVE uint8 = 30,
  SERVER_SIDE_TIMEOUT_GRACE uint8 = 5, // a grace perioid for timeout

  HEADER_AGENT_MOVE string = "move",
  HEADER_SERVER_GAME_START string = "game start",
  HEADER_SERVER_GAME_END string = "game end",
  HEADER_SERVER_MOVE string = "move",
)
