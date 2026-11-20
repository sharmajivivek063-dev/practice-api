# Go Task Management API

A RESTful Task Management API built using **Golang, Gin Framework, MySQL, and JWT Authentication**. This project provides user registration, login, profile management, task CRUD operations, filtering, pagination, and task statistics.

## Features

- User registration with password hashing using bcrypt
- User login with JWT authentication
- Secure protected API routes
- User profile retrieval and updates
- Create, read, update, and delete tasks
- Task filtering by status and priority
- Search tasks by title
- Pagination support
- Task statistics dashboard
- User-specific task management
- MySQL database integration

## Tech Stack

| Technology | Purpose |
|---|---|
| Go (Golang) | Backend programming language |
| Gin | REST API framework |
| MySQL | Relational database |
| JWT | Authentication and authorization |
| bcrypt | Password hashing |
| Postman | API testing |

## Project Structure

```text
practice-api/
│
├── connection/
│   └── connection.go
│
├── controllers/
│   ├── auth_controller.go
│   ├── user_controller.go
│   └── task_controller.go
│
├── middleware/
│   └── auth_middleware.go
│
├── models/
│   ├── user.go
│   └── task.go
│
├── routes/
│   ├── auth_routes.go
│   ├── user_routes.go
│   └── task_routes.go
│
├── postman/
│   └── globals/
│       └── workspace.globals.yaml
│
├── .postman/
│   └── resources.yaml
│
├── .env
├── .gitignore
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## Installation and Setup

### 1. Clone the Repository

```bash
git clone https://github.com/sharmajivivek063-dev/practice-api.git
```

### 2. Navigate to the Project Directory

```bash
cd practice-api
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Configure Environment Variables

Create a `.env` file in the project root directory.

Add your database configuration and JWT settings using the environment variable names expected by `connection/connection.go` and `main.go`.

Example:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_mysql_password
DB_NAME=your_database_name

JWT_SECRET=your_secret_key
JWT_EXPIRY=3600
```

**Important:** Replace the example values with your actual configuration. Never upload your `.env` file or expose your database credentials and JWT secret on GitHub.

### 5. Database Setup

Create a MySQL database:

```sql
CREATE DATABASE practice_api;
USE practice_api;
```

Create the users table:

```sql
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    mobile VARCHAR(10) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP
);
```

Create the tasks table:

```sql
CREATE TABLE tasks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    title VARCHAR(150) NOT NULL,
    description TEXT NOT NULL,
    status ENUM('pending','in_progress','completed')
        NOT NULL DEFAULT 'pending',
    priority ENUM('low','medium','high') NOT NULL,
    due_date DATE NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_tasks_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE,

    INDEX idx_tasks_user_created (user_id, created_at, id),
    INDEX idx_tasks_user_status (user_id, status)
);
```

### 6. Run the Application

```bash
go run main.go
```

Server runs at:

```text
http://localhost:8080
```

## API Endpoints

Base URL:

```text
http://localhost:8080/api/v1
```

### Authentication APIs

| Method | Endpoint | Description | Authentication |
|---|---|---|---|
| POST | `/auth/register` | Register a new user | Public |
| POST | `/auth/login` | Login and generate JWT | Public |

### User APIs

| Method | Endpoint | Description |
|---|---|---|
| GET | `/users/profile` | Get user profile |
| PUT | `/users/profile` | Update user profile |

### Task APIs

| Method | Endpoint | Description |
|---|---|---|
| POST | `/tasks` | Create a task |
| GET | `/tasks` | Get all user tasks |
| GET | `/tasks/:id` | Get task by ID |
| PUT | `/tasks/:id` | Update task |
| DELETE | `/tasks/:id` | Delete task |
| GET | `/tasks/stats` | Get task statistics |

All user and task endpoints require JWT authentication.

## Authentication

This project uses JSON Web Tokens (JWT) to protect private routes.

### 1. Register User

**POST** `/auth/register`

Request body:

```json
{
    "name": "John Smith",
    "email": "john@example.com",
    "password": "john@12345",
    "mobile": "9876543210"
}
```

### 2. Login User

**POST** `/auth/login`

Request body:

```json
{
    "email": "john@example.com",
    "password": "john@12345"
}
```

Copy the `access_token` received in the response.

### 3. Pass JWT Token

For protected endpoints, add this HTTP header:

```http
Authorization: Bearer YOUR_JWT_TOKEN
```

## Task Management

### Create Task

**POST** `/tasks`

```json
{
    "title": "Develop Login API",
    "description": "Implement JWT authentication using Gin",
    "priority": "high",
    "due_date": "2026-11-25"
}
```

### Filter Tasks

The task listing endpoint supports pagination and filters.

Examples:

```text
GET /tasks?page=1&limit=10
```

```text
GET /tasks?status=pending
```

```text
GET /tasks?priority=high
```

```text
GET /tasks?search=login
```

Combined filters:

```text
GET /tasks?page=1&limit=10&status=pending&priority=high&search=login
```

### Task Statistics

**GET** `/tasks/stats`

Provides:

- Total tasks
- Pending tasks
- In-progress tasks
- Completed tasks
- High-priority tasks
- Overdue tasks

## Task Status and Priority

**Task Status:**

- `pending`
- `in_progress`
- `completed`

**Task Priority:**

- `low`
- `medium`
- `high`

## Testing

API endpoints can be tested using [Postman](https://www.postman.com/).

1. Start the Go server.
2. Register a user.
3. Log in to obtain a JWT token.
4. Set the Bearer Token in Postman's Authorization tab.
5. Test user and task endpoints.

## Security

- Passwords are stored as bcrypt hashes.
- JWT authentication protects private routes.
- Tasks are associated with individual users.
- Users can access and manage only their own tasks.
- Environment variables are used for sensitive configuration.

## Future Improvements

- Unit and integration testing
- Role-based access control
- Docker support
- API documentation using Swagger
- Task reminders and notifications
- Deployment to a cloud server

## Author

**Vivek Sharma**

GitHub: [sharmajivivek063-dev](https://github.com/sharmajivivek063-dev)

Project Repository: [Go Task Management API](https://github.com/sharmajivivek063-dev/practice-api)

---

If you find this project useful, consider giving the repository a ⭐ on GitHub.
