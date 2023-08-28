# Mock Streaming Application

## Purpose
Purpose of this application is to generate mock data for practicing streaming data pipelines.

## Architecture
The app which generates the data is written in Go.
The Dockerfile is invoked in `docker-compose.yaml` file to build the image.

The application is starting with initial data and then generating a continuous flow of data.
The data is related to workouts that users register in the system, similar to fitness trackers, but simplified.

The data is flowing into a MariaDB in two tables: user and activity.

## Prerequisities
Docker and Docker Compose.

## Quickstart

1. Create a .env file with the name of local.env and register the following secrets:

MYSQL_ROOT_PASSWORD=...  
MYSQL_USER=...  
MYSQL_PASS=...  
MYSQL_PORT=...  
MYSQL_DATABASE=...

2. Start the application and the database in Docker:  
`docker compose --env-file local.env up -d`

3. Add new MariaDB user with privileges:
- login to the mariadb container shell:  
`docker exec -it {MARIADB_CONTAINER_ID} sh

- login to mariadb with root user:
```
mysql -uroot -p{MYSQL_ROOT_PASSWORD}
```

- add new user:
```
CREATE USER '{YOUR_NEW_USER}'@'localhost' IDENTIFIED BY '{YOUR_NEW_PASSWORD}';
```

- add privileges:
```
GRANT SELECT, RELOAD, SHOW DATABASES, REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO '{YOUR_NEW_USER}' IDENTIFIED BY '{YOUR_NEW_PASSWORD}';
```

4. Register MySQL connector in the terminal:
```
curl -i -X POST -H "Accept:application/json" -H "Content-Type:application/json" localhost:8083/connectors/ -d '{ "name": "${CONNECTOR_NAME}", "config": { "connector.class": "io.debezium.connector.mysql.MySqlConnector", "tasks.max": "1", "database.hostname": "db", "database.port": "3306", "database.user": "${YOUR_NEW_USER}", "database.password": "${YOUR_NEW_PASSWORD}", "database.server.id": "184054", "topic.prefix": "${SERVER_NAME}", "database.include.list": "streaming", "schema.history.internal.kafka.bootstrap.servers": "kafka:9092", "schema.history.internal.kafka.topic": "schemahistory.streaming" } }'
```

5. Install Spark

6. Install Python requirements:
- Go to ./spark folder  
- Create a virtual environment: `python3 -m venv venv`  
- Activate the environment: `source venv/bin/activate`  
- Install dependencies: `pip install requirements.txt`  

7. Run Spark streaming:
`python3 main.py`

8. Turn down the whole infrastructure:  
`docker compose down --rmi all -v`