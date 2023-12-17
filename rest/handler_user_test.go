package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"rest_app/service"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

var userInfo = service.UserInfo{Id: uuid.FromStringOrNil("2f3d92a9-f4db-497b-bf16-367bf2ad7e20"), Email: "test@company.com"}

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
