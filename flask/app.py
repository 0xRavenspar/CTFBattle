from flask import Flask, request, jsonify
from dotenv import load_dotenv
import requests
import secrets
import string
import os

load_dotenv()

app = Flask(__name__)

base_url = "http://localhost:1337/api/v1"
key = os.getenv('API_KEY')
headers = {
        "Authorization": f"Bearer {key}",
        "Content-Type": "application/json"
}

@app.route('/')
def index():
    return "Hello, CTFBattle!"

@app.route('/create_user', methods=['POST'])
def create_user():
    user_data = request.json

    characters = string.ascii_letters + string.digits + string.punctuation
    password = ''.join(secrets.choice(characters) for _ in range(8))

    payload = {
        "name": user_data["name"],
        "email": user_data["email"],
        "password": password,
        "type": "user",
        "verified": False,
        "hidden": False,
        "banned": False,
        "fields": []
    }

    url = f"{base_url}/users"

    response = requests.post(url, headers=headers, json=payload)

    if response.status_code == 200:
        return password
    else:
        return jsonify({"error": response.json()}), response.status_code

@app.route('/user_score', methods=['POST'])
def user_score():
    user_data = request.json

    name = user_data["name"]

    url = f"{base_url}/scoreboard"

    response = requests.get(url, headers=headers)
    response = response.json()["data"]

    user_info = next((entry for entry in response if entry["name"] == name), None)
    
    if user_info:
        return {"position": user_info["pos"], "score": user_info["score"]}
    return {"error": "User not found"}

@app.route('/room_info', methods=['GET'])
def room_info():
    url = f"{base_url}/statistics/users"

    response = requests.get(url, headers=headers)
    response = response.json()["data"]

    no_of_users = int(response["registered"])-1

    url = f"{base_url}/challenges"

    response = requests.get(url, headers=headers)
    response = response.json()["data"]
    
    no_of_challenges = len(response)

    url = f"{base_url}/scoreboard"

    response = requests.get(url, headers=headers)
    response = response.json()["data"]

    highest_score = response[0]["score"]

    return {"no_of_users": no_of_users, "no_of_challenges": no_of_challenges, "highest_score": highest_score}


if __name__ == '__main__':
    app.run(debug=True)