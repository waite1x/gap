package main

import (
	"log"

	"github.com/waite1x/gapp"
)

func main() {
	ab := gapp.NewAppBuilder()

	ab.Use(helloSlim)

	app, err := ab.Build()
	if err != nil {
		log.Fatal(err)
	}
	app.Run()
}

func helloSlim(ab *gapp.AppBuilder) {
	ab.Configure(func(app *gapp.AppContext) {
		log.Println("Hello Slim will start!")
	})

	ab.PostRun(func(app *gapp.Application) error {
		log.Println("Hello Slim will exit!")
		return nil
	})

	ab.Run(func(app *gapp.Application) error {
		log.Println("Hello Slim!")
		return nil
	})
}
