#!/usr/bin/python3



if __name__ == "__main__":
    import logging
    import sys
    logging.basicConfig(format='%(asctime)2s %(levelname)-2s %(message)s',stream=sys.stdout, level=logging.INFO)
    logging.info("Start")

    from client.client import Client
    pretext_id = 'x'
    if len(sys.argv) >= 2:
        pretext_id = sys.argv[1]
    client = Client(pretext_id)
    client.run()
