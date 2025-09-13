# Project Joy - Currency Conversion API

A Go-based REST API for currency conversion, historical rates, forecasting, and favorite currency pair monitoring with AWS services integration.

## Features

- **Currency Conversion**: Convert amounts between different currencies
- **Historical Data**: Retrieve historical exchange rates for date ranges
- **Forecasting**: Basic probability forecasts for next-day exchange rates
- **Favorite Monitoring**: Save favorite currency pairs with threshold alerts
- **Email Notifications**: Get notified when exchange rates exceed your thresholds
- **AWS Integration**: Built with AWS services (DynamoDB, SES, SQS)

## API Endpoints

### 1. Currency Conversion
```
GET /api/v1/convert?origin={ORIGIN}&destination={DEST}&amount={AMOUNT}
```

### 2. Historical Exchange Rates
```
GET /api/v1/history?origin={ORIGIN}&destination={DEST}&start_date={YYYY-MM-DD}&end_date={YYYY-MM-DD}
```

### 3. Exchange Rate Forecast
```
GET /api/v1/forecast?origin={ORIGIN}&destination={DEST}
```

### 4. Available Destinations
```
GET /api/v1/origins/{ORIGIN}/destinations
```

### 5. Save Favorite
```
POST /api/v1/favorites
```

### 6. Check Favorites
```
POST /api/v1/favorites/check
```

### 7. Send Notification
```
POST /api/v1/notifications/email
```

## Project Structure

```
├── main.go                    # Application entry point
├── go.mod                     # Go module definition
├── domain/                    # Domain layer
│   ├── models.go             # Domain models and DTOs
│   └── services.go           # Service interfaces
├── handlers/                  # HTTP handlers layer
│   └── currency_handler.go   # REST API handlers
└── infrastructure/           # Infrastructure layer
    ├── aws_services.go       # AWS services initialization
    ├── currency_service.go   # Currency service implementation
    ├── favorite_service.go   # Favorite service implementation
    └── notification_service.go # Notification service implementation
```

## Architecture

The project follows a clean architecture pattern with three main layers:

- **Domain Layer**: Contains business models and service interfaces
- **Handlers Layer**: HTTP request/response handling and validation
- **Infrastructure Layer**: AWS services integration and external API calls

## AWS Services Used

- **DynamoDB**: Store favorites and historical data
- **SES**: Send email notifications
- **SQS**: Queue email notifications for async processing
- **Lambda**: (Future) Process queued notifications and scheduled tasks

## Getting Started

### Prerequisites

- Go 1.21 or later
- AWS CLI configured with appropriate permissions
- AWS account with access to DynamoDB, SES, and SQS

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd joy-currency-conversion-private
```

2. Install dependencies:
```bash
go mod tidy
```

3. Configure AWS credentials:
```bash
aws configure
```

4. Run the application:
```bash
go run main.go
```

The API will be available at `http://localhost:8080`

### Health Check

```bash
curl http://localhost:8080/health
```

## Configuration

### Environment Variables

- `AWS_REGION`: AWS region (default: us-east-1)
- `PORT`: Server port (default: 8080)

### AWS Resources

The following AWS resources need to be created:

1. **DynamoDB Tables**:
   - `favorites`: Store user favorite currency pairs
   - `exchange_rates`: Store current and historical exchange rates

2. **SES Configuration**:
   - Verify sender email address
   - Configure bounce and complaint handling

3. **SQS Queues**:
   - `email-notifications`: Queue for email notifications

## Development Status

This is a basic implementation with mock data. The following features need to be implemented:

- [ ] Real exchange rate API integration
- [ ] DynamoDB table creation and operations
- [ ] SES email sending implementation
- [ ] SQS message queuing
- [ ] Error handling and logging
- [ ] Input validation and sanitization
- [ ] Rate limiting and security
- [ ] Unit and integration tests
- [ ] Docker containerization
- [ ] CI/CD pipeline

## API Documentation

For detailed API documentation, see:
- `apiDesign.md`: Detailed API specification
- `joy_openapi.yaml`: OpenAPI 3.0 specification
- `project_joy_swagger.html`: Interactive API documentation

## Requirements

* I'm using the https://app.exchangerate-api.com/ API to query currencies and perform conversion, when you log in, you receive an API key that you can use when running the application
* In the currency-conversion root folder, create a .env file. Then set the `EXCHANGE_RATE_API_KEY` variable with the key you received

## Run the application
In the root folder, the Dockerfile is defined, it builds the application binary and then uses it in Alpine image.

### Dockerfile (This is not the recommended approach because the .env file is automatically loaded by Docker Compose. In addition, we will need to run more images.)
You can build the docker image as follows:
> docker build -t joy-v1 .

Then, you can run the application with the following command. The application will listen on port 8080:
> docker run -p 8080:8080 joy-v1

### Docker compose (Recommended)
To run the application locally, I recommend using `docker compose`, it loads the environment variables from the .env file and allows you to build and run all the necessary images

> docker-compose run

## License

See LICENSE file for details.