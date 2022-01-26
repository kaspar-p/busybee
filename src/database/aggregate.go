package database

// General methods that define functions that are independent on table or data structure
// Usually relate to an entire guild, or all data

func (database *Database) RemoveAllDataInGuild(guildId string) error {
	err := database.RemoveAllUsersInGuild(guildId)
	if err != nil {
		return err
	}

	err = database.RemoveAllBusyTimesInGuild(guildId)
	if err != nil {
		return err
	}

	err = database.RemoveGuildRolePairByGuild(guildId)

	return err
}
