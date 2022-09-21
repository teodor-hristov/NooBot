package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "voice-chat-top",
			Description: "Print users who talked in the server for all times.",
		},
		{
			Name:        "talk-history",
			Description: "Command gives voice chat history for a given user.",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "user",
					Description: "Username",
					Required:    true,
				},
			},
		},
		{
			Name:        "play",
			Description: "IntroToGo bot starts playing youtube song from given link.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "channel",
					Description: "Voice channel for the bot to enter",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "url",
					Description: "Youtube song url",
					Required:    true,
				},
			},
		},
		{
			Name:        "download",
			Description: "IntroToGo bot starts downloading youtube song from given link and sends it to chat.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "url",
					Description: "Youtube/Soundcloud song url",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "format",
					Description: "Output format for the song",
					Required:    true,
				},
			},
		},
		{
			Name:        "stop",
			Description: "If there is music playing the command stops it and then disconnects.",
		},
	}

	commandsHandler = map[string]func(session *discordgo.Session, event *discordgo.InteractionCreate){
		"voice-chat-top": voiceChatTopCommand,
		"talk-history":   talkHistoryCommand,
		"play":           playCommand,
		"stop":           stopCommand,
		"download":       downloadSongCommand,
	}
)

func stopCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var sb strings.Builder
	player := mp[i.GuildID]

	if s == nil || i == nil {
		return
	}

	if player == nil || !player.isPlaying {
		sb.WriteString(CHATI + "🔸There is no song playing!\n")
	} else {
		sb.WriteString(CHATI + "🔸Song stopped.\n")
		player.isPlaying = false
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: sb.String(),
		},
	})
}

func playCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if s == nil || i == nil {
		return
	}

	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var sb strings.Builder

	if mp[i.GuildID] != nil && mp[i.GuildID].isPlaying {
		sb.WriteString(CHATE + "🔸Song is already going!\n")
	}

	// When the option exists, ok = true
	channelopt, ok := optionMap["channel"]
	if !ok {
		sb.WriteString(CHATE + "🔸Not valid channel!\n")
	}

	urlopt, ok := optionMap["url"]
	if !ok {
		sb.WriteString(CHATE + "🔸Not valid url!\n")
	}

	channel := GetChannel(i.GuildID, channelopt.StringValue())
	if channel == nil {
		sb.WriteString(CHATE + "🔸Channel not found")
	}

	if sb.Len() == 0 {
		sb.WriteString("🔊  Playing: ")
		sb.WriteString(urlopt.StringValue())
		sb.WriteString("\n")
		go func() {
			err := discordPlayMusic(i.GuildID, channel.ID, urlopt.StringValue())
			if err != nil {
				fmt.Print(err)
				session.ChannelMessageSend(i.ChannelID, CHATE+" Something went wrong!")
			}

			mp[i.GuildID].voiceConn.Disconnect()
		}()
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: sb.String(),
		},
	})
}

func downloadSongCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if s == nil || i == nil {
		return
	}

	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var sb strings.Builder

	// When the option exists, ok = true
	formatopt, ok := optionMap["format"]
	if !ok {
		sb.WriteString(CHATE + "🔸Not valid channel!\n")
	}

	urlopt, ok := optionMap["url"]
	if !ok {
		sb.WriteString(CHATE + "🔸Not valid url!\n")
	}

	if sb.Len() == 0 {
		sb.WriteString("🎥 Downloading wanted song.\nI'll send download link when finished!👽\n")
		go func() {
			err := downloadSong(i.ChannelID, urlopt.StringValue(), formatopt.StringValue())
			if err != nil {
				fmt.Print(err)
				session.ChannelMessageSend(i.ChannelID, CHATE+" Something went wrong... 👀")
			}
		}()
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: sb.String(),
		},
	})
}

func talkHistoryCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if s == nil || i == nil {
		return
	}
	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var sb strings.Builder

	// When the option exists, ok = true
	if option, ok := optionMap["user"]; ok {
		sb.WriteString("Voice history for user: ")
		sb.WriteString(option.StringValue())
		sb.WriteString("\n")
		sb.WriteString(UserVoiceHistory(i.GuildID, option.StringValue()))
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: sb.String(),
		},
	})
}

func voiceChatTopCommand(session *discordgo.Session, event *discordgo.InteractionCreate) {
	if session == nil || event == nil {
		return
	}
	var sb strings.Builder
	if len(voiceStats[event.GuildID]) > 0 {
		top := make([]Pair, len(voiceStats))

		i := 0
		for key, value := range voiceStats[event.GuildID] {
			top[i] = Pair{key, value.secondsTalked}
			i++
		}

		SortBySecond(top)

		for count, item := range top {
			if strings.Compare(item.First, session.State.User.Username) != 0 {
				fmt.Fprintf(&sb, "%d : %s	%d seconds\n", count+1, item.First, item.Second)
			}
		}
	} else {
		sb.WriteString("No users talked.")
	}

	session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: sb.String(),
		},
	})
}

func createCommands() {
	len := len(commands)

	fmt.Printf("Register %d commands.\n", len)

	for _, item := range session.State.Guilds {
		for i := 0; i < len; i++ {
			ccmd, err := session.ApplicationCommandCreate(session.State.User.ID, item.ID, commands[i])
			if err != nil {
				fmt.Printf("%s failed registering with error %d.", ccmd.Name, err)
			}
		}
	}

}

func addHandlersToCommands() {
	len := len(commandsHandler)
	fmt.Printf("Add %d command handlers.\n", len)

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if cmd, ok := commandsHandler[i.ApplicationCommandData().Name]; ok {
			cmd(s, i)
		}
	})
}
