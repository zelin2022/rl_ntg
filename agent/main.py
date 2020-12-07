#!/usr/bin/python3


if __name__ == "__main__":
    import logging
    import sys

    logging.basicConfig(format='%(asctime)2s %(levelname)-2s %(message)s',stream=sys.stdout, level=logging.INFO)
    logging.info("Start")

    QUEUE_AGENT_2_SERVER = "server_in_0"
    QUEUE_SERVER_2_AGENT = "server_out_0"
    TIMEOUT_BETWEEN_WAITING = 5

    from client.client import Client
    client = Client(QUEUE_AGENT_2_SERVER, QUEUE_SERVER_2_AGENT, TIMEOUT_BETWEEN_WAITING)
    client.run()
