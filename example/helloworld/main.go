package main

import (
	"log"

	slim "github.com/waite1x/gapp"
)

func main() {
	ab := slim.NewAppBuilder()

	ab.Use(helloSlim)

	app, err := ab.Build()
	if err != nil {
		log.Fatal(err)
	}
	app.Run()
}

func helloSlim(ab *slim.AppBuilder) {
	ab.Configure(func() error {
		log.Println("Hello Slim will start!")
		return nil
	})

	ab.PostRun(func() error {
		log.Println("Hello Slim will exit!")
		return nil
	})

	ab.Run(func() error {
		log.Println("Hello Slim!")
		return nil
	})
}
