package constants

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	BusyRoleID string
	BusyRoleName string
	BotToken string
	AppID string
	ConnectionURL string
	DatabaseName string
	UsersCollectionName string
	BusyTimesCollectionName string
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

	// Get environment variables
	BotToken = viper.GetString("BOT.TOKEN");
	AppID = viper.GetString("BOT.APP_ID");

	// Database constants
	ConnectionURL = viper.GetString("MONGO_DB.CONNECTION_URL")
	DatabaseName = viper.GetString("MONGO_DB.DATABASE_NAME")
	UsersCollectionName = viper.GetString("MONGO_DB.COLLECTIONS.USERS_NAME")
	BusyTimesCollectionName = viper.GetString("MONGO_DB.COLLECTIONS.BUSYTIMES_NAME")

	// Role constants
	BusyRoleName = viper.GetString("CONSTANTS.ROLE_NAME");
}