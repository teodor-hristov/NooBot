package main

import (
	"strings"
	"time"
)

var (
	voiceStats = make(map[string]map[string]Statistics) //guild -> username -> stat
	startTime  time.Time
)

type Statistics struct {
	isTalking     bool
	secondsTalked uint64
	guildId       string
	timeJoined    []time.Time
	timeLeft      []time.Time
}

func (stat *Statistics) UserJoined() {
	if stat == nil {
		return
	}

	stat.isTalking = true
	stat.timeJoined = append(stat.timeJoined, time.Now())
}

func (stat *Statistics) UserLeft() {
	if stat == nil {
		return
	}

	stat.isTalking = false
	stat.timeLeft = append(stat.timeLeft, time.Now())

	//the only reason for len(stat.timeJoined) to be 0 is when bot is after the other person (who left)
	if len(stat.timeJoined) == 0 {
		stat.timeJoined = append(stat.timeJoined, startTime)
	}

	stat.secondsTalked += uint64(time.Now().Sub(stat.timeJoined[len(stat.timeJoined)-1]).Abs().Seconds())
}

func (stat Statistics) GetUserVoiceHistory() string {
	var res strings.Builder

	arrLen := len(stat.timeJoined) + len(stat.timeLeft)
	join := 0
	left := 0

	for i := 0; i < arrLen; i++ {
		res.WriteString("\n")
		if len(stat.timeJoined) > join && stat.timeJoined[join].Before(stat.timeLeft[left]) {
			res.WriteString("**Joined at:** ")
			res.WriteString(stat.timeJoined[join].String())
			if len(stat.timeJoined)-1 != join {
				join++
			}

		} else {
			res.WriteString("**Left at:** ")
			res.WriteString(stat.timeLeft[left].String())
			if len(stat.timeLeft)-1 != left {
				left++
			}
		}

	}

	return res.String()
}
