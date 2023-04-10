package server

import (
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/waite1x/gap/common/di"
)

func DefaultMiddleware(sb *ServerBuiler) {
	sb.PreConfigure(func(s *Server) {
		s.Route.Use(Log(sb.Options.LogLevel))
		s.Route.Use(DependencyInjection())
		s.Route.Use(ErrorMiddleware)
		s.Route.Use(UnitWorkMiddleware())
	})
}

func Log(lvl zerolog.Level) gin.HandlerFunc {
	return logger.SetLogger(
		logger.WithDefaultLevel(lvl),
	)
}

func DependencyInjection() gin.HandlerFunc {
	return func(c *gin.Context) {
		p := di.GetContainer().CreateScope(c)
		c.Set(di.ProviderKey, p)
		c.Next()
		p.Close()
	}
}
