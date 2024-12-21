# CTFBattle 🎯

CTFBattle is a platform designed to simplify the process of launching and managing Capture The Flag (CTF) instances. Built with Go and the Gofr framework, CTFBattle leverages Docker containers to provide a seamless experience for hosting real-time CTF challenges. The platform’s unique feature is the ability to create private and public rooms, enabling participants to compete with friends or other users globally.

## Features 🚀

- **Easy Deployment of CTF Instances**: Quickly spin up and tear down CTF challenges using Docker containers.
- **Room-based Competition**: Create private rooms to compete with friends or join public rooms to battle against participants worldwide.
- **Global Leaderboard**: Track your progress and compare your performance with others globally. 🏆
- **Real-time Collaboration and Competition**: Experience smooth and responsive gameplay in a highly interactive environment. 💻
- **Secure Environment**: Ensures the integrity and security of the platform and CTF instances. 🔒

## Technologies Used 🛠️

- **Programming Language**: Go 🐹
- **Framework**: Gofr 📚
- **Containerization**: Docker 🐳
- **Database**: (Specify the database used, e.g., PostgreSQL, MySQL, etc.) 🗄️
- **Frontend (if applicable)**: (Specify if there is a frontend framework, e.g., React, Vue.js, etc.) 🎨
- **Other Tools**: (List any other relevant tools or libraries used in the project) 🧰

## Getting Started 🎉

### Prerequisites ✅

- **Go**: Ensure you have Go installed. [Install Go](https://golang.org/doc/install)
- **Docker**: Install Docker. [Get Docker](https://docs.docker.com/get-docker/)
- **Gofr**: Familiarity with the Gofr framework is recommended. [Gofr Documentation](https://gofr.dev/)

### Installation 🛠️

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/ctfbattle.git
   cd ctfbattle
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the application:
   ```bash
   go build -o ctfbattle
   ```

4. Run the application:
   ```bash
   ./ctfbattle
   ```

### Docker Setup 🐳

CTFBattle uses Docker to manage CTF instances. Ensure Docker is running before starting the application. The platform will automatically create and manage Docker containers for each challenge.

## Usage 📖

1. **Create a Room**:
   - Navigate to the platform’s UI or use the CLI (if applicable) to create a room.
   - Share the room code with friends or open it for public participation. 🤝

2. **Join a Room**:
   - Enter the room code to join a private room or browse available public rooms. 🌍

3. **Compete**:
   - Solve challenges, submit flags, and earn points. 🏅

4. **Leaderboard**:
   - Check your rankings and compare with other participants. 📊

## Development 🖥️

### Running in Development Mode 🛠️

Start the application in development mode:
   ```bash
   go run main.go
   ```

### Testing 🧪

Run tests to ensure the application works as expected:
```bash
   go test ./...
```

## License 📜

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments 🙌

- [Gofr Framework](https://gofr.dev/)
- [Docker](https://www.docker.com/)
- Community contributions and feedback. 💡

---

Feel free to reach out if you have any questions or suggestions for the project. Happy hacking! 💻

