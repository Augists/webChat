import socket
import json

HOST = "127.0.0.1"
PORT = 5556

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.connect((HOST, PORT))
    msg={
        "type": 2,
        "content": [{
                "creat_time": "21:47",
                 "name":"TEST1",
                "text": "test message"
                }]
          }
    msg = str(msg)
    s.send(str.encode(json.dumps(msg)))#发送请求

print("Sent")
