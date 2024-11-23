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
        return jsonify(response.json()), 200
    else:
        return jsonify({"error": response.json()}), response.status_code


if __name__ == '__main__':
    app.run(debug=True)