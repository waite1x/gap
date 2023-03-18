package slim

type ConfigureFunc = func()
type RunFunc = func() error

type runFunInfo struct {
	order int
	run   RunFunc
}

type AppContext struct {
	Runs []runFunInfo

	Name        string
	Description string
	Version     string

	data map[string]interface{}
}

func NewAppContext() *AppContext {
	return &AppContext{
		Runs: make([]runFunInfo, 0),
		data: make(map[string]interface{}),
	}
}

func (a *AppContext) RunOrder(order int, action RunFunc) {
	a.Runs = append(a.Runs, runFunInfo{
		run:   action,
		order: order,
	})
}

func (a *AppContext) Get(key string) (interface{}, bool) {
	v, ok := a.data[key]
	return v, ok
}

func (a *AppContext) Set(key string, value interface{}) {
	a.data[key] = value
}

func (a *AppContext) TrySet(key string, value interface{}) {
	if _, ok := a.data[key]; !ok {
		a.data[key] = value
	}
}
