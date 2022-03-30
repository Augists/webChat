"""
This script runs the client application using a development server.
"""

from os import environ
from client import app,socketio
import _thread
import socket
from flask_socketio import SocketIO
import eventlet


def run_server():
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.bind(('localhost', 5556))
        s.listen()
        while True:
            print("Waiting")
            conn, addr = s.accept()
            with conn:
                print('Connected by', addr)
                while True:
                    data = conn.recv(1024)
                    if not data:
                        break
                # data.store()




if __name__ == '__main__':
    #_thread.start_new_thread(run_server, ())
    socketio.run(app,host="0.0.0.0")