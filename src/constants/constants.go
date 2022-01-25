package constants

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	BeeColor int

	BotReady bool
	BotToken string
	AppID string

	ConnectionURL string
	DatabaseName string
	UsersCollectionName string
	BusyTimesCollectionName string
	GuildsCollectionName string
)

func InitializeViper() {
	viper.SetConfigName("env");
	viper.AddConfigPath(".");
	viper.AutomaticEnv();
	viper.SetConfigType("yml");

	err := viper.ReadInConfig();
	if err != nil {
		fmt.Println("Error reading from environment variables file: ", err);
	}

	// Set constants not dependent on Viper
	BeeColor = 15122779; // Yellow
	BotReady = false;

	// Get environment variables
	BotToken = viper.GetString("BOT.TOKEN");
	AppID = viper.GetString("BOT.APP_ID");

	// Database constants
	ConnectionURL = viper.GetString("MONGO_DB.CONNECTION_URL")
	DatabaseName = viper.GetString("MONGO_DB.DATABASE_NAME")
	UsersCollectionName = viper.GetString("MONGO_DB.COLLECTIONS.USERS_NAME")
	BusyTimesCollectionName = viper.GetString("MONGO_DB.COLLECTIONS.BUSYTIMES_NAME")
	GuildsCollectionName = viper.GetString("MONGO_DB.COLLECTIONS.GUILDS_NAME")
}