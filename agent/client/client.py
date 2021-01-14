from myamqp.myamqp import MyAmqp
from game.game import Game
from datetime import datetime
import json
import logging
import time
# when agent sends message to server:
# {
#   header:
#   body:
#   agentID:
#   sendtime:
# }
#
#
#
# when server sends message to agent:
# {
#   header:
#   body:
#   sendtime:
# }
QUEUE_AGENT_2_SERVER = "server_in_0"
LISTENING_TIMEOUT = 5
HEADER_AGENT_2_SERVER_MOVE = "move"
HEADER_SERVER_2_AGENT_START = "game start"
HEADER_SERVER_2_AGENT_MOVE = "move"
HEADER_SERVER_2_AGENT_END = "game end"

class Client:
    def __init__(self):
        import uuid
        self.agentID = str(uuid.uuid4())
        self.amqp = MyAmqp(QUEUE_AGENT_2_SERVER, self.amqp_listener_callback)
        self.amqp.setup()
        self.ingame = False
        self.next_waiting_send_time = 0
        self.game = Game()

    def run(self):
        logging.info("Client.run has started")
        self.send_sign_in()
        self.hook_send_sign_out()
        while True:
            self.amqp.try_recv(5)
            if not self.ingame:
                logging.info("not in game")
                if self.should_send_waiting():
                    self.send_waiting()
            else:
                logging.info("in game")
                to_server = self.game.step()
                if to_server:
                    self.amqp.send_something(self.create_move_msg(to_server))

    def should_send_waiting(self):
        current_time = int(time.time())
        if current_time > self.next_waiting_send_time:
            self.next_waiting_send_time = current_time + LISTENING_TIMEOUT
            return True
        else:
            return False

    def reset_next_waiting_send_time(self): # when game start we call this, so agent sends watiing as soon as game is over
        self.next_waiting_send_time = 0

#######################################################


# when agent sends message to server:
# {
#   header:
#   body:
#   agentID:
#   sendtime:
# }

    def create_move_msg(self, body):
        output = {}
        output["header"] = HEADER_AGENT_2_SERVER_MOVE
        output["body"] = json.dumps(body)
        output["aid"] = self.agentID
        output["stime"] = self.time_stamp()
        return json.dumps(output)

    def create_status_msg(self, status):
        output = {}
        output["header"] = status
        output["body"] = self.amqp.get_agent_queue()
        output["aid"] = self.agentID
        output["stime"] = self.time_stamp()
        return json.dumps(output)

    @staticmethod
    def time_stamp():
        return datetime.utcnow().strftime('%Y-%m-%d %H:%M:%S.%f')[:-3]



#######################################################

# when server sends message to agent:
# {
#   header:
#   body:
#   sendtime:
# }

    def amqp_listener_callback(self, ch, method, properties, body):
        print(" [x] Received %r" % body)
        loaded_msg = json.loads(body)
        print(loaded_msg)
        header_to_function={
        HEADER_SERVER_2_AGENT_START : self.recv_start_game,
        HEADER_SERVER_2_AGENT_MOVE : self.recv_others_move,
        HEADER_SERVER_2_AGENT_END : self.recv_end_game,
        }
        loaded_body = json.loads(loaded_msg["body"])
        print(loaded_body)
        header_to_function[loaded_msg["header"]](loaded_body)

#######################################################

    def recv_start_game(self, msg):
        self.ingame = True
        self.game.new_game(msg, self.agentID)
        self.reset_next_waiting_send_time()

    def recv_others_move(self, msg):
        self.game.update_with_others_move(msg)

    def recv_end_game(self, msg):
        self.ingame = False
        self.game.end_game()

#######################################################



    def hook_send_sign_out(self):
        import atexit
        def exit_handler():
            self.send_sign_out()
        atexit.register(exit_handler)

########################################################
    def send_sign_in(self):
        self.amqp.send_something(self.create_status_msg("sign in"))

    def send_sign_out(self):
        self.amqp.send_something(self.create_status_msg("sign out"))

    def send_waiting(self):
        self.amqp.send_something(self.create_status_msg("waiting"))
