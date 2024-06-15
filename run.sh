#! /bin/bash

PROJECT="go_starter"

# run docker compose
docker-compose up -d db

docker-compose up -d --build ${PROJECT}
