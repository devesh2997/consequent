package config

import "flag"

// GetEnvironment returns the current env of app
func GetEnvironment() string {
	var env string

	flag.StringVar(&env, "e", "local", "App environment")
	flag.Parse()

	if env == "" {
		env = "local"
	}

	return env
}
