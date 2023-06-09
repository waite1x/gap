package gap

import "sort"

type VersionInfo struct {
	Version string
	Date    string
	Commit  string
	BuiltBy string
}

func (v *VersionInfo) String() string {
	return v.Version
}

type Application struct {
	Ctx *AppContext
}

func NewDefaultApp(ctx *AppContext) *Application {
	return &Application{
		Ctx: ctx,
	}
}

func (app *Application) Run() error {
	if err := app.runConfigures(); err != nil {
		return err
	}
	if err := app.runApp(); err != nil {
		return err
	}
	return nil
}

func (app *Application) runConfigures() error {
	sort.SliceStable(app.Ctx.Configures, func(i, j int) bool {
		return app.Ctx.Configures[i].order < app.Ctx.Configures[j].order
	})
	for _, run := range app.Ctx.Configures {
		if err := run.run(app.Ctx); err != nil {
			return err
		}
	}
	return nil
}

func (app *Application) runApp() error {
	sort.SliceStable(app.Ctx.Runs, func(i, j int) bool { return app.Ctx.Runs[i].order < app.Ctx.Runs[j].order })
	for _, run := range app.Ctx.Runs {
		if err := run.run(app); err != nil {
			return err
		}
	}
	return nil
}
