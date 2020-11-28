from myamqp.myamqp import MyAmqp
from comms.comms import Comms

class Client:
    def __init__(self):
        import uuid
        self.agentID = str(uuid.uuid4())
        self.serverQueue = 'go_in'
        self.amqp = MyAmqp(self.serverQueue)
        self.amqp.setup(amqp_listener_callback)
        self.comms = Comms(self.agentID, self.amqp.queue)
        self.ai = MockAI(self.amqp.ai_output, self.agentID)

    def hook_send_sign_out(self):
        import atexit
        def exit_handler():
            self.send_sign_out()
        atexit.register(exit_handler)

    def run(self):
        self.send_sign_in()
        self.hook_send_sign_out()
        while True :
            if !self.ingame:
                self.send_waiting()
                sleep(30)
            else


    def run_in_game(self):
        players

    def amqp_listener_callback(ch, method, properties, body):
        print(" [x] Received %r" % body)
        loaded_msg = json.loads(body)
        header_to_function={
        "game start" : self.recv_start_game,
        "move" : self.recv_others_move,
        "game over" : self.recv_end_game,
        "session interrupted" : self.recv_on_session_interrupted
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
        self.amqp.send_something(self.comms.generate_sign_in())

    def send_sign_out(self):
        self.amqp.send_something(self.comms.generate_sign_out())

    def send_waiting(self):
        self.amqp.send_something(self.comms.generate_waiting())
