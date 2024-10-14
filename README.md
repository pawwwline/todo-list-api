# To Do list API

This project is a RESTful API for managing user to-do lists. It supports CRUD operations (Create, Read, Update, Delete) for tasks, along with user authentication. The API allows users to register, log in, and manage their personal to-do lists.

## Key features
- User registration and login using JWT.
- CRUD operations for managing tasks (creating, retrieving, updating, and deleting tasks).
- Authentication and authorization for protecting user-specific data.

## Technologies used
- Language: Go
- Database: PostgreSQL
- Authentication: JWT
- Containerization: Docker and Docker Compose

## How to use
### Prerequisites
1. Install Dependencies: Make sure you have Docker, Docker Compose, Go v1.22 or later installed on your computer.
```bash
docker --version
docker-compose --version
go version
```

### Setup
1. Clone the Repository:

```bash
git clone https://github.com/pawwwline/todo-list-api
cd <repository-directory>
```
2. Create Configuration Files:

- Create a directory for configuration files and copy example:
```bash
mkdir config_files
cp -r config_files_example/* config_files/
```
3. Create a .env File and copy example:

``` bash
cp .env.example .env
```
Edit the .env file to set your environment variables, such as database connection strings and API keys.
### Run application

1. Build and run app container

```bash
docker-compose up --build
```

2. Apply Migrations:

- Access the running container (replace <container_name> with the name of your service):
```bash
docker exec -it <container_name> /bin/sh
```
- Run the migration command:
```bash
make migrate-up
```
- if you need to cancel migrations run:
```bash
make migrate-down
```
## API Endpoints

- **POST** `/register` - Register a new user
- **POST** `/login` - User login
- **POST** `/todos` - Create a new task
- **PUT** `/todos/{id}` - Update a task by ID
- **DELETE** `/todos/{id}` - Delete a task by ID
- **GET** `/todos` - Retrieve all tasks
