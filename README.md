# Tournament Service

Tournament service uses pretty_logger package from Account service.
To install pretty_logger package separately try this:
```
go get github.com/WebChads/AccountService/pkg/pretty_logger@0.0.1
```

If there are still problems with pretty_logger, you can try this:
```
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

### Running the service:
Tournament service uses SERVER_ENV environment variable to choose
in which environment (local or docker container) the service gonna run.

### Run service locally:
SERVER_ENV=local is set by default, so the export from command
line can be omitted. 
```
export SERVER_ENV=local
go run cmd/app/main.go (or go build cmd/app/main.go && ./main)
```

### Run service in the Docker container:
SERVER_ENV environment variable is set by default to "docker" in
the docker-compose file.
(see docker-compose file for more details).
```
sudo docker compose up -d --build
```

To make project work, create database "tournament_service":
1. Through docker-compose
2. Manually in your db instance
