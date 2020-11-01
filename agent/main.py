#!/usr/bin/python3

from myamqp.myamqp import MyAmqp
from comms.comms import Comms
import uuid

def main():
    agentID = str(uuid.uuid4())
    serverQueue = 'go_in'
    amqp = MyAmqp(serverQueue, agent_ID)
    amqp.setup(amqp_listener_callback)
    comms = Comms(agentID, amqp.queue)
    amqp.send_sign_in()

    def outside_game():

    def inside_game():



def amqp_listener_callback(ch, method, properties, body):
    print(" [x] Received %r" % body)
    msg_json = json.loads(body)
    header_to_function={
    "game start" : start_game,
    "move" : do_move,
    "game over" : end_game,
    "session interrupted" : on_session_interrupted
    "drop notification" : on_drop_notification
    }

def start_game():
    pass

def do_move():

def end_game():

def on_session_interrupted():

def on_drop_notification():




if __name__ == "__main__":
    main()
