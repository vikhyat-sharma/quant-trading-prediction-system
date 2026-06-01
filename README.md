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

#### Stock Search & Filtering

Query parameters for `/stocks`:
- `search` - Search by stock symbol or name (case-insensitive)
- `exchange` - Filter by exchange (case-insensitive)

**Example:**
```bash
# Search for stocks matching "apple"
GET /stocks?search=apple

# Filter stocks by exchange
GET /stocks?exchange=NYSE

# Combine search and filter
GET /stocks?search=tech&exchange=NASDAQ
```

### Users

| Method | Endpoint       | Description     |
| ------ | -------------- | --------------- |
| GET    | `/users`       | List all users  |
| GET    | `/users/{id}`  | Get user by ID  |

#### User Search & Filtering

Query parameters for `/users`:
- `search` - Search by name or email (case-insensitive)

**Example:**
```bash
# Search for users matching "john"
GET /users?search=john

# Search by email
GET /users?search=john@example.com
```

### Predictions

| Method | Endpoint                                 | Description                 |
| ------ | ---------------------------------------- | --------------------------- |
| GET    | `/stocks/{stockID}/predictions`          | Get predictions for a stock |
| POST   | `/stocks/{stockID}/predictions/generate` | Generate a new prediction   |
| POST   | `/stocks/{stockID}/sentiment`            | Analyze sentiment for stock |

#### Prediction Search & Filtering

Query parameters for `/stocks/{stockID}/predictions`:
- `start_date` - Filter predictions from this date (format: YYYY-MM-DD)
- `end_date` - Filter predictions until this date (format: YYYY-MM-DD)
- `min_price` - Filter predictions with minimum price
- `max_price` - Filter predictions with maximum price

**Example:**
```bash
# Get predictions between two dates
GET /stocks/1/predictions?start_date=2024-01-01&end_date=2024-12-31

# Get predictions within a price range
GET /stocks/1/predictions?min_price=100&max_price=200

# Combine filters
GET /stocks/1/predictions?start_date=2024-06-01&min_price=150&max_price=300
```

### Price History

| Method | Endpoint                       | Description                  |
| ------ | ------------------------------ | ---------------------------- |
| GET    | `/stocks/{stockID}/price-history` | Get price history for stock |
| POST   | `/stocks/{stockID}/price-history` | Record new price            |
| GET    | `/stocks/{stockID}/price-history/range` | Get price history by date range |

#### Price History Search & Filtering

Query parameters for `/stocks/{stockID}/price-history`:
- `start_date` - Filter prices from this date (format: YYYY-MM-DD)
- `end_date` - Filter prices until this date (format: YYYY-MM-DD)
- `min_price` - Filter prices with minimum value
- `max_price` - Filter prices with maximum value

**Example:**
```bash
# Get price history for a date range
GET /stocks/1/price-history?start_date=2024-01-01&end_date=2024-12-31

# Get prices within a range
GET /stocks/1/price-history?min_price=50&max_price=150

# Get volatile prices
GET /stocks/1/price-history?start_date=2024-06-01&min_price=100&max_price=200
```


### User Watchlists

| Method | Endpoint                                                    | Description                        |
| ------ | ----------------------------------------------------------- | ---------------------------------- |
| GET    | `/users/{userID}/watchlists`                                | List all watchlists for user       |
| POST   | `/users/{userID}/watchlists`                                | Create a new watchlist             |
| DELETE | `/users/{userID}/watchlists/{watchlistID}`                  | Delete a watchlist                 |
| GET    | `/users/{userID}/watchlists/{watchlistID}/items`            | List all stocks in a watchlist     |
| POST   | `/users/{userID}/watchlists/{watchlistID}/items`            | Add a stock to a watchlist         |
| DELETE | `/users/{userID}/watchlists/{watchlistID}/items/{stockID}`  | Remove a stock from a watchlist    |

### User Alert Rules

| Method | Endpoint                                         | Description                        |
| ------ | ------------------------------------------------ | ---------------------------------- |
| GET    | `/users/{userID}/alert-rules`                    | List all alert rules for user      |
| POST   | `/users/{userID}/alert-rules`                    | Create a new alert rule            |
| DELETE | `/users/{userID}/alert-rules/{ruleID}`           | Delete an alert rule               |

### Portfolios

| Method | Endpoint                                    | Description                |
| ------ | ------------------------------------------- | -------------------------- |
| GET    | `/users/{userID}/portfolios`                | List user portfolios       |
| POST   | `/users/{userID}/portfolios`                | Create new portfolio       |
| GET    | `/users/{userID}/portfolios/{portfolioID}` | Get portfolio by ID        |
| PUT    | `/users/{userID}/portfolios/{portfolioID}` | Update portfolio           |
| DELETE | `/users/{userID}/portfolios/{portfolioID}` | Delete portfolio           |

#### Portfolio Search & Filtering

Query parameters for `/users/{userID}/portfolios`:
- `search` - Search by portfolio name (case-insensitive)

**Example:**
```bash
# Search for portfolios by name
GET /users/1/portfolios?search=conservative

# Filter specific user's portfolios
GET /users/1/portfolios?search=tech
```

### Alerts

| Method | Endpoint                                | Description               |
| ------ | --------------------------------------- | ------------------------- |
| GET    | `/stocks/{stockID}/alerts`              | Get alerts for stock      |
| POST   | `/stocks/{stockID}/alerts`              | Create alert              |
| DELETE | `/stocks/{stockID}/alerts/{alertID}`    | Delete alert              |
| POST   | `/stocks/{stockID}/alerts/evaluate`     | Evaluate all alerts       |
| GET    | `/stocks/{stockID}/notifications`       | Get notifications         |

---

## 🔧 Configuration

Environment variables:

| Variable       | Description                  | Default |
| -------------- | ---------------------------- | ------- |
| `PORT`         | Server port                  | 8080    |
| `DATABASE_URL` | PostgreSQL connection string | —       |

---

## 🔍 Filtering & Search Features

The API supports flexible filtering and searching across multiple endpoints:

### Search Syntax
- **Text Search**: Case-insensitive partial matching on specified fields
- **Date Filtering**: ISO 8601 format (YYYY-MM-DD)
- **Range Filtering**: Numeric ranges with min/max boundaries
- **Combined Filters**: All filters can be combined in a single request

### Common Filter Parameters
- `search` - Text-based search across relevant fields
- `start_date` / `end_date` - Date range filtering
- `min_price` / `max_price` - Price range filtering

All filters are optional and can be combined for more specific queries.

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
