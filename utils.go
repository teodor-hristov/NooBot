package main

import (
	"os/exec"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	CHATI      = "👋 Hey! "
	CHATE      = "👳‍♂️Ooops! "
	NEW_MEMBER = "👋 Glad to see you here "
)

type Pair struct {
	First  string
	Second uint64
}

func SortBySecond(users []Pair) {
	sort.SliceStable(users, func(i, j int) bool {
		return users[i].Second >= users[j].Second
	})
}

func GetChannel(guildId string, channelID string) *discordgo.Channel {
	gc, err := session.GuildChannels(guildId)
	if err != nil {
		return nil
	}

	for _, i := range gc {
		if i.Type == discordgo.ChannelTypeGuildVoice && strings.Compare(i.Name, channelID) == 0 {
			return i
		}
	}

	return nil
}

func GetSongName(url string) (string, error) {
	youtubedl := exec.Command("yt-dlp", "--get-title", url, "-o", "-")
	name, err := youtubedl.CombinedOutput()

	if err != nil {
		return "", err
	}

	defer youtubedl.Process.Kill()
	youtubedl.Wait()

	name = name[:len(name)-1]
	return string(name), nil
}
