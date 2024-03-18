# Movie Server

## Description

This project is a movie server that provides endpoints for managing movies and actors.

## Prerequisites

- Docker
- Docker Compose
- Make

## Getting Started

1. :Docker Compose to start the database:

   ```bash
   docker-compose up -d
   ```

2. :Run the migrations:

   ```bash
   make migrateup
   ```

3. :Run the server:

   ```bash
   go run main.go
   ```

4. :Swagger UI:

   ```bash
   make swagger
   ```
