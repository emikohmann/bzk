package main

import (
	"github.com/emikohmann/bzk/api"
	"github.com/mercadolibre/go-meli-toolkit/gingonic/mlhandlers"
)

var (
	service = api.NewServiceImpl()
	handler = api.NewHandlerImpl(service)
	router  = mlhandlers.DefaultMeliRouter()
)

func main() {
	api.Route(router, handler)
	api.Run(router)
}
