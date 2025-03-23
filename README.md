# TrendZone - NFL Stats API

TrendZone is a RESTful API service built with Go that provides access to NFL statistics through data sourced from the SportsData.io API. The application retrieves, processes, and serves NFL data including teams, players, games, schedules, and standings.

## Features

- **Comprehensive NFL Data**: Access to teams, players, games, schedules, and standings
- **Flexible Filtering**: Query data by various parameters (team, season, week, etc.)
- **MongoDB Integration**: Persistent data storage with MongoDB
- **RESTful API**: Clean, intuitive API endpoints
- **Authentication**: JWT-based authentication for protected endpoints
- **Docker Support**: Easy deployment with Docker and docker-compose
- **Swagger Documentation**: Interactive API documentation

## Technology Stack

- **Go 1.22.2**: Core programming language
- **Gin**: HTTP web framework
- **MongoDB**: Database for storing NFL data
- **SportsData.io API**: External data source for NFL statistics
- **JWT**: Authentication mechanism
- **Docker**: Containerization
- **Swagger**: API documentation

## Architecture

The application follows a clean architecture pattern with the following components:

- **API Handlers**: Process HTTP requests and responses
- **Repositories**: Data access layer for MongoDB
- **Service Layer**: Business logic and integration with SportsData.io
- **Models**: Data structures representing NFL entities
- **Middleware**: Authentication, logging, and error handling

## Project Structure

```
├── cmd/
│   └── server/              # Main application entry point
├── config/                  # Configuration handling
├── internal/                # Application internal packages
│   ├── api/                 # API related code
│   │   ├── handlers/        # Request handlers
│   │   ├── middleware/      # HTTP middleware
│   │   └── routes/          # API route definitions
│   ├── db/                  # Database related code
│   │   ├── models/          # Data models
│   │   └── mongodb/         # MongoDB specific code
│   │       └── repositories/# Data repositories
│   ├── logger/              # Logging functionality
│   └── sportsdata/          # SportsData.io API integration
├── .env                     # Environment variables
├── docker-compose.yml       # Docker compose configuration
├── Dockerfile               # Docker build instructions
└── go.mod, go.sum           # Go module files
```

## Prerequisites

- Go 1.22 or higher
- MongoDB 6.0 or higher
- SportsData.io API key
- Docker and docker-compose (optional)

## Installation and Setup

### Option 1: Using Docker (Recommended)

1. Clone the repository:
   ```
   git clone https://github.com/web-dev-jesus/trendzone.git
   cd trendzone
   ```

2. Create a `.env` file in the root directory with your SportsData.io API key:
   ```
   SPORTSDATA_API_KEY=your-api-key-here
   ```

3. Build and start the containers:
   ```
   docker-compose up -d
   ```

### Option 2: Manual Setup

1. Clone the repository:
   ```
   git clone https://github.com/web-dev-jesus/trendzone.git
   cd trendzone
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Configure environment variables in `.env` file:
   ```
   # Application
   APP_ENV=development
   APP_PORT=8080
   APP_SECRET=your-secret-key-here
   LOG_LEVEL=debug

   # MongoDB
   MONGO_URI=mongodb://localhost:27017
   MONGO_DB_NAME=sportsdata_nfl
   MONGO_TIMEOUT=10

   # SportsData.io API
   SPORTSDATA_API_KEY=your-api-key-here
   SPORTSDATA_API_BASE_URL=https://api.sportsdata.io/v3/nfl
   ```

4. Ensure MongoDB is running locally on port 27017

5. Build and run the application:
   ```
   go build -o trendzone ./cmd/server
   ./trendzone
   ```

## API Endpoints

The API provides the following endpoints:

### Public Endpoints

- `GET /health` - Health check
- `GET /swagger/*any` - Swagger documentation

### Teams Endpoints

- `GET /api/v1/teams` - Get all teams
- `GET /api/v1/teams/:id` - Get team by ID
- `GET /api/v1/teams/key/:key` - Get team by key (abbreviation)

### Players Endpoints

- `GET /api/v1/players` - Get all players
- `GET /api/v1/players/:id` - Get player by ID
- `GET /api/v1/players/pid/:playerID` - Get player by PlayerID

### Games Endpoints

- `GET /api/v1/games` - Get all games
- `GET /api/v1/games/:id` - Get game by ID
- `GET /api/v1/games/key/:gameKey` - Get game by GameKey

### Standings Endpoints

- `GET /api/v1/standings` - Get all standings
- `GET /api/v1/standings/:id` - Get standing by ID
- `GET /api/v1/standings/team/:team` - Get standing by team

### Schedules Endpoints

- `GET /api/v1/schedules` - Get all schedules
- `GET /api/v1/schedules/:id` - Get schedule by ID
- `GET /api/v1/schedules/key/:gameKey` - Get schedule by GameKey

### Protected Endpoints (require JWT authentication)

- `POST /api/v1/admin/sync` - Sync all data from SportsData.io API

## Filtering Data

Many endpoints support filtering by query parameters:

- Teams: No filters
- Players: `?team=XXX` (filter by team abbreviation)
- Games: `?team=XXX` or `?season=2023&week=1` 
- Standings: `?conference=AFC&division=East`
- Schedules: `?team=XXX` or `?season=2023&week=1`

## Authentication

Protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer {your-jwt-token}
```

## Docker Deployment

The application includes Docker support for easy deployment:

```
docker-compose up -d
```

This will start:
- The API service on port 8080
- MongoDB instance on port 27017

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.