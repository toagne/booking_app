# Booking App – Backend #
This repository contains the backend of a simple booking application written in Go, using Gin as the web framework.
The backend exposes APIs for user authentication, match retrieval, and match booking.
The application is containerized using Docker and orchestrated with docker-compose.

🚀 Tech Stack
- Go
- Gin (HTTP web framework)
- JWT for authentication
- bcrypt for password hashing
- Docker & Docker Compose
- Thunder Client (for API testing)

📁 Project Structure
```
booking_app/
├── backend/
│   ├── db/            # Database logic and queries
│   ├── handlers/      # HTTP handlers
│   ├── middleware/    # Auth middleware (JWT)
│   ├── utils/         # Utilities (JWT, password hashing, email workers)
│   ├── main.go        # Application entry point
│   ├── Dockerfile     # Backend Docker image
│   ├── go.mod
│   └── go.sum
└── docker-compose.yml # Backend + database services
```

⚙️ Application Overview
- On startup, the application:
  - Initializes the database
  - Starts background email workers (simulated)
  - Starts the Gin HTTP server on port 8080
```go
func main() {
	db.InitDb()
	utils.StartEmailWorkers(3)

	router := gin.Default()

	router.GET("/matches/matchday/:id", handlers.GetMatchesByMatchday)
	router.GET("/matches/team/:id", handlers.GetMatchesByTeam)
	router.GET("/matches/match/:id", handlers.GetMatchByMatchId)
	router.POST("/signup", handlers.Signup)
	router.POST("/login", handlers.Login)

	auth := router.Group("/auth")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/book_match", handlers.BookMatch)
	}

	router.Run(":8080")
}
```

🔐 Authentication
- Passwords are hashed using bcrypt
- Authentication is handled using JWT
- Protected routes require a valid JWT token in the Authorization header:

🐳 Running the Application (Docker)
- Prerequisites
  - Docker
  - Docker Compose
- Start the application
- From the project root:
`docker-compose up --build`

🧪 Testing the API
- Endpoints were tested using Thunder Client (VS Code extension).
- Example login request:
```json
{
  "email": "user@example.com",
  "password": "password"
}
```

📄 Environment Variables
- Create a .env file (not committed to git) for secrets such as:
```
DB_PORT=3306
DB_NAME=db
DB_USER=user
DB_PASSWORD=secret_password
JWT_SECRET=your_secret_key
```
- An `.env.example` file can be used as a reference.

📝 Notes
- Email sending is simulated using background workers (goroutines)
- The project is intended for learning purposes and incremental improvement
- Vendor dependencies are not committed; the project relies on go.mod and go.sum

✅ Future Improvements
- Email confirmation service
- Refresh tokens
- Role-based access control
- Better error handling
- Unit and integration tests
