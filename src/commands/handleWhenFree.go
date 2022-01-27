package commands

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/entities"
	"github.com/kaspar-p/bee/src/utils"
)

type TimePair struct {
	Hour     int
	TimeText string
}

func validateStructure(message *discordgo.MessageCreate) string {
	if len(message.Mentions) != len(strings.Split(message.Content, " "))-1 {
		return "message not of the form: `.whenfree <@ mention> <@ mention> ...` \\:("
	}

	if len(message.Mentions) == 0 {
		return "message must have at least one mention \\:("
	}

	return ""
}

func getNextFreeIntervalOfSize(user *entities.User, startingAt time.Time, numberOfHours int) (time.Time, bool) {
	for i := 1; i < len(user.BusyTimes); i++ {
		currentBusyTime := user.BusyTimes[i]

		if currentBusyTime.Start.After(startingAt) {
			previousBusyTime := user.BusyTimes[i-1]

			var intervalLength int
			if previousBusyTime.End.Before(startingAt) {
				// If we are currently in the interval, only use the interval of [startingAt, currentBusyTime.Start]
				intervalLength = int(math.Floor(currentBusyTime.Start.Sub(startingAt).Hours()))
			} else {
				// If we are NOT currently in the interval, use the whole interval [previousBusyTime.End, currentBusyTime.Start]
				intervalLength = int(math.Floor(currentBusyTime.Start.Sub(previousBusyTime.End).Hours()))
			}

			if intervalLength >= numberOfHours {
				return previousBusyTime.End, true
			}
		}
	}

	// Now() as a standin time.Time, not relevant
	return time.Now(), false
}

func getLaterTime(time1, time2 time.Time) time.Time {
	switch {
	case time1.Equal(time2):
		return time1
	case time1.After(time2):
		return time1
	default:
		return time2
	}
}

func getNextCommonFreeNumberOfHoursTwo(user1, user2 *entities.User, numberOfHours int) (time.Time, bool) {
	latestEndTime := getLaterTime(user1.GetLatestEndTime(), user2.GetLatestEndTime())
	numberOfHoursUntilLatestEndTime := int(time.Until(latestEndTime).Hours())

	for hourOffset := 0; hourOffset < numberOfHoursUntilLatestEndTime; hourOffset++ {
		startTime := time.Now().Add(time.Duration(hourOffset) * time.Hour)
		user1NextStartFreeTime, foundStart1 := getNextFreeIntervalOfSize(user1, startTime, numberOfHours)
		user2NextStartFreeTime, foundStart2 := getNextFreeIntervalOfSize(user2, startTime, numberOfHours)

		log.Printf(
			"# hours: %d. # hour offset: %d. user1 next free start time %v, user 2 next free start time %v\n",
			numberOfHours,
			numberOfHours,
			user1NextStartFreeTime.Format(time.Layout),
			user2NextStartFreeTime.Format(time.Layout),
		)

		// Exit if they are not free at all for this number of hours
		if !foundStart1 || !foundStart2 {
			return time.Now(), false
		}

		// If they are both equal at the same time
		if user1NextStartFreeTime.Equal(user2NextStartFreeTime) {
			return user2NextStartFreeTime, true
		} else if (user1NextStartFreeTime.After(user2NextStartFreeTime)) ||
			(user2NextStartFreeTime.After(user1NextStartFreeTime)) {
			// If the structure is (intervals end anywhere, length >= numberOfHours):
			// .....[........time2............
			// .........[....time1............
			// or
			// .........[....time2............
			// .....[........time1............
			laterStartTime := getLaterTime(user1NextStartFreeTime, user2NextStartFreeTime)

			return laterStartTime, true
		}
	}

	return time.Now(), false
}

func getNextCommonFreeNumberOfHoursMany(users []*entities.User, numberOfHours int) (time.Time, bool) {
	var latestFreeCommonTime time.Time

	set := false

	if len(users) == 1 {
		return getNextFreeIntervalOfSize(users[0], time.Now(), numberOfHours)
	}

	for i := 0; i < len(users)-1; i++ {
		for j := i + 1; j < len(users); j++ {
			userI := users[i]
			userJ := users[j]

			nextCommonFreeTime, found := getNextCommonFreeNumberOfHoursTwo(userI, userJ, numberOfHours)
			if !found {
				return time.Now(), false
			}

			if !set || latestFreeCommonTime.Before(nextCommonFreeTime) {
				set = true
				latestFreeCommonTime = nextCommonFreeTime
			}
		}
	}

	return latestFreeCommonTime, set
}

func toNiceDateTimeString(eventTime time.Time) string {
	return eventTime.Format("03:04 PM 01/02")
}

func HandleWhenFree(discord *discordgo.Session, message *discordgo.MessageCreate) error {
	errorMessage := validateStructure(message)
	if errorMessage != "" {
		log.Println("Command .whenFree error with message:", errorMessage)
		err := SendSingleMessage(discord, message.ChannelID, errorMessage)

		return err
	}

	// Convert mentions into 'User's
	mentionedUsers := make([]*entities.User, 0)

	for _, mentionedUser := range message.Mentions {
		if user, ok := entities.Users[message.GuildID][mentionedUser.ID]; ok {
			log.Println("Got mentioned user", user.Name)
			mentionedUsers = append(mentionedUsers, user)
		} else {
			err := SendSingleMessage(
				discord,
				message.ChannelID,
				"the @ mentioned user `"+mentionedUser.Username+"` isn't in the system. ask them to enrol pls \\:)",
			)

			return err
		}
	}

	maxHours := 6

	timePairs := make([]TimePair, maxHours)

	for hour := 1; hour < maxHours+1; hour++ {
		timeFound, found := getNextCommonFreeNumberOfHoursMany(mentionedUsers, hour)

		var timeText string
		if !found {
			timeText = "NONE \\:("
		} else {
			timeText = toNiceDateTimeString(timeFound)
		}

		timePairs[hour-1] = TimePair{
			Hour:     hour,
			TimeText: timeText,
		}
	}

	embed := GenerateWhenFreeEmbed(timePairs)
	err := SendSingleEmbed(discord, message.ChannelID, embed)

	return err
}

func GenerateWhenFreeEmbed(timePairs []TimePair) *discordgo.MessageEmbed {
	var bodyString string

	for _, pair := range timePairs {
		hourText := "s"
		spaceText := ""

		if pair.Hour == 1 {
			hourText = ""
			spaceText = " "
		}

		bodyString += fmt.Sprintf("%d hour%s:%s %s\n", pair.Hour, hourText, spaceText, pair.TimeText)
	}

	descriptionString := utils.WrapStringInCodeBlock(bodyString)

	return CreateTableEmbed("hours fwee \\:)", descriptionString)
}
