from myamqp.myamqp import MyAmqp
from comms.comms import Comms
import logging

class Client:
    def __init__(self, queue_out, queue_in):
        logging.info("Client __init__ begin")
        import uuid
        self.agentID = str(uuid.uuid4())
        self.amqp = MyAmqp(queue_out, self.amqp_listener_callback)
        self.amqp.setup()
        self.comms = Comms(self.agentID, self.amqp.queue)
        # self.ai = MockAI(self.amqp.ai_output, self.agentID)
        self.ingame = False
        logging.info("Client __init__ finished")

    def hook_send_sign_out(self):
        import atexit
        def exit_handler():
            self.send_sign_out()
        atexit.register(exit_handler)

    def run(self):
        logging.info("Client.run has started")
        self.send_sign_in()
        self.hook_send_sign_out()
        while True :
            if not self.ingame:
                logging.info("not in game")
                self.send_waiting()
                self.amqp.try_recv(30)
            else:
                logging.info("in game")
                pass
                # do nothing?


    def run_in_game(self):
        pass

    def amqp_listener_callback(self, ch, method, properties, body):
        print(" [x] Received %r" % body)
        loaded_msg = json.loads(body)
        header_to_function={
        "game start" : self.recv_start_game,
        "move" : self.recv_others_move,
        "game over" : self.recv_end_game,
        "session interrupted" : self.recv_on_session_interrupted,
        "drop notification" : self.recv_on_drop_notification
        }
        header_to_function[loaded_msg.header](loaded_msg)



#######################################################

    def recv_start_game(self, msg):
        self.ingame = True
        self.ai.new_game(msg.players)

    def recv_others_move(self, msg):
        self.ai.others_move(msg.move)

    def recv_end_game(self, msg):
        self.ingame = False
        self.ai.end_game()

    def recv_on_session_interrupted(self, msg):
        self.ingame = False

    def recv_on_drop_notification(self, msg):
        self.ingame = False




########################################################
    def send_sign_in(self):
        self.amqp.send_something(self.comms.create_status_msg("sign in"))

    def send_sign_out(self):
        self.amqp.send_something(self.comms.create_status_msg("sign out"))

    def send_waiting(self):
        self.amqp.send_something(self.comms.create_status_msg("waiting"))
