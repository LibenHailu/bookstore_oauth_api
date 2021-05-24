package http

import (
	"net/http"
	"strings"

	atDomain "github.com/LibenHailu/bookstore_oauth_api/src/domain/access_token"
	"github.com/LibenHailu/bookstore_oauth_api/src/services/access_token"
	"github.com/LibenHailu/bookstore_users_api/utils/errors"
	"github.com/gin-gonic/gin"
)

type AccessTokenHandler interface {
	GetById(*gin.Context)
	Create(*gin.Context)
}
type accessTokenHandler struct {
	service access_token.Service
}

func NewHandler(service access_token.Service) AccessTokenHandler {
	return &accessTokenHandler{
		service: service,
	}
}

func (handler *accessTokenHandler) GetById(c *gin.Context) {
	accessTokenId := strings.TrimSpace(c.Param("access_token_id"))
	access_token, err := handler.service.GetById(accessTokenId)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, access_token)
}

func (handler *accessTokenHandler) Create(c *gin.Context) {
	var at atDomain.AccessTokenRequest
	if err := c.ShouldBindJSON(&at); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	accessToken, err := handler.service.Create(at)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusCreated, accessToken)
}
