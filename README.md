# Library API

Welcome to the Library API, a Golang-based API server for managing a library.

## Overview

The Library API is designed with a clean and modular architecture, following a 3-layered structure for better maintainability and scalability.


- **Application Layer:** Orchestrates HTTP requests without containing domain logic.
- **Domain Layer:** Encompasses all domain logic, independent of infrastructure.
- **Infrastructure Layer:** Handles generic functionality, including caching.

## Getting Started

### Prerequisites

Make sure you have Docker installed on your machine.

### Running the Project

```bash
docker-compose up -d
```
## Accessing Swagger Documentation
Navigate to http://localhost:3000 in your browser to explore the Swagger documentation.

## Accessing the API Server
The API server can be accessed at http://localhost:8080.

## Postman collection
The postman collection for the api can be found at: ```build/postman/Library-api.postman_collection.json``` in this project.

# Important Notes
Allow a few moments for the server to establish a connection to the MySQL database when running Docker Compose.
The database initialization may take a minute or two. The server takes about 4 to 6 connection attempts after it starts up, 
view the logs to see the attempts and verify that connection has been established before testing.

# Architecture
The architecture promotes a modular and maintainable codebase.

## Application Layer
Coordinates HTTP requests, serving as the entry point to the API.

## Domain Layer
Contains all business logic without being aware of infrastructure details.

## Infrastructure Layer
Provides generic functionality such as database interaction and caching.

