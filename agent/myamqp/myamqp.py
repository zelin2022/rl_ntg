import pika
import logging

class MyAmqp:
    HEADER_GAME_START = "game start"
    HEADER_GAME_MOVE = "move"
    HEADER_GAME_END = "game end"

    def __init__(self, target_queue, callback_method):
        self.connection = None
        self.channel = None
        self.agent_queue = None
        self.server_queue = target_queue
        self.recv_timed_out = True # add this variable so we know if process_data_events() times out or recvs
        self.callback_method = callback_method

    def setup(self):
        logging.info("MyAmqp setup begin")
        self.connection = pika.BlockingConnection(
        pika.ConnectionParameters(host='localhost', blocked_connection_timeout=3.0)) #modify this for wait time
        self.channel = self.connection.channel()

        result = self.channel.queue_declare('', exclusive=True)
        self.agent_queue = result.method.queue

        self.channel.basic_consume(queue=self.agent_queue, on_message_callback=self.amqp_listener_callback, auto_ack=True)

        logging.info("MyAmqp setup finish")

    def try_recv(self, timeout):
        logging.info("running MyAmqp.try_recv with timeout of "+ str(timeout))
        self.recv_timed_out = True
        self.channel.connection.process_data_events(time_limit=timeout)


    def send_something(self, msg):
        logging.info("Sending: " + msg)
        self.channel.basic_publish(exchange='', routing_key=self.server_queue, body=msg)

    def amqp_listener_callback(self, ch, method, properties, body):
        self.recv_timed_out = False
        self.callback_method(ch, method, properties, body)

    def get_agent_queue(self):
        return self.agent_queue
