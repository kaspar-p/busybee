package commands

// SLASH COMMAND CODE
// createdCommands []*discordgo.ApplicationCommand
// commands = []*discordgo.ApplicationCommand {
// 	{
// 		Name: "register",
// 		Description: "Add .ics calendar to register events with the bot. Be sure to attach a .ics file to this message!",
// 		Type: discordgo.ChatApplicationCommand,
// 	},.

// }.

// commandHandlers = map[string]func(s * discordgo.Session, i *discordgo.InteractionCreate) {
// 	"register": handleEnrolment,
// }.

type Command struct {
	Handler CommandHandler
	Trigger string
}

var (
	ENROL    = Command{HandleEnrol, "enrol"}
	WHOBUSY  = Command{HandleWhoBusy, "whobusy"}
	WYD      = Command{HandleWyd, "wyd"}
	WHENFREE = Command{HandleWhenFree, "whenfree"}
)

var commandHandlers = []Command{ENROL, WHOBUSY, WYD, WHENFREE}
