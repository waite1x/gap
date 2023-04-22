package gap

import (
	"github.com/waite1x/gap/di"
	"github.com/waite1x/gap/log"
)

// Logger 注册日志组件
func Logger(ab *AppBuilder) *AppBuilder {
	ab.Configure(func(app *AppContext) error {
		di.TryAddValue(log.Opts)
		di.TryAddValue(log.GetLogger())
		return nil
	})
	return ab
}
