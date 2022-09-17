package flags

import "flag"

var Create = flag.String("create", "", "create new migrations files (up and down) for the given title")
var Steps = flag.Int("steps", 0, "migrate up if steps > 0, and down if steps < 0")
var ENV = flag.String("e", "local", "App environment")

func Parse() {
	flag.Parse()
}

// GetEnvironment returns the current env of app
func GetEnvironment() string {
	Parse()

	if *ENV == "" {
		*ENV = "local"
	}

	return *ENV
}
