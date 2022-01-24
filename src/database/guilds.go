package database

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type GuildRolePair struct {
	GuildId string;
	RoleId string;
}

func (pair GuildRolePair) ConvertGuildRolePairToDocument() bson.D {
	return bson.D {
		{ Key: "GuildId", Value: pair.GuildId },
		{ Key: "RoleId", Value: pair.RoleId },
	}
}

func (database *Database) GetRoleIdsForGuilds() ([]GuildRolePair) {
	cursor, err := database.guilds.Find(database.context, bson.D{{ }});
	if err != nil {
		fmt.Println("Error getting cursor when finding all users. Error: ", err)
		panic(&GetUserError{Err: err})
	}

	var results []bson.M
	if err = cursor.All(database.context, &results); err != nil {
		fmt.Println("Error getting results from cursor when getting all users. Error: ", err)
		panic(&GetUserError{Err: err})
	}

	pairs := make([]GuildRolePair, 0)
	for _, result := range results {
		guildId := result["GuildId"].(string)
		roleId := result["RoleId"].(string)
		pair := GuildRolePair{
			GuildId: guildId,
			RoleId: roleId,
		}
		pairs = append(pairs, pair)
	}

	return pairs;
}

func (database *Database) RemoveGuildRolePairByGuild(guildId string) error {
	if database == nil {
		return &DatabaseUninitializedError{};
	}

	filter := bson.D {{ Key: "GuildId", Value: guildId }};
	_, err := database.busyTimes.DeleteOne(database.context, filter);
	fmt.Println("Deleted GuildRolePair that belonged to guild", guildId);

	return err;
}

func (database *Database) IsGuildInPairMap(guildId string) bool {
	filter := bson.D{{ Key: "GuildID", Value: guildId }};
	result := database.guilds.FindOne(database.context, filter);

	if result.Err() != nil {
		fmt.Println("No pairs in database found for guildId:", guildId);
		return false;
	} else {
		return true;
	}
} 

func (database *Database) AddGuildRolePair(guildId string, roleId string) {
	if database == nil {
		panic(&DatabaseUninitializedError{})
	}

	guildRolePair := GuildRolePair{
		GuildId: guildId,
		RoleId: roleId,
	}

	pairDocument := guildRolePair.ConvertGuildRolePairToDocument();

	_, err := database.guilds.InsertOne(database.context, pairDocument);
	if err != nil {
		fmt.Println("Error inserting guild role pair: ", pairDocument, ". Error: ", err);
		panic(&AddGuildRolePairError{ Err: err })
	}
}