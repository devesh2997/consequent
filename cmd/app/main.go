package main

import (
	"github.com/devesh2997/consequent/app"
	"github.com/devesh2997/consequent/cmd/flags"
)

func main() {
	env := flags.GetEnvironment()
	app.InitApp(env)
}
