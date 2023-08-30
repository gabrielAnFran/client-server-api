# Client-Server

This Go project consists of a simple client-server application that retrieves and stores exchange rate data for the USD to BRL (Brazilian Real) currency pair. The server retrieves the latest exchange rate from an external API and stores it in a SQLite database, while the client fetches the exchange rate from the server and saves it to a text file.

## Server

The server component (`server.go`) exposes an HTTP server that listens on port 8080. It has the following functionalities:

- **Endpoint**: The server exposes an endpoint `/cotacao` that responds to HTTP requests.

- **Retrieve Exchange Rate**: When the `/cotacao` endpoint is accessed, the server fetches the latest USD to BRL exchange rate from the "https://economia.awesomeapi.com.br" API. It then sends this data to the client.

- **Store Exchange Rate**: The fetched exchange rate is also persisted in an SQLite database named "test.db". The server utilizes the GORM library for database management.

## Client

The client component (`client.go`) fetches the exchange rate from the server using the `/cotacao` endpoint and performs the following actions:

- **Request Exchange Rate**: The client sends an HTTP GET request to the server's `/cotacao` endpoint.

- **Retrieve and Store Data**: Upon receiving the exchange rate data from the server, the client parses the response and extracts the bid price. It then creates a text file named "cotacao.txt" and writes the exchange rate information to it.

## How to Use

1. Make sure you have Go installed on your system.

2. Copy and paste the provided server code into a `.go` file (e.g., `server.go`), and the client code into another `.go` file (e.g., `client.go`).

3. Open two separate terminals, navigate to the directory containing each file, and run the following commands to start the server and client:

   - Terminal 1 (Server):
     ```
     go run server.go
     ```

   - Terminal 2 (Client):
     ```
     go run client.go
     ```

4. The client will send a request to the server and receive the exchange rate. It will then create a `cotacao.txt` file with the exchange rate information.

## Libraries Used

- `context`: Used for managing timeouts and cancellations in HTTP requests.
- `encoding/json`: Used for JSON encoding and decoding.
- `fmt`: Used for formatted printing.
- `io`: Used for reading and writing data.
- `net/http`: Used for making HTTP requests and serving HTTP on the server side.
- `os`: Used for file operations.
- `time`: Used for managing timeouts and timestamps.
- `gorm`: Used for database operations. Requires the `gorm.io/gorm` package and a SQLite database driver (`gorm.io/driver/sqlite`).

## Note

- Ensure that the server is running before executing the client.
- The project demonstrates basic client-server communication and data retrieval. For production use, consider error handling, security, and more advanced application design.
- The provided URLs and endpoints should be checked for correctness and availability. Endpoint URLs might need adjustments based on external API changes.
- The SQLite database (`test.db`) will be created in the same directory as the server executable. Make sure you have write permissions in that directory.

**Disclaimer**: This code is provided as-is, and best practices may need to be applied for production use, including securing endpoints, handling errors, and optimizing performance.
