package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"rest_app/cache"
	"rest_app/cfg"
	"rest_app/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type UserTokenModel struct {
	Token *string `json:"token"`
}

type UserInfoModel struct {
	ID    uuid.UUID
	Email string
}

type UserCreateTokenRequest struct {
	Email string `json:"email"`
}

const keyPatternUserInfo = "user-info-%s"
const cachingTime = time.Duration(20 * float64(time.Minute))

func UserCreateTokenHandler(c *gin.Context, _cache cache.Cache, userTokenService service.UserToken,
	userService service.User, cfg cfg.JwtConfiguration) {
	bytes, err := io.ReadAll(c.Request.Body)

	if err != nil {
		log.Println("error during reading a request for creation of a user's token: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var request = UserCreateTokenRequest{}
	err = json.Unmarshal(bytes, &request)

	if err != nil {
		log.Println("error during unmarshalling a request for creation of a user's token: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userInfoCacheKey := fmt.Sprintf(keyPatternUserInfo, request.Email)
	userInfo, err := cache.GetAndStoreIfMissed[service.UserInfo](_cache, userInfoCacheKey, func() (*service.UserInfo, time.Duration, error) {
		info, err := userService.FindByEmail(request.Email)
		return info, cachingTime, err
	})

	token, err := userTokenService.Create(userInfo.Email, cfg.JwtSecret())

	if err != nil {
		log.Println("error during a user's token creation: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, UserTokenModel{Token: token})
}

func UserInfoHandler(c *gin.Context, _cache cache.Cache, userService service.User) {
	token, exists := c.Get(MiddlewareUserTokenDecodedKey)

	if !exists {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userTokenDecoded := token.(*service.UserTokenDecoded)

	userInfoCacheKey := fmt.Sprintf(keyPatternUserInfo, userTokenDecoded.Email)
	userInfo, err := cache.GetAndStoreIfMissed[service.UserInfo](_cache, userInfoCacheKey, func() (*service.UserInfo, time.Duration, error) {
		info, err := userService.FindByEmail(userTokenDecoded.Email)
		return info, cachingTime, err
	})

	if err != nil {
		log.Println("no user info found: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, UserInfoModel{userInfo.Id, userInfo.Email})
}
