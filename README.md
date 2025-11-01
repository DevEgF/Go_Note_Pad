# Go Note Pad

A simple note-taking web application built with Go, using the MVC architecture and a MySQL database.

## Setup

### 1. Database

Create a MySQL database and a user with privileges to access it. Then, create the `notes` table using the following SQL statement:

```sql
CREATE TABLE notes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### 2. Environment Variables

Set the following environment variables to configure the database connection:

```
export DB_USER="your_database_user"
export DB_PASSWORD="your_database_password"
export DB_HOST="your_database_host"
export DB_PORT="your_database_port"
export DB_NAME="your_database_name"
```

### 3. Run the Application

```
go run main.go
```

The application will be available at `http://localhost:8080`.
