# Money API

## Design and Architecture

The Money API application is designed as a microservice to handle banking operations such as account validation and transfer creation. It follows a clean architecture pattern, separating concerns into different layers:

- **Domain Layer**: Defines the core business logic and entities (`domain.Account` and `domain.Transfer`).
- **Repository Layer**: Handles data access and storage operations, abstracting the database interactions.
- **Service Layer**: Contains business logic that operates on the data sent to and from the Repository layer and is consumed by the Handler layer.
- **Handler Layer**: Manages HTTP request and response, acting as the interface between the HTTP server and the service layer.
- **Libraries**: Includes external libraries and utilities, such as database connection ([lib/pq/db.go](/lib/pq/db.go#1%2C1-1%2C1)) and HTTP client ([pkg/http_client/client.go](/pkg/http_client/client.go#1%2C1-1%2C1)).

## How It Works

1. **Starting the Application**: The entry point is [cmd/main.go](/cmd/main.go#1%2C1-1%2C1), where dependencies are set up, routes are configured, and the HTTP server is started.

2. **HTTP Server and Routing**: The server listens on port 8080 and routes incoming HTTP requests to the appropriate handlers for health checks, account validation, and transfer creation.

3. **Health Check**: The application provides endpoints for checking the health of the service and the database connection.

4. **Account and Transfer Operations**: The application supports validating bank accounts and creating transfers between accounts. It interacts with an external bank API for some operations.

5. **Database Interaction**: The application uses GORM for ORM and interacts with a PostgreSQL database. Database connection configuration is loaded from environment variables.

6. **Dockerization**: The application is containerized using Docker, facilitating deployment and environment consistency.

## Running the Project

1. **Environment Setup**: Ensure Docker is installed on your system.

2. **Build the Docker Image**:
   ```bash
   docker build -t money-api .
   ```
   This command builds the Docker image using the `Dockerfile`, which compiles the Go application and prepares the final image based on Alpine Linux.

3. **Run the Container**:
   ```bash
   docker run -p 8080:8080 money-api
   ```
   This command runs the application container, making the API accessible on port 8080 of the host machine.

4. **Environment Variables**: The application requires a `.env` file at the root with the following variables for database configuration:
   ```
   POSTGRES_HOST=
   POSTGRES_USER=
   POSTGRES_PASSWORD=
   POSTGRES_DB=
   POSTGRES_PORT=
   POSTGRES_SSLMODE=
   ```
   Ensure these are set before running the application.

5. **Database Migration**: On startup, the application automatically performs database migration to ensure the schema is up to date.

## Project Structure

- **cmd/**: Contains the application entry point and router setup.
- **internal/**: Houses the domain models, handlers, services, and repositories.
- **lib/**: Includes utility libraries such as the database connector.
- **pkg/**: Contains external packages like the HTTP client.
- **.gitignore**: Configures files and directories to be ignored by git.

# Concurrency in Money API Application

The Money API application leverages Go's concurrency model, primarily through the use of goroutines, to perform non-blocking operations and improve the efficiency of tasks that can be executed in parallel. Here's how concurrency is utilized in the application:

## 1. Validating Accounts Concurrently

When creating a transfer, the application needs to validate both the sender's and receiver's account numbers with an external bank API. This operation is performed concurrently using goroutines to avoid blocking the main thread while waiting for the external API responses.


```52:85:internal/repository/transfer_repository.go
func (r *transferRepository) validateAccountsWithTimeout(ctx context.Context, fromAccountNumber, toAccountNumber string) (bool, error) {
	results := make(chan bool, 2)
	errors := make(chan error, 2)

	validate := func(accountNumber string) {
		select {
		case <-ctx.Done():
			errors <- ctx.Err()
		default:
			valid, err := r.bankAPI.Validate(accountNumber)
			if err != nil {
				errors <- fmt.Errorf("error validating account %s: %w", accountNumber, err)
				return
			}
			results <- valid
		}
	}

	go validate(fromAccountNumber)
	go validate(toAccountNumber)

	validCount := 0
	for i := 0; i < 2; i++ {
		select {
		case err := <-errors:
			return false, err
		case valid := <-results:
			if valid {
				validCount++
			}
		}
	}

	return validCount == 2, nil
```


In the `validateAccountsWithTimeout` function, two goroutines are spawned to validate each account number concurrently. Channels are used for synchronization and error handling. This approach significantly reduces the validation time, especially when external API calls are involved.

## 2. Asynchronous Transfer Callback Handling

After initiating a transfer, the application handles callbacks from the bank API asynchronously. This is achieved by processing the callback in a separate goroutine, allowing the main thread to respond immediately without waiting for the callback operation to complete.


```42:49:internal/repository/transfer_repository.go
func (r *transferRepository) HandleTransferCallback(transactionID string, status string) error {
	go func() {
		err := r.bankAPI.Callback(transactionID, status)
		if err != nil {
			fmt.Printf("Error handling transfer callback: %v\n", err)
		}
	}()
	return nil
```


The `HandleTransferCallback` method uses a goroutine to call the external bank API's callback endpoint. This design ensures that the application remains responsive and can handle other requests while dealing with potentially slow external API calls.

