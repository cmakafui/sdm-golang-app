# Sparse Distributed Memory (SDM) Application

This project is an implementation of Sparse Distributed Memory (SDM) using Go for the backend and a simple HTML/CSS/JavaScript interface enhanced with HTMX and TailwindCSS for the frontend. The application allows storing and retrieving high-dimensional binary vectors and visually highlights differences between stored and retrieved data.

## Table of Contents

- [Description](#description)
- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Architecture](#architecture)
- [API Endpoints](#api-endpoints)
- [Troubleshooting](#troubleshooting)
- [License](#license)

## Description

Sparse Distributed Memory (SDM) is a model of long-term memory storage inspired by how the human brain may store information. This application demonstrates how to store and retrieve data using SDM principles, including random binary vectors, Hamming distances, and convergence.

## Features

- Store binary data using a randomly generated or user-provided address.
- Retrieve stored binary data based on the provided address.
- Visual comparison between stored and retrieved data.
- Simple and intuitive web interface with TailwindCSS for modern styling.
- Copy retrieved data to the clipboard.
- Clear memory functionality.
- Logging for debugging and monitoring operations.
- Responsive design for a seamless experience across devices.

## Requirements

- Go 1.16 or later
- A modern web browser

## Installation

1. **Clone the repository:**

   ```sh
   git clone https://github.com/yourusername/sdm-golang-app.git
   cd sdm-golang-app
   ```

2. **Install Go dependencies:**

   ```sh
   go mod tidy
   ```

3. **Build the project:**
   ```sh
   go build -o sdm-app cmd/main.go
   ```

## Usage

1. **Run the application:**

   ```sh
   ./sdm-app
   ```

2. **Open your browser and navigate to:**

   ```
   http://localhost:5080
   ```

3. **Use the web interface to generate test data or input your own address and data to store and retrieve information. The output will be displayed below the form, with differences between stored and retrieved data highlighted.**

## Architecture

### Backend (Go)

- **SDM Module:** Contains the core logic for storing and retrieving data using SDM principles.
- **Main Module:** Sets up HTTP server and handles routing.

### Frontend (HTML/CSS/JavaScript)

- **HTML Template:** Provides the structure of the web interface.
- **CSS (TailwindCSS):** Provides modern and responsive styling for the web interface.
- **JavaScript (HTMX):** Handles dynamic interactions and data generation.

## API Endpoints

- **GET /**: Serves the main web interface.
- **POST /**: Handles storing and retrieving data.
- **GET /generate-random**: Generates random binary address and data.
- **POST /clear-memory**: Clears the memory.
- **GET /memory-stats**: Retrieves memory statistics and history of stored addresses.

### Example Payload

```json
{
  "address": "1000101010... (1000 bits)",
  "data": "1101010101... (1000 bits)"
}
```

## Troubleshooting

If you encounter issues running the application, ensure that all dependencies are installed correctly and that you are using a compatible version of Go. Check the logs for detailed error messages.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
