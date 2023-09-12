# publicTodo

# Todo App API

## Overview

The Todo App is a basic web application that implements user authentication and authorization using JSON Web Tokens (JWT) to secure API endpoints. This project is a demonstration of how to implement user authentication and authorization in a Go web application using modern tools and best practices.

## Key Technologies & Components:

1. `Go (Golang):` The primary programming language for the backend.
2. `Gin Web Framework:` Used for building the backend APIs.
3. `Postgresql:` Database for storing user data and todos.
4. `GORM:` An ORM library in Golang for database interactions.
5. `JWT & Cookie:` For authentication and authorization.
6. `bcrypt:` For hashing passwords securely.
7. `godotenv:` For .env file interactions.

## Prerequisites

<ol><li><strong>Go Installation</strong>: Ensure that Go is <a href="https://golang.org/doc/install" target="_new">installed</a> on your machine. If you haven't installed it yet, <a href="https://golang.org/dl/" target="_new">download Go</a> and follow the installation instructions for your operating system.</li><li><strong>Postgresql</strong>: This project uses Postgres as its primary database. Make sure it's installed on your machine. If not, you can <a href="https://www.postgresql.org/download/" target="_new">download and install PostgreSQL</a>. After installation, you'll need to create a database and user for this project.</li><li><strong>Git</strong>: For cloning the repository.</li></ol>
## Getting Started

1.  Clone the repository:
    ```plaintext
    git clone https://github.com/mertkaradayi/myToDoApp
    ```
2.  Install dependencies:
    To fetch all the dependencies for the project, run it where you clone the project:
    ```plaintext
    go mod tidy
    ```
3.  Set up your .env file:
    Copy the sample .env.example to .env and fill in your configurations based on your PostgreSQL setup and desired port. If running locally, make sure to use the appropriate database connection string. If using Docker, use the provided Docker connection string.

        PORT=3000

        # #If you run local use this.
        # DB="host=/tmp user=imertkaradayi password=123456 dbname=imertkaradayi port=5433 sslmode=disable"
        #If you run it with docker, use this.
        DB="host=host.docker.internal user=imertkaradayi password=123456 dbname=imertkaradayi port=5433 sslmode=disable"

        SECRET=gowitinc

4.  Run the project:

    ```plaintext
    go run main.go
    ```

    The server will start on port 3000 (or whatever port you specify in your .env).

5.  Sending Requests:
    1. `Postman:` You can use Postman to send requests to the server. Import the Postman collection provided in the project (if any).
    2. `cURL:` Alternatively, you can use cURL from the command line.

## Docker Instructions:

For those who want to use Docker, follow these steps:

1. Build the Docker image:
   ```plaintext
   docker build -t go-todo-app .
   ```
1. Run the Docker container:
   Map the container's port 3000 to the host's port 3000:
   ```plaintext
   docker run --rm -p 3000:3000 go-todo-app
   ```

## API Documentation:

To understand the available endpoints and their functionalities, refer to the Swagger documentation at `docs/api.yaml`

## Useful Links:

<ul><li><a href="https://gin-gonic.com/" target="_new">Gin Web Framework</a></li><li><a href="https://gorm.io/" target="_new">GORM Library</a></li><li><a href="https://jwt.io/" target="_new">JSON Web Tokens (JWT)</a></li><li><a href="https://pkg.go.dev/golang.org/x/crypto/bcrypt" target="_new">bcrypt</a></li><li><a href="https://github.com/joho/godotenv" target="_new">godotenv</a></li></ul>
