# Event Connect Application

This application is built using Go and PostgreSQL, and it is containerized using Docker.

## Prerequisites

- Docker (https://www.docker.com/get-started)
- Docker Compose (usually bundled with Docker)

## Getting Started

1. Make sure you have the application files in a directory on your machine.

2. Open a terminal and navigate to the directory containing the application files.

3. Build and run the application using Docker Compose:

   docker-compose up --build

   This command will build the Docker image for the application and start the containers defined in the `docker-compose.yml` file.

4. Access the application in your web browser at `http://localhost:8000`.

5. To stop the application and the database containers, press `Ctrl+C` in the terminal where you ran `docker-compose up`.

## Configuration

The application uses environment variables for configuration. The following variables are defined in the `docker-compose.yml` file:

- `DB_HOST`: The hostname of the PostgreSQL database container (default: `db`).
- `DB_PORT`: The port number of the PostgreSQL database (default: `5432`).
- `DB_USER`: The username for connecting to the PostgreSQL database (default: `postgres`).
- `DB_PASSWORD`: The password for connecting to the PostgreSQL database (default: `admin`).
- `DB_NAME`: The name of the PostgreSQL database (default: `postgres`).

You can modify these variables in the `docker-compose.yml` file if needed.

## Database Initialization

The application uses a database initializer to create the necessary tables and schemas. The database initialization is handled automatically when the application container starts.

The database initializer is defined in the `db_initializer.go` file located in the `models` package.


## Contact

If you have any questions or issues, please contact Maurice Jarvis at 16043988@stu.mmu.ac.uk 






