# JWT Backend Application
This is a backend development app that utilizes JSON Web Tokens (JWT) for authentication. It provides an API that allows clients to authenticate and access protected resources.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Architecture Layers](#architecture-layers)
  - [API Flow without JWT Authentication Middleware](#api-flow-without-jwt-authentication-middleware)
  - [API Flow with JWT Authentication Middleware](#api-flow-with-jwt-authentication-middleware)
- [Linting and Code Quality](#linting-and-code-quality)
  - [Linting Installation](#linting-installation)
  - [Linting Usage](#linting-usage)
- [Contributing](#contributing)
- [License](#license)

## Prerequisites

Before running this application, ensure that you have the following prerequisites installed:

- Go: [Install Go](https://go.dev/doc/install/)
- Docker: [Install Docker](https://docs.docker.com/get-docker/)
- Docker Compose: [Install Docker Compose](https://docs.docker.com/compose/install/)

## Installation

1. Clone this repository
  ```bash
    git clone https://github.com/kemalkochekov/JWT-Backend-Application.git
  ```
2. Navigate to the project directory:
  ```
    cd JWT-Backend-Application
  ```
3. Build the Docker image:
  ```
    docker-compose build
  ```

## Usage
1. Start the Docker containers:
  ```
    docker-compose up
  ```
2. The application will be accessible at:
  ```
    localhost:8080
  ```

## API Endpoints
The following API endpoints are available:
- POST http://localhost:8080/users/signup
- POST http://localhost:8080/users/login
- GET http://localhost:8080/logout
- GET http://localhost:8080/users
- GET http://localhost:8080/admin

For detailed API documentation, including examples, request/response structures, and authentication details, please refer to the

<a href="https://documenter.getpostman.com/view/31073105/2s9YeMzng6" target="_blank">
    <img alt="View API Doc Button" src="https://github.com/kemalkochekov/JWT-Backend-Development-App/assets/85355663/0c231cef-ee76-4cdf-bc41-e900845da493" width="200" height="60"/>
</a>

## Architecture Layers
I've designed a structured Go (Golang) backend architecture using Fiber, PostgreSQL, Redis, JWT auth middleware, and Docker, ensuring a robust and organized system.

![JWTbackend](https://github.com/kemalkochekov/JWT-Backend-Development-App/assets/85355663/e934493d-5568-401d-9810-0c71ffde3c43)

### API Flow without JWT Authentication Middleware
![Untitled-2023-11-28-1052](https://github.com/kemalkochekov/JWT-Backend-Development-App/assets/85355663/53ff225a-1c5c-4d4d-b06c-f21c96c968d3)

### API Flow with JWT Authentication Middleware
![with](https://github.com/kemalkochekov/JWT-Backend-Development-App/assets/85355663/bbfa0665-2c7b-45d9-aae1-d12bb87d783b)

## Linting and Code Quality

This project maintains code quality using `golangci-lint`, a fast and customizable Go linter. `golangci-lint` checks for various issues, ensures code consistency, and enforces best practices, helping maintain a clean and standardized codebase.

### Linting Installation

To install `golangci-lint`, you can use `brew`:

```bash
  brew install golangci-lint
```

### Linting Usage
1. Configuration: 

After installing golangci-lint, create or use a personal configuration file (e.g., .golangci.yml) to define specific linting rules and settings:
```bash
  golangci-lint run --config=.golangci.yml
```
This command initializes linting based on the specified configuration file.

2. Run the linter:

Once configuration is completed, you can execute the following command at the root directory of your project to run golangci-lint:

```bash
  golangci-lint run
```
This command performs linting checks on your entire project and provides a detailed report highlighting any issues or violations found.

3. Customize Linting Rules:

You can customize the linting rules by modifying the `.golangci.yml` file.

For more information on using golangci-lint, refer to the golangci-lint documentation.

## Contributing
Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request. Ensure that you follow the existing code style and conventions.

## License
This project is licensed under the [MIT License](LICENSE).
