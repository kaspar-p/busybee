package constants

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	GuildID string
	ChannelID string
	BotToken string
	AppID string
	ConnectionURL string
	DatabaseName string
	UsersCollectionName string
	CoursesCollectionName string
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
	activeServer := viper.GetString("BOT.ACTIVE_SERVER");
	GuildID = viper.GetString("BOT.GUILD_IDS." + activeServer);
	ChannelID = viper.GetString("BOT.CHANNEL_IDS." + activeServer);
	AppID = viper.GetString("BOT.APP_ID");

	// Database constants
	ConnectionURL = viper.GetString("MONGO_DB.CONNECTION_URL")
	DatabaseName = viper.GetString("MONGO_DB.DATABASE_NAME")
	UsersCollectionName = viper.GetString("MONGO_DB.COLLECTIONS.USERS_NAME")
	CoursesCollectionName = viper.GetString("MONGO_DB.COLLECTIONS.COURSES_NAME")
	BusyTimesCollectionName = viper.GetString("MONGO_DB.COLLECTIONS.BUSYTIMES_NAME")
}