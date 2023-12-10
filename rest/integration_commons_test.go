package rest

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const postgresUser = "postgres"
const postgresPassword = "postgres"
const postgresDb = "postgres"
const postgresHost = "localhost"
const postgresSchema = "rest_app"
const postgresPort = "5432"
const redisPort = "6379"

var config *ITConfigurationKeys = nil

type ITConfigurationKeys struct {
	dbUrl     string
	jwtSecret string
}

func (c ITConfigurationKeys) DbConnectionUrl() string {
	return c.dbUrl
}

func (c ITConfigurationKeys) JwtSecret() string {
	return c.jwtSecret
}

func (c ITConfigurationKeys) RedisConnectionUrl() string {
	return "redis-url"
}

func createDatabaseTestContainer() testcontainers.Container {
	ctx := context.Background()
	pgRequest := testcontainers.ContainerRequest{
		Image:        "postgres:14-alpine",
		ExposedPorts: []string{postgresPort},
		WaitingFor:   wait.ForExposedPort(),
		Env: map[string]string{
			"POSTGRES_USER":     postgresUser,
			"POSTGRES_PASSWORD": postgresPassword,
			"POSTGRES_DB":       postgresDb,
		},
	}
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: pgRequest,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	port := getContainerMappedPort(&pgContainer, postgresPort)
	url := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", postgresUser, postgresPassword, port, postgresDb)

	path, err := os.Getwd()
	m, err := migrate.New(fmt.Sprintf("file://%s/../migration", path), url)
	m.Up()

	return pgContainer
}

func createRedisContainer() testcontainers.Container {
	ctx := context.Background()
	redisRequest := testcontainers.ContainerRequest{
		Image:        "redis:7.2.3",
		ExposedPorts: []string{redisPort},
		WaitingFor:   wait.ForExposedPort(),
	}
	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: redisRequest,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	return redisContainer
}

func getContainerMappedPort(container *testcontainers.Container, exposedPort nat.Port) string {
	portMapping, _ := (*container).MappedPort(context.Background(), exposedPort)
	return portMapping.Port()
}
