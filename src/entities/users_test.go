package entities_test

import (
	"testing"

	. "github.com/kaspar-p/bee/src/entities"
	"github.com/kaspar-p/bee/src/test"
)

func TestInitializeUsers(t *testing.T) {
	t.Parallel()

	guildIds := []string{
		"guild 1",
		"guild 2",
		"guild 3",
	}

	test.Assert(t, len(Users) == 0, "Length of users == 0 before insertion")
	InitializeUsers(guildIds)
	test.Assert(t, len(Users) == len(guildIds), "Length of Users == 3 after insertion")
}
