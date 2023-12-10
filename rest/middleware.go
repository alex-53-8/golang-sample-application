package rest

import (
	"net/http"
	"rest_app/cfg"
	"rest_app/service"

	"github.com/gin-gonic/gin"
)

func Authenticator(userTokenService service.UserToken, cfg cfg.Configuration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var authorization []string = c.Request.Header["Authorization"]

		if len(authorization) == 1 {
			token := authorization[0]
			decodedToken, err := userTokenService.Decode(token, cfg.JwtSecret())
			if err == nil && userTokenService.Validate(decodedToken) {
				c.Set(MiddlewareUserTokenDecodedKey, decodedToken)
				c.Next()
				return

			}
		}

		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
