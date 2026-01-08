# HelloBanglaTTS Backend

A Go-based backend API for Hello Bangla TTS, an AI platform for Bangla language text-to-speech services.

## Description

Hello Bangla TTS is a Bangla language-based AI platform written in Go. This backend provides RESTful APIs for user authentication, API key management, user profiles, and support ticket handling.

## Features

- User authentication with JWT tokens
- API key management for external integrations
- User profile management
- Support ticket system
- Swagger API documentation
- Docker containerization
- MongoDB database integration

## Tech Stack

- **Language**: Go 1.25.4
- **Framework**: Fiber v2
- **Database**: MongoDB
- **Authentication**: JWT
- **Documentation**: Swagger
- **Containerization**: Docker & Docker Compose

## Installation

### Prerequisites

- Go 1.25.4 or later
- Docker and Docker Compose
- MongoDB (or use the provided Docker setup)

### Clone the Repository

```bash
git clone https://github.com/abdullah-mobin/helloBanglaTTS-backend.git
cd helloBanglaTTS-backend
```

### Environment Setup

Create a `.env` file in the root directory with the following variables:

```env
APP_PORT=8080
MONGODB_URI=mongodb://localhost:27017
DATABASE_NAME=helloBanglaTTS
JWT_SECRET=your_jwt_secret
# Add other required environment variables
```

### Using Docker Compose

1. Build and run the services:

```bash
docker-compose up --build
```

This will start the MongoDB database and the Go application.

### Local Development

1. Install dependencies:

```bash
go mod download
```

2. Generate Swagger documentation:

```bash
make swag
```

3. Run the application:

```bash
make run
```

Or for development with live reload:

```bash
make dev
```

## Usage

The API server will be running on `http://localhost:8080`.

### API Endpoints

- **Authentication**: `/api/v1/auth`
  - POST `/login` - User login
  - POST `/refresh` - Refresh JWT token
  - POST `/logout` - User logout

- **Users**: `/api/v1/user`
  - User profile management endpoints

- **API Keys**: `/api/v1/api-keys`
  - Manage API keys for integrations

- **Support**: `/api/v1/support`
  - Support ticket handling

### API Documentation

Access the Swagger documentation at: `http://localhost:8080/swagger/index.html`

## Building

### Using Makefile

- `make build` - Build the application
- `make run` - Run the application
- `make dev` - Run with live reload
- `make swag` - Generate Swagger docs
- `make clean` - Clean build artifacts

### Using Docker

```bash
docker build -t hello-bangla-tts-backend .
docker run -p 8080:8080 --env-file .env hello-bangla-tts-backend
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.