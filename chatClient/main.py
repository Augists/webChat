import socket

HOST = "20.24.46.99"
PORT = 1201

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.connect((HOST, PORT))
    s.sendall(b'Hello, world')
    print("Sent")
