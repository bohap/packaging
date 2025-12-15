# Server application for packing

This is a golang application that is responsible for determine the correct packs for a given number of items package.

## Structure

* `main.go` - entry point of the application
* `internal` - folder containing the main code
* `internal/appcontext` - builds the application context (DB, service, repo)
* `internal/contoller` - defines the endpoints and handlers for the app
* `internal/model` - holds the models for the app (request, response, ORM...)
* `internal/repository` - holds the repository files
* `internal/service` - holds the business logic
* `test` - folder containing the tests
* `test/test` - contains the unit tests
* `test/itest` - contains the integration tests
* `test/stub` - contains the stubs used during testing
* `openapi.yaml` - OpenAPI definition of the API

## Running
Makefile is present that can be used for building, testing and running the application

In order for the app to work, DB config needs to be provided. They are read as REQUIRED env variables:
* DB_HOST
* DB_PORT
* DB_USERNAME
* DB_PASSWORD
* DB_NAME

This can be a local DB, or the provided db-docker-compose.yml file can be used to start a docker container.
```bash
docker-compose -f db-docker-compose.yml up -d
```

`make run` will start the application on port 8080 after this. DEFAULT OS values are set for this

IF different DB config is needed, that values can be provided to `make run` command.

```bash
make run DB_HOST=localhost DB_PORT=6432 DB_USERNAME=server DB_PASSWORD=server DB_NAME=server
```

## Testing
Prerequirements: as the integration tests start a PostgreSQL container, docker is needed on the machine where tests are run.

`make test` - will run all the tests
`make testreport` - will generate a test report under build/testresults