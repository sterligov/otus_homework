#!/bin/bash

DIR="$(dirname "$0")"

docker rm calendar_migrations
docker-compose --env-file deployments/.env -f deployments/docker-compose.yml up -d --build
while [ "$(docker inspect calendar_migrations --format='{{.State.ExitCode}}')" != "0" ]; do :; done
/bin/bash ${DIR}/wait-for-it.sh localhost:5672 -t 60

docker-compose -f deployments/docker-compose.test.yml up --build
EXIT_CODE=$?

docker-compose -f deployments/docker-compose.test.yml down
docker-compose --env-file deployments/.env -f deployments/docker-compose.yml down
exit ${EXIT_CODE}
