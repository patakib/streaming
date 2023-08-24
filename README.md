# Mock Streaming Application

## Purpose
Purpose of this application is to generate mock data for practicing streaming data pipelines.

## Architecture
The app which generates the data is written in Go.
The Dockerfile is invoked in `docker-compose.yaml` file to build the image.

The application is starting with initial data and then generating a continuous flow of data.
The data is related to workouts that users register in the system, similar to fitness trackers, but simplified.

The data is flowing into a MariaDB in two tables: user and activity.

## Quickstart

1. Create a .env file with the name of local.env and register the following secrets:

MYSQL_ROOT_PASSWORD=...
MYSQL_USER=...
MYSQL_PASS=...
MYSQL_PORT=...
MYSQL_DATABASE=...

2. Start the application and the database in Docker:  
`docker compose --env-file local.env up -d`

3. Turn down the whole infrastructure:  
`docker compose down --rmi all -v`