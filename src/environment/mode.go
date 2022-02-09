package environment

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Mode int64

const (
	PRODUCTION Mode = iota
	DEVELOPMENT
	TESTING
)

func (mode Mode) IsProduction() bool {
	return mode == PRODUCTION
}

func (mode Mode) IsTesting() bool {
	return mode == TESTING
}

func (mode Mode) IsDevelopment() bool {
	return mode == DEVELOPMENT
}

func loadConfigurationFromFile(filepath, filename string) {
	viper.SetConfigName(filename)
	viper.AddConfigPath(filepath)
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic("Error reading from environment variables file: ", err)
	}
}

func (mode Mode) ConfigureEnvironmentVariables() {
	switch {
	case mode.IsDevelopment():
		// Environment variables stored on-file
		loadConfigurationFromFile(getProjectRootPath(), mode.ConfigFile())
	case mode.IsProduction():
		// Heroku stores the environment variables for us
		viper.AutomaticEnv()
	case mode.IsTesting():
		// The 'CI' key is one Github always sets to true
		if os.Getenv("CI") != "" {
			// Running tests from Github actions
			viper.AutomaticEnv()
		} else {
			// Running tests from terminal or local machine
			loadConfigurationFromFile(getProjectRootPath(), mode.ConfigFile())
		}
	}
}

func (m Mode) String() string {
	switch m {
	case PRODUCTION:
		return "production"
	case DEVELOPMENT:
		return "development"
	case TESTING:
		return "testing"
	}

	return ""
}

func (m Mode) ConfigFile() string {
	switch m {
	case PRODUCTION:
		// Currently unused - the variables are in the deploy location
		// Currently, that's Heroku
		return "env.prod"
	case DEVELOPMENT:
		return "env.dev"
	case TESTING:
		// Only used on local tests
		// Github Actions has separate tests
		return "env.test"
	}

	return ""
}
