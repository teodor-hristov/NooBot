package main

import (
	"strings"
	"testing"
)

func TestGetSongName(t *testing.T) {
	cases := []string{"", "12345678", "alabalaportokala",
		"https://www.youtube.com/watch?v=-D9XcT9N73Y", "www.google.com"}

	expected := []string{"", "", "", "Qvkata DLG - Extra , a ti ?", ""}

	for i := 0; i < len(cases); i++ {
		res, _ := GetSongName(cases[i])

		if strings.Compare(res, expected[i]) != 0 {
			t.Errorf("Case %d failed", i)
		}
	}
}

func TestIsValidURL(t *testing.T) {
	cases := []string{"", "alabalaportokala",
		"https://www.youtube.com/watch?v=-D9XcT9N73Y", "www.google.com",
		"http://alabala.com", "http://", "wwww.google.bg"}

	expected := []bool{false, false, true, false, true, false, false}

	for i := 0; i < len(cases); i++ {
		res := IsValidURL(cases[i])

		if res != expected[i] {
			t.Errorf("Case %d failed", i)
		}
	}
}
