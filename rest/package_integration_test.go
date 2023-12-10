package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"rest_app/cache"
	"rest_app/database"
	"rest_app/service"
	"testing"

	"github.com/cucumber/godog"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const userId = "2f3d92a9-f4db-497b-bf16-367bf2ad7e20"
const userEmail = "boss@company.com"

func TestIntegration(t *testing.T) {
	t.Run("test/integration", func(t *testing.T) {
		pgContainer := createDatabaseTestContainer()
		defer func() {
			if err := pgContainer.Terminate(context.Background()); err != nil {
				panic(err)
			}
		}()

		port := getContainerMappedPort(&pgContainer, postgresPort)
		dbUrl := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?search_path=%s", postgresUser, postgresPassword, port, postgresDb, postgresSchema)

		var db = database.DatabaseService{DbConnectionUrl: dbUrl}
		db.Init()

		redisContainer := createRedisContainer()
		defer func() {
			if err := redisContainer.Terminate(context.Background()); err != nil {
				panic(err)
			}
		}()

		redisMappedPort := getContainerMappedPort(&redisContainer, redisPort)
		redisUrl := fmt.Sprintf("redis://localhost:%s/0", redisMappedPort)
		var _cache = cache.RedisCache{RedisConnectionUrl: redisUrl}
		_cache.Init()

		config = &ITConfigurationKeys{dbUrl, "secret"}
		restServer := CreateRestServer(&_cache, &db, config)
		go func() {
			restServer.Start()
		}()
		defer restServer.Stop()

		suite := godog.TestSuite{
			ScenarioInitializer: InitializeScenario,
			Options: &godog.Options{
				Format:   "pretty",
				Paths:    []string{"features"},
				TestingT: t,
			},
		}

		if suite.Run() != 0 {
			t.Fatal("non-zero status returned, failed to run feature tests")
		}
	})
}

type userIT struct {
	UserName string
	Password string
	Token    *string
}

func userProvideValidUserCredentials(ctx context.Context) (context.Context, error) {
	return context.WithValue(ctx, "user", userIT{userEmail, "123456", nil}), nil
}

func testUrl(path string) string {
	return fmt.Sprintf("http://localhost:8081%s", path)
}

func userCallUserTokenEndpoint(ctx context.Context) (context.Context, error) {
	var body []byte = []byte(fmt.Sprintf("{\"email\": \"%s\"}", userEmail))

	request, err := http.NewRequest("POST", testUrl("/user/token"), bytes.NewBuffer(body))
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return ctx, err
	}

	if response.StatusCode != 200 {
		return ctx, fmt.Errorf("not OK")
	}

	var result = UserTokenModel{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	json.Unmarshal(buf.Bytes(), &result)

	var user userIT = ctx.Value("user").(userIT)
	user.Token = result.Token

	return context.WithValue(ctx, "user", user), nil
}

func userHttpResponseContainsValidUserToken(ctx context.Context) error {
	var user userIT = ctx.Value("user").(userIT)

	if user.Token == nil {
		return fmt.Errorf("generated token is null")
	}

	return nil
}

func userProvideValidToken(ctx context.Context) (context.Context, error) {
	userTokenService := service.UserTokenService{}
	token, err := userTokenService.Create(userEmail, config.JwtSecret())

	return context.WithValue(ctx, "user-token", token), err
}

func userProvideInvalidToken(ctx context.Context) (context.Context, error) {
	invalidUserToken := "invalid token"
	return context.WithValue(ctx, "user-token", &invalidUserToken), nil
}

func userCallUserInfoEndpoint(ctx context.Context) (context.Context, error) {
	token := ctx.Value("user-token").(*string)
	request, err := http.NewRequest("GET", testUrl("/user/info"), bytes.NewBuffer([]byte{}))
	request.Header["Authorization"] = []string{*token}
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, "user-info-response", response), nil
}

func userHttpResponseContainsValidUserInformation(ctx context.Context) error {
	response := ctx.Value("user-info-response").(*http.Response)

	if response.StatusCode != 200 {
		return fmt.Errorf("cannot get user's information")
	}

	var userInfo = UserInfoModel{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	json.Unmarshal(buf.Bytes(), &userInfo)

	if userInfo.ID.String() != userId {
		return errors.New(fmt.Sprintf("user id is not as expected %v", userId))
	}

	if userInfo.Email != userEmail {
		return errors.New(fmt.Sprintf("user email is not as expected [%s]", userEmail))
	}

	return nil
}

func _401StatusIsReturned(ctx context.Context) error {
	response := ctx.Value("user-info-response").(*http.Response)

	if response.StatusCode != 401 {
		return errors.New("response is not 401")
	} else {
		return nil
	}
}

func userToken(sc *godog.ScenarioContext) {
	sc.Step(`^valid username and password$`, userProvideValidUserCredentials)
	sc.Step(`^call /user/token endpoint$`, userCallUserTokenEndpoint)
	sc.Step(`^user token is returned$`, userHttpResponseContainsValidUserToken)

	sc.Step(`^valid user token$`, userProvideValidToken)
	sc.Step(`^call /user/info endpoint$`, userCallUserInfoEndpoint)
	sc.Step(`^user information is returned$`, userHttpResponseContainsValidUserInformation)

	sc.Step(`^invalid user token$`, userProvideInvalidToken)
	sc.Step(`^401 status is returned$`, _401StatusIsReturned)
}

func InitializeScenario(sc *godog.ScenarioContext) {
	userToken(sc)
}
