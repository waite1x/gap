package main

import (
	"server"

	"github.com/gin-gonic/gin"
	"github.com/waite1x/gapp"
)

func main() {
	ab := gapp.NewAppBuilder()
	server.UseServer(ab).
		Use(configureServer)

	app, err := ab.Build()
	if err != nil {
		panic(err)
	}
	app.Run()
}

func configureServer(sb *server.ServerBuiler) {
	sb.Configure(func(s *server.Server) {
		s.Route.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	})
}
