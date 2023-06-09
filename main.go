package main

import (
	"go-api-meli/routes"
)

func main() {
	build := routes.BuildControllers()
	routes.Routers(build)

}
