package commands_test

import (
	"testing"

	"github.com/kaspar-p/busybee/src/environment"
	"github.com/kaspar-p/busybee/src/test"
	"github.com/kaspar-p/gourd"
)

func TestWydTalkToBot(t *testing.T) {
	t.Parallel()

	_, teardown := test.SetupDiscordRequiredTests()
	config := environment.InitializeViper(environment.TESTING)

	gourdConfig := gourd.Config{
		AppId:       config.TestingConfig.AppId,
		BotToken:    config.TestingConfig.BotToken,
		TestChannel: "938447408923299890",
		TestingBot:  config.DiscordConfig.AppId,
	}
	tester, disconnect := gourd.CreateTester(gourdConfig)

	tester.ExpectSending(".wyd <@" + config.DiscordConfig.AppId + ">").ToReturn("nothing much \\;)")

	disconnect()
	teardown()
}

func TestWydNoArgs(t *testing.T) {
	t.Parallel()

	_, teardown := test.SetupDiscordRequiredTests()
	config := environment.InitializeViper(environment.TESTING)

	gourdConfig := gourd.Config{
		AppId:       config.TestingConfig.AppId,
		BotToken:    config.TestingConfig.BotToken,
		TestChannel: "938447408923299890",
		TestingBot:  config.DiscordConfig.AppId,
	}
	tester, disconnect := gourd.CreateTester(gourdConfig)

	tester.ExpectSending(".wyd").ToReturn("command must have a single argument of the @ of a user \\:)")

	disconnect()
	teardown()
}
