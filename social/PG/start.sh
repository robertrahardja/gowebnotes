#!/bin/bash

docker build -t socialnetwork-db .
docker run -d --name socialnetwork_db -p 5432:5432 socialnetwork-db

echo "Database started. Use: export DB_ADDR=\"postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable\""
