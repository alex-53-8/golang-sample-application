package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticator_whenNoAuthorizationHeader_then401Returned(t *testing.T) {
	var userTokenService = MockUserTokenService{ValidateValue: true}
	var cfg = MockConfigurationKeys{}
	var authenticator gin.HandlerFunc = Authenticator(&userTokenService, &cfg)
	var request http.Request = http.Request{Header: map[string][]string{}}
	w := httptest.NewRecorder()
	var context, _ = gin.CreateTestContext(w)

	context.Request = &request

	authenticator(context)

	assert.Equal(t, 401, w.Code, "statuses do not equal")
	assert.False(t, userTokenService.ValidateWasCalled)
}

func TestAuthenticator_whenWrongTokenInAuthorizationHeader_then401Returned(t *testing.T) {
	var userTokenService = MockUserTokenService{ValidateValue: false, DecodedValueValid: false}
	var cfg = MockConfigurationKeys{}
	var authenticator gin.HandlerFunc = Authenticator(&userTokenService, &cfg)
	var request http.Request = http.Request{Header: map[string][]string{
		"Authorization": {"wrong-token"},
	}}
	w := httptest.NewRecorder()
	var context, _ = gin.CreateTestContext(w)

	context.Request = &request

	authenticator(context)

	assert.Equal(t, 401, w.Code, "statuses do not equal")
	assert.True(t, userTokenService.DecodeWasCalled)
}

func TestAuthenticator_whenProperAuthorizationHeader_thenContextChainProceed(t *testing.T) {
	var userTokenService = MockUserTokenService{ValidateValue: true, DecodedValueValid: true}
	var cfg = MockConfigurationKeys{}
	var authenticator gin.HandlerFunc = Authenticator(&userTokenService, &cfg)
	var request http.Request = http.Request{Header: map[string][]string{
		"Authorization": {"mocked-token"},
	}}
	w := httptest.NewRecorder()
	var context, _ = gin.CreateTestContext(w)

	context.Request = &request

	authenticator(context)

	assert.Equal(t, 200, w.Code, "statuses do not equal")
	assert.True(t, userTokenService.DecodeWasCalled)
}
