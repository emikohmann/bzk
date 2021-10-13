package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/go-meli-toolkit/goutils/apierrors"
	"net/http"
)

type Handler interface {
	Ping(c *gin.Context)
	Run(c *gin.Context)
}

type HandlerImpl struct {
	Service Service
}

func NewHandlerImpl(service Service) HandlerImpl {
	return HandlerImpl{
		Service: service,
	}
}

func (handlerImpl HandlerImpl) Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func (handlerImpl HandlerImpl) Run(c *gin.Context) {
	var loadTest LoadTest
	if err := c.ShouldBindJSON(&loadTest); err != nil {
		apiErr := apierrors.NewBadRequestApiError(err.Error())
		c.JSON(apiErr.Status(), apiErr)
		return
	}
	result, apiErr := handlerImpl.Service.Run(loadTest)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}
	c.JSON(http.StatusOK, result)
}
