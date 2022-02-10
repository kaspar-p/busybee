package commands_test

import (
	"os"
	"testing"

	"github.com/kaspar-p/busybee/src/environment"
	"github.com/kaspar-p/busybee/src/test"
	"github.com/kaspar-p/gourd"
)

var (
	tester gourd.Tester
	config *environment.Config
)

func TestMain(m *testing.M) {
	_, teardown := test.SetupDiscordRequiredTests()
	config = environment.InitializeViper(environment.TESTING)

	gourdConfig := gourd.Config{
		AppId:       config.TestingConfig.AppId,
		BotToken:    config.TestingConfig.BotToken,
		TestChannel: "938447408923299890",
		TestingBot:  config.DiscordConfig.AppId,
	}

	var disconnect func()
	tester, disconnect = gourd.CreateTester(gourdConfig)

	code := m.Run()

	disconnect()
	teardown()

	os.Exit(code)
}

func TestWydTalkToBot(t *testing.T) {
	tester.ExpectSending(".wyd <@" + config.DiscordConfig.AppId + ">").ToReturn("nothing much \\;)")
	tester.ExpectSending(".wyd").ToContain("@ of a user")
}

func TestWydNoArgs(t *testing.T) {
	tester.ExpectSending(".wyd").ToContain("@ of a user \\:)")
}
