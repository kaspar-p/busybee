package environment

import (
	"log"
	"os"

	"github.com/kaspar-p/busybee/src/discord"
	"github.com/kaspar-p/busybee/src/persist"
	"github.com/spf13/viper"
)

type Config struct {
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

func InitializeViper(mode Mode) *Config {
	if mode == PRODUCTION {
		viper.AutomaticEnv()
	} else {
		viper.SetConfigName(mode.ConfigFile())
		viper.AddConfigPath(".")
		viper.SetConfigType("yml")

		err := viper.ReadInConfig()
		if err != nil {
			log.Panic("Error reading from environment variables file: ", err)
		}
	}

	log.Println("second", viper.GetString("BUSYBEE_BOT.TOKEN"))

	return &Config{
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
}
