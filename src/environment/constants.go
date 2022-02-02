package environment

import (
	"fmt"
	"log"

	"github.com/kaspar-p/bee/src/discord"
	"github.com/kaspar-p/bee/src/persist"
	"github.com/spf13/viper"
)

type Config struct {
	DiscordConfig  *discord.DiscordConfig
	DatabaseConfig *persist.DatabaseConfig
}

func InitializeViper(mode Mode) *Config {
	viper.SetConfigName("env.prod")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic("Error reading from environment variables file: ", err)
	}

	return &Config{
		DiscordConfig: &discord.DiscordConfig{
			BotToken: viper.GetString(fmt.Sprintf("BUSYBEE_BOT.%s.TOKEN", mode.ConfigString())),
			AppId:    viper.GetString(fmt.Sprintf("BUSYBEE_BOT.%s.APP_ID", mode.ConfigString())),
		},
		DatabaseConfig: &persist.DatabaseConfig{
			ConnectionUrl: viper.GetString(fmt.Sprintf("MONGO_DB.%s.CONNECTION_URL", mode.ConfigString())),
			DatabaseName:  viper.GetString(fmt.Sprintf("MONGO_DB.%s.DATABASE_NAME", mode.ConfigString())),
			CollectionNames: &persist.CollectionNames{
				Users:     viper.GetString(fmt.Sprintf("MONGO_DB.%s.COLLECTIONS.USERS_NAME", mode.ConfigString())),
				BusyTimes: viper.GetString(fmt.Sprintf("MONGO_DB.%s.COLLECTIONS.BUSYTIMES_NAME", mode.ConfigString())),
				Guilds:    viper.GetString(fmt.Sprintf("MONGO_DB.%s.COLLECTIONS.GUILDS_NAME", mode.ConfigString())),
			},
		},
	}
}
