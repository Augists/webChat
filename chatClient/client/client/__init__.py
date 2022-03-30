"""
The flask application package.
"""

from flask import Flask
from flask_socketio import SocketIO
import eventlet
app = Flask(__name__)
socketio = SocketIO(app,async_mode='eventlet')
#socketio = SocketIO(app, async_mode='threading')
import client.views
