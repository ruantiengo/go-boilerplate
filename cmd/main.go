package main

import "ruantiengo/config"

func main() {
	app := config.NewApplication()
	defer app.Shutdown()

	app.Run()
}
