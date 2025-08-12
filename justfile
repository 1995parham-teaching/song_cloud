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
migrate: (dev "up")
    go run main.go migrate

# run golangci-lint
lint:
    golangci-lint run -c .golangci.yml

# setup dev environment
setup: (dev "up") migrate
    #!/usr/bin/env bash
    set -euxo pipefail

    go run main.go serve &
    PID=$!

    echo "creating users..."
    curl -d '{ "username": "elahe", "first_name": "Elahe", "last_name": "Dastan", "email": "elahe.dstn@gmail.com", "password": "123456abc" }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/signup
    curl -d '{ "username": "foo", "first_name": "John", "last_name": "Doe", "email": "john.doe@gmail.com", "password": "123456abc", "introducer": "elahe" }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/signup

    echo "creating songs..."
    curl -d '{ "name": "elahe", "file": "elahe.mp3", "production_year": 2021, "explanation": "new awesome song" }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/song
    curl -d '{ "name": "elahe-p", "file": "elahe.mp3", "production_year": 2021, "explanation": "new awesome song", "price": 100 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/song

    echo "creating categories..."
    curl 127.0.0.1:8080/api/category/pop

    echo "increase user credit"
    curl -d '{ "username": "elahe", "credit": 10 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/wallet

    echo "buy some songs without credit"
    curl -d '{ "username": "elahe", "song": 2 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/buy

    echo "increase user credit"
    curl -d '{ "username": "elahe", "credit": 100 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/wallet

    echo "buy some songs with credit"
    curl -d '{ "username": "elahe", "song": 2 }' -H 'Content-Type: application/json' 127.0.0.1:8080/api/buy

    kill "$PID"
