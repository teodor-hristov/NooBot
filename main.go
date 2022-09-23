package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	BotToken = flag.String("token", "", "Bot access token")
	session  *discordgo.Session
)

func addIntents() {
	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
}

func addHandlers() {
	addHandlersToCommands()

	//when someone enters voice chat
	session.AddHandler(voiceStatusUpdate)
}

func voiceStatusUpdate(session *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	if session == nil || event == nil {
		return
	}

	userStat := voiceStats[event.GuildID][event.Member.User.Username]
	if event.ChannelID == "" {
		userStat.UserLeft()

		//if bot is disconnected while playing -> turn off the music
		if strings.Compare(event.Member.User.Username, session.State.User.Username) == 0 {
			userStat := mp[event.GuildID]
			userStat.isPlaying = false
			mp[event.GuildID] = userStat
		}

		fmt.Printf("User left voice channel %s %s\nTotal talking time %d seconds\n", event.Member.User.Username, userStat.timeLeft[len(userStat.timeLeft)-1].String(),
			userStat.secondsTalked)
		voiceStats[event.GuildID][event.Member.User.Username] = userStat
		return
	}

	userStat.UserJoined()
	fmt.Printf("User joined %s %s\n", event.Member.User.Username, userStat.timeJoined[len(userStat.timeJoined)-1].String())

	voiceStats[event.GuildID][event.Member.User.Username] = userStat
	return
}

func connectToDiscord() {
	//Connect to discord
	var err error
	session, err = discordgo.New("Bot" + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
		return
	}
	fmt.Print("Instance created.\n")

	err = session.Open()
	if err != nil {
		log.Fatalf("Connection open error: %v", err)
		return
	}
	fmt.Print("Connection created.\n")
	startTime = time.Now()
}

func disconnectFromDiscord() {
	err := session.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Print("Connection closed.\n")
}

func main() {
	flag.Parse()

	connectToDiscord()

	for _, item := range session.State.Guilds {
		voiceStats[item.ID] = make(map[string]Statistics)
	}

	addHandlers()
	addIntents()

	createCommands()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("----------------------------\nBot is now running.  Press CTRL-C to exit.\n----------------------------")

	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-stopBot

	// Cleanly close down the Discord session.
	disconnectFromDiscord()
}
