package main

import (
	"testing"
)

func TestGetUserVoiceHistory(t *testing.T) {
	st1 := Statistics{
		isTalking:     false,
		secondsTalked: 0,
		guildId:       "somerndgid",
	}

	//it is possible to left before join
	st1.UserLeft()
	//len(st1.timeJoined) == 0 is needed because if there is no join times we need to add bot's join time
	if st1.isTalking || len(st1.timeLeft) == 0 || len(st1.timeJoined) == 0 {
		t.Error("Update failed when user joined!")
	}

	st2 := Statistics{
		isTalking:     false,
		secondsTalked: 0,
		guildId:       "somerndgid",
	}

	//no need to test multuple joins because this is impossible situation on daily basis
	st2.UserJoined()
	if !st2.isTalking || len(st2.timeJoined) != 1 {
		t.Error("Update failed when user joined!")
	}

	st2.UserLeft()
	if st2.isTalking || len(st2.timeLeft) == 0 {
		t.Error("Update failed when user joined!")
	}
}
