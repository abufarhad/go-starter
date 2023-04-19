#! /bin/bash

PROJECT_NAME="go_starter"

# run docker compose
docker-compose up -d consul db

# wait for consul container be ready
while ! curl --request GET -sL --url 'http://localhost:8500/' > /dev/null 2>&1; do printf .; sleep 1; done
# setting KV, dependency of app
curl --request PUT --data-binary @config.example.yml http://localhost:8500/v1/kv/${PROJECT_NAME}

docker-compose up -d --build ${PROJECT_NAME}
