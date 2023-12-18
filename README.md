# Demo application written in Golang
It is a demo application and nothing more - all similarities with real life are coincidental 😎

The main purpose of the application is to use common software development practices in Golang and run & debug the Golang application with VS Code.

Main things to try:
- create REST application
- try Test Driven Development with Golang
- implement mocks
- add integration tests as part of Behavior Driven Development.

## Application overview

The application exposes REST interface on HTTP port 8081,has connects to PostgreSQL database and Redis. 

Two REST endpoint:
- `/user/token`
  creates a new JWT token out of provided by a client's email. See DB migration for default data.

- `/user/info`
  returns user's information stored in DB, authentication is checked based on JWT token issues by the first endpoint

Source code is covered with unit and integration tests. The root folder contains 
a `docker-compose.yml` file with all required dependencies as well as database migration. `.env` file contain all the environmental variables are 
currently used by the project. 

### Dependencies

- Golang version: 1.21.4 or higher
- Docker & Docker compose for running dependencies and integration tests
- Make for building and running tests in a command line

### Make commands

See `Makefile` for more information

#### Run the application

```shell
docker-compose up
make run
```

#### Run unit tests

```shell
make test-unit
```

#### Run integration tests

```shell
make test-integration
```

#### Run both unit & integration tests

```shell
make test-all
```

#### Build

```shell
make build
```
