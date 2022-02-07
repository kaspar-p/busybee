package environment

import (
	"log"

	"github.com/kaspar-p/busybee/src/discord"
	"github.com/kaspar-p/busybee/src/persist"
	"github.com/spf13/viper"
)

type Config struct {
	DiscordConfig  *discord.DiscordConfig
	DatabaseConfig *persist.DatabaseConfig
}

func InitializeViper(mode Mode) *Config {
	viper.SetConfigName(mode.ConfigFile())
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic("Error reading from environment variables file: ", err)
	}

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
