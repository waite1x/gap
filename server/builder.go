package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/waite1x/gapp"
)

const ServerBuilderName string = "GinServerBuilder"

type ServerBuiler struct {
	App     *gapp.AppBuilder
	preRuns []ServerConfigureFunc
	initors []ServerConfigureFunc
	Options *ServerOptions
	// 只在程序启动过程中进行操作，以保证协程安全
	Items map[string]any
}

func UseServer(ab *gapp.AppBuilder) *ServerBuiler {
	return addServer(ab)
}

func newServerBuilder(builder *gapp.AppBuilder) *ServerBuiler {
	return &ServerBuiler{
		App:     builder,
		preRuns: make([]ServerConfigureFunc, 0),
		initors: make([]ServerConfigureFunc, 0),
		Options: &ServerOptions{
			LogLevel: zerolog.InfoLevel,
		},
		Items: make(map[string]any),
	}
}

func (b *ServerBuiler) PreConfigure(action ServerConfigureFunc) *ServerBuiler {
	b.preRuns = append(b.preRuns, action)
	return b
}

func (b *ServerBuiler) Configure(action ServerConfigureFunc) *ServerBuiler {
	b.initors = append(b.initors, action)
	return b
}

func (b *ServerBuiler) Use(module func(*ServerBuiler)) *ServerBuiler {
	module(b)
	return b
}

func (b *ServerBuiler) Build() *Server {
	g := gin.Default()
	server := NewServer(g, b.Options)
	for _, action := range b.preRuns {
		action(server)
	}
	for _, action := range b.initors {
		action(server)
	}
	return server
}

func addServer(ab *gapp.AppBuilder) *ServerBuiler {
	sb, ok := ab.Get(ServerBuilderName)
	if !ok {
		sb = newServerBuilder(ab)
		ab.RunOrder(gapp.OrderAfterRun-1, func(app *gapp.Application) error {
			server := sb.(*ServerBuiler).Build()
			return server.Run()
		})
		ab.Set(ServerBuilderName, sb)
	}
	return sb.(*ServerBuiler)
}
