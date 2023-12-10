package rest

import (
	"rest_app/cache"
	"rest_app/cfg"
	"rest_app/service"

	"github.com/gin-gonic/gin"
)

func AddUserEndpoints(
	groups *RouterGroups,
	_cache cache.Cache,
	userTokenService service.UserToken,
	userService service.User,
	cfg cfg.Configuration,
) {
	groups.public.POST("/user/token", func(c *gin.Context) {
		UserCreateTokenHandler(c, _cache, userTokenService, userService, cfg)
	})

	groups.auth.GET("/user/info", func(c *gin.Context) {
		UserInfoHandler(c, _cache, userService)
	})
}
