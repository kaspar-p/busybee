package ingest

import (
	"log"

	dbLib "github.com/kaspar-p/bee/src/database"
	"github.com/kaspar-p/bee/src/entities"
	"github.com/kaspar-p/bee/src/update"
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

func FillMapsWithDatabaseData(guildIds []string) {
	guildRolePairs := FetchRoleIdData(guildIds)
	users := FetchUserData()
	allBusyTimes := FetchBusyTimesData()

	for _, user := range users {
		user.SortBusyTimes()
	}

	log.Println("Got data: \n\tUsers:", len(GetAllSubKeysOfUsersMap(entities.Users)),
		"\n\tEvents:", len(allBusyTimes),
		"\n\tGuild-Role pairs:", len(guildRolePairs),
	)
}

func FetchRoleIdData(guildIds []string) []dbLib.GuildRolePair {
	guildRolePairs := dbLib.DatabaseInstance.GetRoleIdsForGuilds(guildIds)

	for _, pair := range guildRolePairs {
		update.GuildRoleMap[pair.GuildId] = pair.RoleId
	}

	return guildRolePairs
}

func FetchUserData() []*entities.User {
	users := dbLib.DatabaseInstance.GetUsers()
	for _, user := range users {
		entities.Users[user.BelongsTo][user.Id] = user
	}

	return users
}

func FetchBusyTimesData() []*entities.BusyTime {
	busyTimesArray := dbLib.DatabaseInstance.GetBusyTimes()
	for _, busyTime := range busyTimesArray {
		user := entities.Users[busyTime.BelongsTo][busyTime.OwnerId]
		user.BusyTimes = append(user.BusyTimes, busyTime)
	}

	return busyTimesArray
}
