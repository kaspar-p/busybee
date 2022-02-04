package ingest

import (
	"log"

	"github.com/kaspar-p/busybee/src/entities"
	"github.com/kaspar-p/busybee/src/persist"
	"github.com/kaspar-p/busybee/src/update"
)

func GetAllSubKeysOfUsersMap(m map[string]map[string]*entities.User) []string {
	allKeys := make([]string, 0)

	for key := range m {
		for subkey := range m[key] {
			allKeys = append(allKeys, subkey)
		}
	}

	return allKeys
}

func FillMapsWithDatabaseData(database *persist.DatabaseType, guildIds []string) {
	guildRolePairs := FetchRoleIdData(database, guildIds)
	users := FetchUserData(database)
	allBusyTimes := FetchBusyTimesData(database)

	for _, user := range users {
		user.SortBusyTimes()
	}

	log.Println("Got data: \n\tUsers:", len(GetAllSubKeysOfUsersMap(entities.Users)),
		"\n\tEvents:", len(allBusyTimes),
		"\n\tGuild-Role pairs:", len(guildRolePairs),
	)
}

func FetchRoleIdData(database *persist.DatabaseType, guildIds []string) []persist.GuildRolePair {
	guildRolePairs := database.GetRoleIdsForGuilds(guildIds)

	for _, pair := range guildRolePairs {
		update.GuildRoleMap[pair.GuildId] = pair.RoleId
	}

	return guildRolePairs
}

func FetchUserData(database *persist.DatabaseType) []*entities.User {
	users := database.GetUsers()
	for _, user := range users {
		entities.Users[user.BelongsTo][user.Id] = user
	}

	return users
}

func FetchBusyTimesData(database *persist.DatabaseType) []*entities.BusyTime {
	busyTimesArray := database.GetBusyTimes()
	for _, busyTime := range busyTimesArray {
		user := entities.Users[busyTime.BelongsTo][busyTime.OwnerId]
		user.BusyTimes = append(user.BusyTimes, busyTime)
	}

	return busyTimesArray
}
