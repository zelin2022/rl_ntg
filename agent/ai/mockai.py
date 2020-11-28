




class MockAI:
    def __init__(self, method, agentID):
        self.agentID = agentID
        self.players = None
        self.current_playing = None
        self.method_to_send_move = method
        self.game_playing = None

    def new_game(self, players):
        self.players = players
        self.current_playing = 0
        self.game_playing = True
        self.run_game()

    def run_game(self):
        while self.game_playing == True :
            if self.players[self.current_playing] == self.agentID    # if it's my turn, then do my move
                self.my_move()

    def others_move(self, player, move_str):
        if player != self.players[self.current_playing]
            # maybe also throw some kind of error
            raise ValueError(f"Move received, but the player who made the move shouldn't be playing right now. \nPlayer: {player} \nExpected: {self.players[self.current_playing]} at position {self.current_playing}")
            return

        # decode others' move from string
            received_move = self.string_to_move(move_str)
        # put this info into our dark magic

        self.next_player()


    def my_move(self):

        # some dark magic to come up with a move
        import random
        move = random.randint(1,2)
        ################################### for mock just random a move

        self.method_to_send_move(self.agentID, self.move_to_string(move)) #send my move to server

        self.next_player()




    def next_player(self): # game specific, but for this game it's just next in the array
        self.whose_turn += 1
        self.whose_turn = 0 if self.whose_turn == len(self.players)

    def end_game(self):
        self.game_playing = False

    @staticmethod
    def string_to_move(str): # is simple because move is just 1 number
        return int(str)

    @staticmethod
    def move_to_string(move): # is simple because move is just 1 number
        return string(move)
