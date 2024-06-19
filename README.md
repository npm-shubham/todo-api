# Go TODO API with ScyllaDB

This repository contains a simple TODO API implemented in Go, using ScyllaDB for data storage. The API supports basic CRUD operations for managing TODO items and includes pagination functionality for listing TODO items.

## Prerequisites

Before running the application, ensure you have the following installed:

- Docker
- Go programming language (if not using Docker for development)

## Getting Started

Follow these steps to set up and run the TODO API:

### 1. Clone the Repository

```bash
git clone <repository-url>
cd todo-api
```

### 2. Verify Keyspace and Table Creation

Ensure the todo_keyspace and todos table are created in ScyllaDB. You can use cqlsh from within the container to do this.

Enter the ScyllaDB container:

```bash
docker exec -it scylla cqlsh
```

Create Keyspace and Table:

```sql
CREATE KEYSPACE IF NOT EXISTS todo_keyspace WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

CREATE TABLE IF NOT EXISTS todo_keyspace.todos (
    id UUID,
    user_id UUID,
    title TEXT,
    description TEXT,
    status TEXT,
    created TIMESTAMP,
    updated TIMESTAMP,
    PRIMARY KEY (user_id, id)
);
```

### 3. Start ScyllaDB Container

```bash
docker run --name scylla -p 9042:9042 -d scylladb/scylla
```

### 4. Build and Run the TODO API

#### Using Docker

```bash
docker build -t todo-api .
docker run --name todo-api --network scylla-net -p 8080:8080 -d todo-api
```

### 5. Verify Containers

Ensure both `scylla` and `todo-api` containers are running:

```bash
docker ps
```

### 6. Run the Application

Ensure if there is no error, and then type this in terminal.

```bash
docker start todo-api
```

### 7. Accessing the API

The API endpoints can be accessed at `http://localhost:8080/todos`.

### 8. API Endpoints

- **POST /todos**: Create a new TODO item
- **GET /todos**: List TODO items with pagination (`page` and `limit` query parameters)
- **GET /todos/{id}**: Retrieve a specific TODO item by ID
- **PUT /todos/{id}**: Update a TODO item
- **DELETE /todos/{id}**: Delete a TODO item

### 9. Postman Testing

Use Postman or any HTTP client to test the API endpoints as described above.
View Documentation [here](https://documenter.getpostman.com/view/34324389/2sA3XTeKsp)

### 10. Troubleshooting

If you encounter any issues:

- Check Docker container logs:
  ```bash
  docker logs scylla
  docker logs todo-api
  ```
- Verify network connectivity:
  ```bash
  docker network inspect scylla-net
  ```

## Additional Notes

- Make sure to configure environment variables or constants in your Go application for database credentials if needed.
- This project assumes ScyllaDB is running on default settings (localhost:9042). Modify connection settings in `db/scylla.go` if necessary.
