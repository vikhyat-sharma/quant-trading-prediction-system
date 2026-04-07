# Quant Trading Prediction System

A clean Go backend architecture for quantitative trading predictions using REST APIs.

## Features

- Clean layered architecture (controllers, services, repositories)
- RESTful API for stock and prediction management
- PostgreSQL database integration
- Docker containerization
- Environment-based configuration

## Architecture

The project follows a clean architecture pattern with the following layers:

- **Controllers**: HTTP request handlers
- **Services**: Business logic
- **Repositories**: Data access layer
- **DB**: Database models and connections
- **Config**: Configuration management

## Prerequisites

- Go 1.21+
- PostgreSQL
- Docker (optional)

## Setup

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/vikhyat-sharma/quant-trading-prediction-system.git
   cd quant-trading-prediction-system
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up environment variables:
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

4. Set up PostgreSQL database:
   The application will automatically create the required tables and indexes on startup. Just ensure your PostgreSQL database is running and accessible.

5. Run the application:
   ```bash
   go run main.go
   ```

### Docker

1. Build and run with Docker Compose:
   ```bash
   docker-compose up --build
   ```

## API Endpoints

### Stocks

- `GET /stocks` - Get all stocks
- `GET /stocks/{id}` - Get stock by ID

### Predictions

- `GET /stocks/{stockID}/predictions` - Get predictions for a stock
- `POST /stocks/{stockID}/predictions/generate` - Generate new prediction

## Configuration

The application uses the following environment variables:

- `PORT`: Server port (default: 8080)
- `DATABASE_URL`: PostgreSQL connection string

## Development

### Running Tests

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Generate coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Building

```bash
go build -o bin/app main.go
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License