# Quant Trading Prediction System

A production-ready Go backend for quantitative trading predictions, built with clean architecture principles and designed for scalability, maintainability, and performance.

---

## 🚀 Overview

This project provides a RESTful API for managing stocks and generating trading predictions. It follows a modular, layered architecture that separates concerns and simplifies testing, extension, and deployment.

---

## ✨ Key Features

* Clean Architecture (Controller → Service → Repository)
* RESTful API design
* PostgreSQL integration with auto-migrations
* Docker & Docker Compose support
* Environment-based configuration
* Testable and modular codebase

---

## 🏗 Architecture

The codebase is organized into clearly defined layers:

```
├── controllers   # HTTP handlers (API layer)
├── services      # Business logic
├── repositories  # Database interaction
├── db            # Models and DB connection
├── config        # Environment configuration
```

### Responsibilities

* **Controllers** → Handle HTTP requests/responses
* **Services** → Implement core business logic
* **Repositories** → Interact with the database
* **DB** → Define models and manage connections
* **Config** → Load and manage environment variables

---

## ⚙️ Prerequisites

* Go 1.21+
* PostgreSQL
* Docker (optional, recommended)

---

## 🛠 Setup

### Local Development

```bash
# Clone the repository
git clone https://github.com/vikhyat-sharma/quant-trading-prediction-system.git
cd quant-trading-prediction-system

# Install dependencies
go mod tidy

# Configure environment variables
cp .env.example .env
```

Update `.env` with your PostgreSQL credentials.

---

### 🗄 Database Setup

Ensure PostgreSQL is running and accessible.

> The application automatically initializes required tables and indexes on startup.

---

### ▶️ Run the Application

```bash
go run main.go
```

Server will start on:

```
http://localhost:8080
```

---

## 🐳 Docker Setup

Run the entire stack using Docker Compose:

```bash
docker-compose up --build
```

---

## 📡 API Endpoints

### Stocks

| Method | Endpoint       | Description     |
| ------ | -------------- | --------------- |
| GET    | `/stocks`      | List all stocks |
| GET    | `/stocks/{id}` | Get stock by ID |

### Predictions

| Method | Endpoint                                 | Description                 |
| ------ | ---------------------------------------- | --------------------------- |
| GET    | `/stocks/{stockID}/predictions`          | Get predictions for a stock |
| POST   | `/stocks/{stockID}/predictions/generate` | Generate a new prediction   |

---

## 🔧 Configuration

Environment variables:

| Variable       | Description                  | Default |
| -------------- | ---------------------------- | ------- |
| `PORT`         | Server port                  | 8080    |
| `DATABASE_URL` | PostgreSQL connection string | —       |

---

## 🧪 Testing

Run all tests:

```bash
go test ./...
```

Verbose output:

```bash
go test -v ./...
```

With coverage:

```bash
go test -cover ./...
```

Generate HTML coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 📦 Build

```bash
go build -o bin/app main.go
```

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Add tests where applicable
5. Open a pull request

---

## 📄 License

This project is licensed under the MIT License.

---

## 📌 Future Improvements (Optional Ideas)

* Authentication & authorization
* Real-time market data integration
* Advanced prediction models (ML/AI)
* Caching layer (Redis)
* Rate limiting & monitoring

---

## 💡 Notes

This project is intended as a clean foundation for building scalable quantitative trading systems. You can extend it with custom strategies, analytics, or machine learning pipelines as needed.
