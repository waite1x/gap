package gap

import "context"

const (
	OrderRun      int = 0
	OrderAfterRun int = 10000
)

type AppBuilder struct {
	context    *AppContext
	appCreator func() *Application
}

func NewAppBuilder() *AppBuilder {
	return &AppBuilder{
		context: NewAppContext(),
	}
}

func (ab *AppBuilder) UseApp(creator func() *Application) *AppBuilder {
	ab.appCreator = creator
	return ab
}

func (ab *AppBuilder) Build() (*Application, error) {
	if ab.appCreator == nil {
		return NewDefaultApp(ab.context), nil
	}
	app := ab.appCreator()
	return app, nil
}

func (ab *AppBuilder) SetVersion(version string) *AppBuilder {
	ab.context.Version = version
	return ab
}

func (ab *AppBuilder) SetName(v string) *AppBuilder {
	ab.context.Name = v
	return ab
}

func (ab *AppBuilder) WithContext(ctx context.Context) *AppBuilder {
	ab.context.WithContext(ctx)
	return ab
}

func (ab *AppBuilder) Description(v string) *AppBuilder {
	ab.context.Description = v
	return ab
}

func (ab *AppBuilder) Set(k string, v any) *AppBuilder {
	ab.context.Set(k, v)
	return ab
}

func (ab *AppBuilder) TrySet(k string, v any) *AppBuilder {
	ab.context.TrySet(k, v)
	return ab
}

func (ab *AppBuilder) Get(k string) (any, bool) {
	return ab.context.Get(k)
}

func (ab *AppBuilder) Configure(action func(*AppContext)) {
	ab.context.Configure(0, action)
}

func (ab *AppBuilder) Run(action RunFunc) {
	ab.context.RunOrder(OrderRun, action)
}

func (ab *AppBuilder) PostRun(action RunFunc) {
	ab.context.RunOrder(OrderAfterRun, action)
}

func (ab *AppBuilder) RunOrder(order int, action RunFunc) {
	ab.context.RunOrder(order, action)
}

func (ab *AppBuilder) Use(action func(*AppBuilder)) *AppBuilder {
	action(ab)
	return ab
}
