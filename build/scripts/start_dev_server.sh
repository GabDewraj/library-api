#!/bin/bash

export MYSQL_USERNAME=root
export MYSQL_PASSWORD=password
export MYSQL_HOST=localhost
export MYSQL_PORT=3306 # mysql5.6
export MYSQL_DATABASE=library_dev
export PERFORM_MIGRATIONS=true
export MYSQL_DRIVER=mysql
export SERVER_MIGRATION_DIRECTORY="./cmd/config/migrations"
export SERVER_PORT="8080"
export SYSTEM_PARTITION="library"
export REDIS_HOST="localhost"
export REDIS_PORT="6389"
export RATE_LIMITER_MAX_REQUESTS=30
export RATE_LIMITER_WINDOW=1
go run cmd/main.go server
# Create a dump for running in compose 
# mysqldump -u root -p --host 127.0.0.1 --port 3306 --ssl-mode=REQUIRED library_dev > dump_file.sql
