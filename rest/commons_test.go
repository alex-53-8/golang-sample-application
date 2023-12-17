package rest

import (
	"context"
	"errors"
	"fmt"
	"os"
	"rest_app/service"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/gofrs/uuid"
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

type MockUserTokenService struct {
	CreateWasCalled   bool
	ValidateWasCalled bool
	ValidateValue     bool
	DecodeWasCalled   bool
	DecodedValueValid bool
}

func (u *MockUserTokenService) Create(userEmail string, secret string) (*string, error) {
	u.CreateWasCalled = true
	var token = "mocked-token"
	return &token, nil
}

func (u *MockUserTokenService) Validate(decodedToken *service.UserTokenDecoded) bool {
	u.ValidateWasCalled = true
	return decodedToken != nil
}

func (u *MockUserTokenService) Decode(token string, secret string) (*service.UserTokenDecoded, error) {
	u.DecodeWasCalled = true
	if u.DecodedValueValid {
		return &service.UserTokenDecoded{}, nil
	} else {
		return nil, errors.New("token is not valid")
	}
}

type MockUserService struct {
}

func (u MockUserService) FindById(userId uuid.UUID) (*service.UserInfo, error) {
	return &userInfo, nil
}
func (u MockUserService) FindByEmail(email string) (*service.UserInfo, error) {
	return &userInfo, nil
}

type MockConfigurationKeys struct {
}

func (c MockConfigurationKeys) DbConnectionUrl() string {
	return "db-url"
}

func (c MockConfigurationKeys) JwtSecret() string {
	return "jwt-secret"
}

func (c MockConfigurationKeys) RedisConnectionUrl() string {
	return "redis-url"
}

type MockCache struct {
	getWasCalled   uint
	storeWasCalled uint
}

func (rc *MockCache) Init() {
}

func (rc *MockCache) Get(key string) ([]byte, error) {
	rc.getWasCalled++
	return nil, nil
}

func (rc *MockCache) Set(key string, bytes []byte, expiration time.Duration) error {
	rc.storeWasCalled++
	return nil
}

func (rc *MockCache) Expire(key string) error {
	return nil
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
