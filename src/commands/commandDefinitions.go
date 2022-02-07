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

var commandHandlers = map[string]CommandHandler{
	"enrol":    HandleEnrol,
	"whobusy":  HandleWhoBusy,
	"wyd":      HandleWyd,
	"whenfree": HandleWhenFree,
}
