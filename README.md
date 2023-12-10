# Demo application written in Golang

## Components

The application has been written for demonstration purposes only. The application exposes REST interface on HTTP port 8081,
has a connection to PostgreSQL database. 

The application contains two endpoint:
- /user/token 
  creates a new JWT token out of provided by a client data

- /user/info
  returns user's information stored in DB, authentication is checked based on JWT token issues by the first endpoint

Source code is covered with unit and integration tests. The root folder contains 
a `docker-compose.yml` file to run all required rependencies. `.env` file contain all the environmental variables are 
currently used by the project. 


## Dependencies

- Golang version: 1.21.4
- Docker & Docker compose
- Make


## Commands

See `Makefile` for more information

### Run the application

```shell
docker-compose up
make run
```

### Run unit tests

```shell
make test-unit
```

### Run integration tests

```shell
make test-integration
```

### Run both unit & integration tests

```shell
make test-all
```

### Build

```shell
make build
```
