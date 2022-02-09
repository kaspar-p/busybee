package environment

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/kaspar-p/busybee/src/discord"
	"github.com/kaspar-p/busybee/src/persist"
	"github.com/spf13/viper"
)

type Config struct {
	TestingConfig  *discord.DiscordConfig
	DiscordConfig  *discord.DiscordConfig
	DatabaseConfig *persist.DatabaseConfig
}

func DecideMode() Mode {
	var mode Mode

	if os.Getenv("MODE") == PRODUCTION.String() {
		// This line is necessary for heroku to run the `web` process correctly
		// PORT needs to be fetched at some point
		log.Println("Port was before:", os.Getenv("PORT"))
		os.Setenv("PORT", "3000")
		os.Setenv("$PORT", "3000")
		log.Println("Port is now:", os.Getenv("PORT"))

		mode = PRODUCTION
	} else {
		mode = DEVELOPMENT
	}

	return mode
}

// Gets the path of the first directory with a .git/ directory in it. This should be the project root.
func getProjectRootPath() string {
	cmdOut, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		log.Panic("Error getting project root directory: ", err)
	}

	return strings.TrimSpace(string(cmdOut))
}

func InitializeViper(mode Mode) *Config {
	if mode == PRODUCTION {
		viper.AutomaticEnv()
	} else {
		viper.SetConfigName(mode.ConfigFile())
		viper.AddConfigPath(getProjectRootPath())
		viper.SetConfigType("yml")

		err := viper.ReadInConfig()
		if err != nil {
			log.Panic("Error reading from environment variables file: ", err)
		}
	}

	config := &Config{
		DiscordConfig: &discord.DiscordConfig{
			BotToken: viper.GetString("BUSYBEE_BOT.TOKEN"),
			AppId:    viper.GetString("BUSYBEE_BOT.APP_ID"),
		},
		DatabaseConfig: &persist.DatabaseConfig{
			ConnectionUrl: viper.GetString("MONGO_DB.CONNECTION_URL"),
			DatabaseName:  viper.GetString("MONGO_DB.DATABASE_NAME"),
			CollectionNames: &persist.CollectionNames{
				Users:     viper.GetString("MONGO_DB.COLLECTIONS.USERS_NAME"),
				BusyTimes: viper.GetString("MONGO_DB.COLLECTIONS.BUSYTIMES_NAME"),
				Guilds:    viper.GetString("MONGO_DB.COLLECTIONS.GUILDS_NAME"),
			},
		},
	}

	if mode == TESTING {
		config.TestingConfig = &discord.DiscordConfig{
			BotToken: viper.GetString("GOURD_BOT.TOKEN"),
			AppId:    viper.GetString("GOURD_BOT.APP_ID"),
		}
	}

	return config
}
