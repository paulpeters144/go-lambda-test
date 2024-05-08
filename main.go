package main

import (
	"echo-server/cmd/api"
)

func main() {
	app := api.New()
	app.Logger.Fatal(app.Start(":1323"))
}
