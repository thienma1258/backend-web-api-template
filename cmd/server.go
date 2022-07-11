package cmd

import (
	"context"
	api "github.com/thienma1258/backend-web-api"
	"log"
)

func main() {
	ctx := context.Background()
	app, err := api.NewApp(ctx)
	if err != nil {
		panic(err)
	}
	defer app.Stop(ctx)
	if err := app.Start(ctx); err != nil {
		log.Panic(err)
	}
}
