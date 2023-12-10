package rest

import (
	"context"
	"log"
	"net/http"
	"rest_app/cache"
	"rest_app/cfg"
	"rest_app/database"
	"rest_app/service"

	"github.com/gin-gonic/gin"
)

type RouterGroups struct {
	public *gin.RouterGroup
	auth   *gin.RouterGroup
}

type RestServer struct {
	restServer *http.Server
}

func (s *RestServer) Start() {
	if err := s.restServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func (s *RestServer) Stop() {
	s.restServer.Shutdown(context.Background())
}

func CreateRestServer(_cache cache.Cache, db database.Database, cfg cfg.Configuration) *RestServer {
	var router = gin.Default()
	var userTokenService = service.UserTokenService{}
	var userService = service.UserService{Db: db}

	routerGroups := RouterGroups{
		public: router.Group("/"),
		auth:   router.Group("/"),
	}
	routerGroups.auth.Use(Authenticator(&userTokenService, cfg))

	AddHealthEndpoints(&routerGroups)
	AddUserEndpoints(&routerGroups, _cache, &userTokenService, &userService, cfg)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	return &RestServer{srv}
}
