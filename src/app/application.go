package app

import (
	"github.com/LibenHailu/bookstore_oauth_api/src/http"
	"github.com/gin-gonic/gin"

	"github.com/LibenHailu/bookstore_oauth_api/src/domain/access_token"
	"github.com/LibenHailu/bookstore_oauth_api/src/repository/db"
)

var (
	router = gin.Default()
)

func StartApplication() {
	atService := access_token.NewService(db.NewRepository())
	atHanler := http.NewHandler(atService)
	router.GET("/oauth/access_token/:access_token_id", atHanler.GetById)
	router.POST("/oauth/access_token", atHanler.Create)

	router.Run(":8080")
}
