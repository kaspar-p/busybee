package database

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func (database *Database) GetRoleIdsForGuilds(guildIds []string) []GuildRolePair {
	filter := bson.D {{
		Key: "GuildId", 
		Value: bson.D{{ 
			Key: "$in", Value: guildIds,
		}},
	}};

	cursor, err := database.guilds.Find(database.context,filter);
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

func (database *Database) GetRoleIdsForGuild(guildId string) []string {
	filter := bson.D{{ Key: "GuildId", Value: guildId }};
	cursor, err := database.guilds.Find(database.context, filter);
	if err != nil {
		fmt.Println("Error getting cursor when finding all users. Error: ", err)
		panic(&GetUserError{Err: err})
	}

	var results []bson.M
	if err = cursor.All(database.context, &results); err != nil {
		fmt.Println("Error getting results from cursor when getting all users. Error: ", err)
		panic(&GetGuildRolePairError{Err: err})
	}

	roleIds := make([]string, 0);
	for _, result := range results {
		roleId := result["RoleId"].(string)
		roleIds = append(roleIds, roleId)
	}

	return roleIds;
}

func (database *Database) RemoveGuildRolePairByGuildAndRole(guildId string, roleId string) error {
	if database == nil {
		return &DatabaseUninitializedError{};
	}

	filter := bson.D {
		{ Key: "GuildId", Value: guildId },
		{ Key: "RoleId", Value: roleId },
	};
	_, err := database.busyTimes.DeleteOne(database.context, filter);
	fmt.Printf("Deleted GuildRolePair that belonged to guild %s and role %s.\n", guildId, roleId);

	return err;
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
	
	var result GuildRolePair;
	err := database.guilds.FindOne(database.context, filter).Decode(&result);

	if err == mongo.ErrNoDocuments {
		fmt.Println("No pairs in database found for guildId:", guildId);
		return false;
	} else if err != nil {
		fmt.Println("Error found when getting single guild pair!");
		// TODO: find a better way of handling an error case!
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

func (database *Database) UpdateGuildRolePairWithNewRole(guildId string, oldRoleId string, newRoleId string) {
	if database == nil {
		panic(&DatabaseUninitializedError{});
	}

	filter := bson.D{
		{ Key: "GuildId", Value: guildId }, 
		{ Key: "RoleId", Value: oldRoleId },
	}
	update := bson.D{{ Key: "RoleId", Value: newRoleId }}
	updateResult, err := database.guilds.UpdateOne(database.context, filter, update);
	if err != nil {
		fmt.Printf("Error while updating guild-role pair with a new role. Guild ID: %s, Old role ID: %s, New role ID: %s. Error: %v.\n", guildId, oldRoleId, newRoleId, err);
		return;
	}
	fmt.Println("Updated", updateResult.ModifiedCount, "guild-role pairs with a new ID");
}