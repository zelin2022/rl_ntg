from myutli import myutli
import json

class Comms:
    def __init__(self, id, queue):
        self.agent_id = id
        self.agent_queue = queue

    def create_status_msg(self, status):
        out = {}
        out["Header"] = status
        out["AgentID"] = self.agent_id
        out["AgentQueue"] = self.agent_queue
        out["SendTime"] = myutli.TimeStamp()
        return json.dumps(out)
