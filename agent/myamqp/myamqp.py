import pika

class MyAmqp:

    def __init__(self, target_queue):
        self.connection = None
        self.channel = None
        self.queue = None
        self.target_queue = target_queue

    def setup(self, consume_callback):
        self.connection = pika.BlockingConnection(
        pika.ConnectionParameters(host='localhost'))
        self.channel = self.connection.channel()

        result = self.channel.queue_declare('', exclusive=True)
        self.queue = result.method.queue

        channel.basic_consume(queue=self.queue, on_message_callback=consume_callback, auto_ack=True)
        self.channel.start_consuming()

    def send_something(self, msg):
        self.channel.basic_publish(exchange='', routing_key=self.target_queue, body=msg)

    def ai_output(self, agentID, move):
        import json
        msg = {}
        msg["Header"] = "move"
        msg["AgentID"] = agentID
        msg["AgentQueue"] = self.queue
        msg["Move"] = move
        msg["SendTime"] = TimeStamp()
        self.send_something(json.dumps(msg))
