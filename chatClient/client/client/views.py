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

HOST = "20.24.46.99"
#HOST = "localhost"
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
class Message:
    def __init__(self,username,create_time,content):
        self.username=username
        self.create_time=create_time
        self.content =content


#主界面
@app.route('/')
def home():
    """Renders the home page."""
    return render_template('index.html')

#登录界面
@app.route('/login')
def login():
    """login page."""
    return render_template('login.html')

#登录数据处理
@app.route('/login_data',methods=['GET', 'POST'])
def login_data():
    """login data"""
    username =request.args.get("username")
    password =request.args.get("password")
    s=socket.socket(socket.AF_INET, socket.SOCK_STREAM)#创建socket
    s.connect((HOST, PORT))#连接服务器端
    msg={
    "type": 0,
    "content": [{
            "id": username,
            "password": password
             }]
    }
    msg = str(msg)
    s.send(str.encode(json.dumps(msg)))#发送请求
    while True:
        print("waiting....")
        data = s.recv(512)#接受请求
        break
    data = bytes.decode(data)
    data = eval(data)#转为字典
    message = data["content"]#获取消息
    return render_template('index.html',message = message)

#注册数据处理   
@app.route('/register_data',methods=['GET', 'POST'])
def register_data():
    """register data"""
    username =request.args.get("username")
    nickname = request.args.get("nickname")
    password =request.args.get("password")
    s=socket.socket(socket.AF_INET, socket.SOCK_STREAM)#创建socket
    s.connect((HOST, PORT))#连接服务器端
    msg={
        "type": 2,
        "content": [{
                "id": username,
                 "name":nickname,
                "password": password
                }]
          }
    msg = str(msg)
    s.send(str.encode(json.dumps(msg)))#发送请求
    while True:
        print("waiting....")
        data = s.recv(512)#接受请求
        break
    data = bytes.decode(data)
    data = eval(data)#转为字典
    message = data["content"]#获取消息
    #if message =="" 注册成功消息
    #return render_template('reg_success',message = message)
    return render_template('register.html',message = message)#注册失败返回
#返回注册界面
@app.route('/register')
def register():
    """register page."""
    return render_template('register.html')

#聊天消息测试
@app.route('/chat')
def chat():
    """chat page."""
    # get message_list from database

    message1  = Message("Niko","21:09","Hello")
    message2  = Message("Simple","21:19","Hello,too")
    message_list=[message1,message2]
    return render_template('chat.html',message_list=message_list)
