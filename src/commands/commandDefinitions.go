package commands

import "github.com/bwmarrin/discordgo"

var (
	// SLASH COMMAND CODE
	// createdCommands []*discordgo.ApplicationCommand
	// commands = []*discordgo.ApplicationCommand {
	// 	{
	// 		Name: "register",
	// 		Description: "Add .ics calendar to register events with the bot. Be sure to attach a .ics file to this message!",
	// 		Type: discordgo.ChatApplicationCommand,
	// 	},
		
	// }

	// commandHandlers = map[string]func(s * discordgo.Session, i *discordgo.InteractionCreate) {
	// 	"register": handleEnrolment,
	// }

	commandHandlers = map[string]func(s * discordgo.Session, i *discordgo.MessageCreate) {
		"enrol": HandleEnrol,
		"whobusy": HandleWhoBusy,
		"wyd": HandleWyd,
	}
)