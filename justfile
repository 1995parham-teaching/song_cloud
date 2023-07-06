default:
    @just --list

# build song_cloud binary
build:
    go build -o song_cloud

# update go packages
update:
    go get -u

# set up the dev environment with docker-compose
dev cmd *flags:
    #!/usr/bin/env bash
    set -euxo pipefail
    if [ {{ cmd }} = 'down' ]; then
      docker compose -f docker-compose.yml down
      docker compose -f docker-compose.yml rm
    elif [ {{ cmd }} = 'up' ]; then
      docker compose -f docker-compose.yml up -d {{ flags }}
    else
      docker compose -f docker-compose.yml {{ cmd }} {{ flags }}
    fi

# connect into the dev environment database
database: (dev "up") (dev "exec" "database psql -U song")

# run migration in the dev environment
migrate:
    go run main.go migrate

# run golangci-lint
lint:
    golangci-lint run -c .golangci.yml
