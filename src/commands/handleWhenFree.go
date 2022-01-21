package commands

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kaspar-p/bee/src/entities"
)

func validateStructure(message *discordgo.MessageCreate) string {
	if (len(message.Mentions) != len(strings.Split(message.Content, " ")) - 1) {
		return "message not of the form: `.whenfree <@ mention> <@ mention> ...` \\:("
	}

	if (len(message.Mentions) == 0) {
		return "message must have at least one mention \\:(";
	} 

	return ""
}

func getNextFreeIntervalOfSize(user *entities.User, startingAt time.Time, numberOfHours int) (time.Time, bool) {
	for i := 1; i < len(user.BusyTimes); i++ {
		currentBusyTime := user.BusyTimes[i];
		
		if currentBusyTime.Start.After(startingAt) {
			previousBusyTime := user.BusyTimes[i-1];
			
			var intervalLength int;
			if previousBusyTime.End.Before(startingAt) {
				// If we are currently in the interval, only use the interval of [startingAt, currentBusyTime.Start]
				intervalLength = int(math.Floor(currentBusyTime.Start.Sub(startingAt).Hours()));
			} else {
				// If we are NOT currently in the interval, use the whole interval [previousBusyTime.End, currentBusyTime.Start]
				intervalLength = int(math.Floor(currentBusyTime.Start.Sub(previousBusyTime.End).Hours()));
			}
			
			if intervalLength >= numberOfHours {
				return previousBusyTime.End, true;
			}
		}
	}

	// Now() as a standin time.Time, not relevant
	return time.Now(), false;
}

func getLaterTime(t1 time.Time, t2 time.Time) time.Time {
	if (t1.Equal(t2)) {
		return t1;
	} else if (t1.After(t2)) {
		return t1;
	} else {
		return t2;
	}
}

func getNextCommonFreeNumberOfHoursTwo(user1 *entities.User, user2 *entities.User, numberOfHours int) (time.Time, bool) {
	latestEndTime := getLaterTime(user1.GetLatestEndTime(), user2.GetLatestEndTime());
	numberOfHoursUntilLatestEndTime := int(time.Until(latestEndTime).Hours());

	for hourOffset := 0; hourOffset < numberOfHoursUntilLatestEndTime; hourOffset++ {
		startTime := time.Now().Add(time.Duration(hourOffset) * time.Hour);
		user1NextStartFreeTime, foundStart1 := getNextFreeIntervalOfSize(user1, startTime, numberOfHours);
		user2NextStartFreeTime, foundStart2 := getNextFreeIntervalOfSize(user2, startTime, numberOfHours);

		fmt.Printf("# hours: %d. # hour offset: %d. user1 next free start time %v, user 2 next free start time %v\n", numberOfHours, numberOfHours, user1NextStartFreeTime.Format(time.Layout), user2NextStartFreeTime.Format(time.Layout))

		// Exit if they are not free at all for this number of hours
		if !foundStart1 || !foundStart2 {
			return time.Now(), false;
		}

		// If they are both equal at the same time
		if user1NextStartFreeTime.Equal(user2NextStartFreeTime) {
			return user2NextStartFreeTime, true;
		} else if 	(user1NextStartFreeTime.After(user2NextStartFreeTime)) ||
					(user2NextStartFreeTime.After(user1NextStartFreeTime)) {
			// If the structure is (intervals end anywhere, length >= numberOfHours):
			// .....[........time2............
			// .........[....time1............
			// or
			// .........[....time2............
			// .....[........time1............
			laterStartTime := getLaterTime(user1NextStartFreeTime, user2NextStartFreeTime);
			return laterStartTime, true;
		}
	}

	return time.Now(), false;
}

func getNextCommonFreeNumberOfHoursMany(users []*entities.User, numberOfHours int) (time.Time, bool) {
	var latestFreeCommonTime time.Time;
	set := false;

	if len(users) == 1 {
		return getNextFreeIntervalOfSize(users[0], time.Now(), numberOfHours);
	}

	for i := 0; i < len(users) - 1; i++ {
		for j := i+1; j < len(users); j++ {
			userI := users[i];
			userJ := users[j];

			nextCommonFreeTime, found := getNextCommonFreeNumberOfHoursTwo(userI, userJ, numberOfHours);
			if !found {
				return time.Now(), false;
			}

			if !set || latestFreeCommonTime.Before(nextCommonFreeTime) {
				set = true;
				latestFreeCommonTime = nextCommonFreeTime;
			}
		}
	}

	return latestFreeCommonTime, set;
}


func toNiceDateTimeString(eventTime time.Time) string {
	return eventTime.Format("03:04 PM 01/02");
}

func HandleWhenFree(discord *discordgo.Session, message *discordgo.MessageCreate) {
	errorMessage := validateStructure(message)
	if errorMessage != "" {
		fmt.Println("Command .whenFree error with message:", errorMessage);
		discord.ChannelMessageSend(message.ChannelID, errorMessage);
		return
	}

	// Convert mentions into 'User's
	mentionedUsers := make([]*entities.User, 0);
	for _, mentionedUser := range message.Mentions {
		if user, ok := entities.Users[message.GuildID][mentionedUser.ID]; ok {
			fmt.Println("Got mentioned user", user.Name);
			mentionedUsers = append(mentionedUsers, user);
		} else {
			discord.ChannelMessageSend(message.ChannelID, "the @ mentioned user `" + mentionedUser.Username + "` isn't in the system. ask them to enrol pls \\:)");
			return;
		}
	}

	resultString := "```\n"
	resultString += "+--------------------------+\n";
	maxHours := 12;

	type TimePair struct {
		Time time.Time
		Found bool
	}

	hourTimesMap := make(map[int]TimePair);
	for hour := 1; hour < maxHours+1; hour++ {
		t, found := getNextCommonFreeNumberOfHoursMany(mentionedUsers, hour);
		if found {
			hourTimesMap[hour] = TimePair{t, found};
		}
	}

	for hour, timePair := range hourTimesMap {
		var hourText string;
		if hour == 1 {
			hourText = "hour ";
		} else {
			hourText = "hours";
		}

		if !timePair.Found {
			resultString += fmt.Sprintf("| %d %s | %s |\n", hour, hourText, "    NONE \\:(   ");
		} else {
			resultString += fmt.Sprintf("| %d %s | %s |\n", hour, hourText, toNiceDateTimeString(timePair.Time));
		}
	}
	resultString += "+--------------------------+\n";
	resultString += "```\n";
	discord.ChannelMessageSend(message.ChannelID, resultString);
}