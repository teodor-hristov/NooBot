package main

import (
	"sort"
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

	//sort all times
	allTimes := append(stat.timeJoined, stat.timeLeft...)
	sort.Slice(allTimes, func(i, j int) bool {
		return allTimes[i].Before(allTimes[j])
	})

	i := 0
	if len(stat.timeJoined) < len(stat.timeLeft) {
		res.WriteString("⛔️ **Left at:** " + allTimes[i].String() + "\n")
		i++
	}

	odd := true
	for ; i < len(allTimes); i++ {
		if odd {
			res.WriteString("✅ **Joined at:** " + allTimes[i].String() + "\n")
			odd = false
			continue
		}

		res.WriteString("⛔️ **Left at:** " + allTimes[i].String() + "\n")
		odd = true
	}

	return res.String()
}
