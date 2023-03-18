package slim

const (
	OrderBeforeRun int = -1000
	OrderRun       int = 0
	OrderAfterRun  int = 1000
)

type AppBuilder struct {
	Context    *AppContext
	appCreator func() *Application
}

func NewAppBuilder() *AppBuilder {
	return &AppBuilder{
		Context: NewAppContext(),
	}
}

func (ab *AppBuilder) UseApp(creator func() *Application) *AppBuilder {
	ab.appCreator = creator
	return ab
}

func (ab *AppBuilder) Build() (*Application, error) {
	if ab.appCreator == nil {
		return NewDefaultApp(ab.Context), nil
	}
	app := ab.appCreator()
	return app, nil
}

func (ab *AppBuilder) SetVersion(version string) *AppBuilder {
	ab.Context.Version = version
	return ab
}

func (ab *AppBuilder) SetName(v string) *AppBuilder {
	ab.Context.Name = v
	return ab
}

func (ab *AppBuilder) Description(v string) *AppBuilder {
	ab.Context.Description = v
	return ab
}

func (ab *AppBuilder) Set(k string, v interface{}) *AppBuilder {
	ab.Context.Set(k, v)
	return ab
}

func (ab *AppBuilder) TrySet(k string, v interface{}) *AppBuilder {
	ab.Context.TrySet(k, v)
	return ab
}

func (ab *AppBuilder) Configure(action RunFunc) {
	ab.Context.RunOrder(OrderBeforeRun, action)
}

func (ab *AppBuilder) Run(action RunFunc) {
	ab.Context.RunOrder(OrderRun, action)
}

func (ab *AppBuilder) PostRun(action RunFunc) {
	ab.Context.RunOrder(OrderAfterRun, action)
}

func (ab *AppBuilder) RunOrder(order int, action RunFunc) {
	ab.Context.RunOrder(order, action)
}

func (ab *AppBuilder) Use(action func(*AppBuilder)) *AppBuilder {
	action(ab)
	return ab
}
