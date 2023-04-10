package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/waite1x/gap"
	"github.com/waite1x/gap/common/di"
)

const ServerBuilderName string = "GinServerBuilder"

type ServerBuiler struct {
	App     *gap.AppBuilder
	preRuns []ServerConfigureFunc
	initors []ServerConfigureFunc
	Options *ServerOptions
	// 只在程序启动过程中进行操作，以保证协程安全
	Items map[string]any
}

func UseServer(ab *gap.AppBuilder) *ServerBuiler {
	sb := addServer(ab)
	sb.Use(DefaultMiddleware)
	return sb
}

func newServerBuilder(builder *gap.AppBuilder) *ServerBuiler {
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

func addServer(ab *gap.AppBuilder) *ServerBuiler {
	sb, ok := ab.Get(ServerBuilderName)
	if !ok {
		serverBuilder := newServerBuilder(ab)
		ab.Configure(func(app *gap.AppContext) {
			di.AddValue(serverBuilder.Options)
		})
		ab.RunOrder(gap.OrderAfterRun-1, func(app *gap.Application) error {
			server := serverBuilder.Build()
			return server.Run()
		})
		ab.Set(ServerBuilderName, sb)
		return serverBuilder
	}
	return sb.(*ServerBuiler)
}
