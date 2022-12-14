package main

import (
	"errors"
	"net/url"
	"os/exec"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	CHATI      = "ð Hey! "
	CHATE      = "ð³ââï¸Ooops! "
	NEW_MEMBER = "ð Glad to see you here "
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
	if !IsValidURL(url) {
		return "", errors.New("Invalid input!")
	}

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

func IsValidURL(songUrl string) bool {
	u, err := url.ParseRequestURI(songUrl)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}
