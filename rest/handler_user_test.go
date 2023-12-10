package rest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"rest_app/service"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

var userInfo = service.UserInfo{Id: uuid.FromStringOrNil("2f3d92a9-f4db-497b-bf16-367bf2ad7e20"), Email: "test@company.com"}

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

func TestUserCreateTokenHandler(t *testing.T) {
	w := httptest.NewRecorder()
	var context, _ = gin.CreateTestContext(w)
	var userTokenService = MockUserTokenService{}
	var userService = MockUserService{}
	var cache = MockCache{}
	var cfg = MockConfigurationKeys{}

	var request *http.Request = &http.Request{}
	request.Body = io.NopCloser(strings.NewReader("{\"email\": \"boss@company.com\"}"))
	context.Request = request

	UserCreateTokenHandler(context, &cache, &userTokenService, &userService, &cfg)

	var result = UserTokenModel{}
	json.Unmarshal(w.Body.Bytes(), &result)

	assert.Equal(t, 200, w.Code)
	assert.NotNil(t, result.Token)
	assert.Equal(t, uint(1), cache.getWasCalled)
	assert.Equal(t, uint(1), cache.storeWasCalled)
}

func TestUserInfoHandler(t *testing.T) {
	w := httptest.NewRecorder()
	var context, _ = gin.CreateTestContext(w)
	var userService = MockUserService{}
	var cache = MockCache{}
	var decodedToken = service.UserTokenDecoded{}

	context.Set(MiddlewareUserTokenDecodedKey, &decodedToken)

	UserInfoHandler(context, &cache, &userService)

	var result = UserInfoModel{}
	json.Unmarshal(w.Body.Bytes(), &result)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "test@company.com", result.Email)
}
