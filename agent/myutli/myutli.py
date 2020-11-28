from datetime import datetime
def TimeStamp():
    return datetime.utcnow().strftime('%Y-%m-%d %H:%M:%S.%f')[:-3]
