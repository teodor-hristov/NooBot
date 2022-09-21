package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Statistics struct {
	isTalking     bool
	secondsTalked uint64
	channelID     string
	guildId       string
	timeJoined    []time.Time
	timeLeft      []time.Time
}

var (
	voiceStats = make(map[string]map[string]Statistics)
	startTime  time.Time
)

func UserJoined(event *discordgo.VoiceStateUpdate) {
	if event == nil {
		return
	}

	userName := event.Member.User.Username
	entry := voiceStats[event.GuildID][userName]

	entry.isTalking = true
	entry.timeJoined = append(entry.timeJoined, time.Now())
	entry.channelID = event.ChannelID

	voiceStats[event.GuildID][userName] = entry

	fmt.Printf("User joined %s %s\n", userName, voiceStats[event.GuildID][userName].timeJoined[len(voiceStats[event.GuildID][userName].timeJoined)-1].String())
}

func UserLeft(event *discordgo.VoiceStateUpdate) {
	if event == nil {
		return
	}

	userName := event.Member.User.Username
	entry := voiceStats[event.GuildID][userName]

	entry.isTalking = false
	entry.timeLeft = append(entry.timeLeft, time.Now())

	//the only reason for len(voiceStats[event.GuildID][userName].timeJoined) to be 0 is when bot is after the other person (who left)
	if len(voiceStats[event.GuildID][userName].timeJoined) == 0 {
		entry.timeJoined = append(entry.timeJoined, startTime)
	}
	//entry.channelID = event.ChannelID on leave no channelID is given
	entry.secondsTalked += uint64(time.Now().Sub(entry.timeJoined[len(entry.timeJoined)-1]).Abs().Seconds())

	voiceStats[event.GuildID][userName] = entry

	//if bot is disconnected while playing -> turn off the music
	if strings.Compare(userName, session.State.User.Username) == 0 {
		entry := mp[event.GuildID]
		entry.isPlaying = false
		mp[event.GuildID] = entry
	}

	fmt.Printf("User left voice channel %s %s\nTotal talking time %d seconds\n",
		userName,
		voiceStats[event.GuildID][userName].timeLeft[len(voiceStats[event.GuildID][userName].timeLeft)-1].String(),
		voiceStats[event.GuildID][userName].secondsTalked)
}

func UserVoiceHistory(guildId string, userName string) string {
	var res strings.Builder

	arrLen := len(voiceStats[guildId][userName].timeJoined) + len(voiceStats[guildId][userName].timeLeft)
	join := 0
	left := 0

	for i := 0; i < arrLen; i++ {
		res.WriteString("\n")
		if len(voiceStats[guildId][userName].timeJoined) > join && voiceStats[guildId][userName].timeJoined[join].Before(voiceStats[guildId][userName].timeLeft[left]) {
			res.WriteString("**Joined at:** ")
			res.WriteString(voiceStats[guildId][userName].timeJoined[join].String())
			join++

		} else {
			res.WriteString("**Left at:** ")
			res.WriteString(voiceStats[guildId][userName].timeLeft[left].String())
			left++
		}

	}

	return res.String()
}
