package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/go-meli-toolkit/goutils/logger"
)

func Route(engine *gin.Engine, handler Handler) {
	engine.GET("/ping", handler.Ping)
	engine.POST("/run", handler.Run)
}

func Run(engine *gin.Engine) {
	if err := engine.Run(); err != nil {
		logger.Panic("error running engine", err)
	}
}
