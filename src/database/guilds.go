package database

import (
	"log"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GuildRolePair struct {
	GuildId string
	RoleId  string
}

func (pair GuildRolePair) ConvertGuildRolePairToDocument() bson.D {
	return bson.D{
		{Key: "GuildId", Value: pair.GuildId},
		{Key: "RoleId", Value: pair.RoleId},
	}
}

func (database *Database) GetRoleIdsForGuilds(guildIds []string) []GuildRolePair {
	filter := bson.D{{
		Key: "GuildId",
		Value: bson.D{{
			Key: "$in", Value: guildIds,
		}},
	}}

	cursor, err := database.guilds.Find(database.context, filter)
	if err != nil {
		log.Panic("Error getting cursor when finding all users. Error: ", err)
		panic(&GetUserError{Err: err})
	}

	var results []bson.M

	if err = cursor.All(database.context, &results); err != nil {
		log.Panic("Error getting results from cursor when getting all users. Error: ", err)
		panic(&GetUserError{Err: err})
	}

	pairs := make([]GuildRolePair, 0)

	for _, result := range results {
		guildId, found := result["GuildId"].(string)
		if !found {
			log.Panic("Key 'GuildId' not found on guild-role pair!")
			panic(&GetGuildRolePairError{})
		}

		roleId, found := result["RoleId"].(string)
		if !found {
			log.Panic("Key 'RoleId' not found on guild-role pair!")
			panic(&GetGuildRolePairError{})
		}

		pair := GuildRolePair{
			GuildId: guildId,
			RoleId:  roleId,
		}
		pairs = append(pairs, pair)
	}

	return pairs
}

func (database *Database) GetRoleIdsForGuild(guildId string) []string {
	filter := bson.D{
		{Key: "GuildId", Value: guildId},
	}

	cursor, err := database.guilds.Find(database.context, filter)
	if err != nil {
		log.Panic("Error getting cursor when finding all users. Error: ", err)
		panic(&GetGuildRolePairError{Err: err})
	}

	var results []bson.M
	if err = cursor.All(database.context, &results); err != nil {
		log.Panic("Error getting results from cursor when getting all users. Error: ", err)
		panic(&GetGuildRolePairError{Err: err})
	}

	roleIds := make([]string, 0)

	for _, result := range results {
		roleId, found := result["RoleId"].(string)
		if !found {
			log.Panic("Key 'RoleId' not found on a guild-role pair!")
			panic(&GetGuildRolePairError{Err: err})
		}

		roleIds = append(roleIds, roleId)
	}

	return roleIds
}

func (database *Database) RemoveGuildRolePairByGuildAndRole(guildId, roleId string) error {
	if database == nil {
		return &DatabaseUninitializedError{}
	}

	filter := bson.D{
		{Key: "GuildId", Value: guildId},
		{Key: "RoleId", Value: roleId},
	}
	_, err := database.busyTimes.DeleteOne(database.context, filter)

	log.Printf("Deleted GuildRolePair that belonged to guild %s and role %s.\n", guildId, roleId)

	return errors.Wrap(err, "Error removing guild-role by both guild ID and role ID!")
}

func (database *Database) RemoveGuildRolePairByGuild(guildId string) error {
	if database == nil {
		return &DatabaseUninitializedError{}
	}

	filter := bson.D{{Key: "GuildId", Value: guildId}}
	_, err := database.busyTimes.DeleteOne(database.context, filter)

	log.Println("Deleted GuildRolePair that belonged to guild", guildId)

	return errors.Wrap(err, "Error removing guild-role pair by guild ID.")
}

func (database *Database) IsGuildInPairMap(guildId string) bool {
	filter := bson.D{{Key: "GuildID", Value: guildId}}

	var result GuildRolePair
	err := database.guilds.FindOne(database.context, filter).Decode(&result)

	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		log.Panic("No pairs in database found for guildId:", guildId)

		return false
	case err != nil:
		log.Panic("Error found when getting single guild pair!")
		// TODO: find a better way of handling an error case!
		return false
	default:
		return true
	}
}

func (database *Database) AddGuildRolePair(guildId, roleId string) {
	if database == nil {
		panic(&DatabaseUninitializedError{})
	}

	guildRolePair := GuildRolePair{
		GuildId: guildId,
		RoleId:  roleId,
	}

	pairDocument := guildRolePair.ConvertGuildRolePairToDocument()

	_, err := database.guilds.InsertOne(database.context, pairDocument)
	if err != nil {
		log.Panic("Error inserting guild role pair: ", pairDocument, ". Error: ", err)
		panic(&AddGuildRolePairError{Err: err})
	}
}

func (database *Database) UpdateGuildRolePairWithNewRole(guildId, oldRoleId, newRoleId string) {
	if database == nil {
		panic(&DatabaseUninitializedError{})
	}

	filter := bson.D{
		{Key: "GuildId", Value: guildId},
		{Key: "RoleId", Value: oldRoleId},
	}
	update := bson.D{{Key: "RoleId", Value: newRoleId}}

	updateResult, err := database.guilds.UpdateOne(database.context, filter, update)
	if err != nil {
		log.Panicf("Error while updating guild-role pair with a new role. Guild ID: %s, "+
			"Old role ID: %s, New role ID: %s. Error: %v.\n", guildId, oldRoleId, newRoleId, err)
		// TODO: find a better way of handling this error case
		return
	}

	log.Println("Updated", updateResult.ModifiedCount, "guild-role pairs with a new ID")
}
