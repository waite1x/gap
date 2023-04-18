package gap

import "context"

type ConfigureFunc = func(c *AppContext) error
type RunFunc = func(c *Application) error

type FuncInfo[T any] struct {
	order int
	run   T
}

type AppContext struct {
	ctx context.Context

	Configures []FuncInfo[ConfigureFunc]
	Runs       []FuncInfo[RunFunc]

	Name        string
	Description string
	Version     string

	data map[string]any
}

func NewAppContext() *AppContext {
	return &AppContext{
		ctx:        context.Background(),
		Configures: make([]FuncInfo[ConfigureFunc], 0),
		Runs:       make([]FuncInfo[RunFunc], 0),
		data:       make(map[string]any),
	}
}

func (a *AppContext) Configure(order int, action ConfigureFunc) {
	a.Configures = append(a.Configures, FuncInfo[ConfigureFunc]{
		run:   action,
		order: order,
	})
}

func (a *AppContext) RunOrder(order int, action RunFunc) {
	a.Runs = append(a.Runs, FuncInfo[RunFunc]{
		run:   action,
		order: order,
	})
}

func (a *AppContext) Context() context.Context {
	return a.ctx
}

func (a *AppContext) WithContext(ctx context.Context) {
	a.ctx = ctx
}

func (a *AppContext) Get(key string) (any, bool) {
	v, ok := a.data[key]
	return v, ok
}

func (a *AppContext) Set(key string, value any) {
	a.data[key] = value
}

func (a *AppContext) TrySet(key string, value any) {
	if _, ok := a.data[key]; !ok {
		a.data[key] = value
	}
}
