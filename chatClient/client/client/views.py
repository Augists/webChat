"""
Routes and views for the flask application.
"""

from datetime import datetime
from flask import render_template
from client import app
from flask import request,jsonify
import pymysql
import socket
import json

# HOST = "20.24.46.99"
HOST = "localhost"
PORT = 1201



# def findData(sql):
#     # 填写必要的参数——>host：数据库地址，默认写本地localhost，user：用户名，password：密码，db：连接的数据库名，port端口：3306
#     db = pymysql.connect(host="localhost",user="root",password="123456",db="chat" ,port=3306)
#     # 创建一个游标对象
#     cur = db.cursor()
#     cur.execute(sql)
#     data = cur.fetchall()
#     db.close()
#     return data




@app.route('/')
def home():
    """Renders the home page."""
    return render_template('index.html')

@app.route('/login')
def login():
    """login page."""
    return render_template('login.html')

@app.route('/login_data',methods=['GET', 'POST'])
def login_data():
    """login data"""
    username =request.args.get("username")
    password =request.args.get("password")
    #print(username)
    #print(password)
    #print(TruePassword)
    s=socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect((HOST, PORT))
    msg={
  "type": 0,
  "content": [{
    "id": username,
    "password": password
        }]
    }
    msg = str(msg)
    print(msg)
    s.send(str.encode(json.dumps("{'type': 0, 'content': [{'id': 'abc', 'password': '123'}]}")))
    while True:
        print("waiting....")
        data = s.recv(512)
        if not data:
            break
    print(data)
    message = "登录"
    return render_template('index.html',message = message)
