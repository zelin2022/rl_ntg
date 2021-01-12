import hashlib

GAME_INITIAL_BOARD_NIM = 11
GAME_MOVE_RESIGN = "resign"

class Game:
    def __init__(self):
        self.my_id = None
        self.players = None
        self.board = None
        self.time_per_move = None
        self.move_count = None
        self.players_to_move = None
        self.i_just_moved = None

    def new_game(self, params, myid):
        self.my_id = myid
        self.players = params["players"]
        self.time_per_move = params["mtime"]
        self.board = GAME_INITIAL_BOARD_NIM
        self.move_count = 0
        self.players_to_move = 0
        self.i_just_moved = False

    def end_game(self):
        pass # for now

    def step(self):
        if self.my_id != self.players[self.players_to_move]: #if not my turn
            return None

        move = 1
        self.board -= move # where we get move


        move_struct = {}
        move_struct["move"] = str(move)
        move_struct["movenum"] = self.move_count

        self.move_count += 1
        self.players_to_move = (self.players_to_move + 1) % len(self.players)
        self.i_just_moved = True

        move_struct["hash"] = self.get_hash()

        return move_struct




    def update_with_others_move(self, params):
        if self.i_just_moved:
            self.i_just_moved = False
            return
        if params["movenum"] == self.move_count:
            self.apply_others_move(params["move"])
            self.move_count += 1
            self.players_to_move = (self.players_to_move + 1) % len(self.players)
            if self.get_hash() != params["hash"]:
                # potentially respond?
                raise ValueError("hash mismatch\nmy hash: %s\nmy state string: %s\ntheir hash: %s\n" % (self.get_hash(), self.get_state_string, params["hash"]))
            else:
                print("hash matches correctly, hash: %s" % params["hash"])

        else:
            raise ValueError("move_num mismatch, received %d expecting %d " % (params["movenum"], self.move_count))

    def apply_others_move(self, move):
        # specific for nim
        self.board -= int(move)

    def get_hash(self):
        val =  self.get_state_string().encode('utf-8')
        hasher = hashlib.new("sha1", val)
        return hasher.hexdigest()

    def get_state_string(self): # the output of this is used for hash
        return "%s,%d,%d" % (self.players[self.players_to_move], self.move_count, self.board)
