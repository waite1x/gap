package main

import (
	"log"

	"github.com/waite1x/gap"
)

func main() {
	ab := gap.NewAppBuilder()

	ab.Use(helloSlim)

	app, err := ab.Build()
	if err != nil {
		log.Fatal(err)
	}
	app.Run()
}

func helloSlim(ab *gap.AppBuilder) {
	ab.Configure(func(app *gap.AppContext) error {
		log.Println("Hello Slim will start!")
		return nil
	})

	ab.PostRun(func(app *gap.Application) error {
		log.Println("Hello Slim will exit!")
		return nil
	})

	ab.Run(func(app *gap.Application) error {
		log.Println("Hello Slim!")
		return nil
	})
}
