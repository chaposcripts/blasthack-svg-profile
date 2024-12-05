package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Profile struct {
	ID         string
	Name       string
	Messages   string
	Solutions  string
	Topics     string
	Reactions  string
	Points     string
	GroupTitle string
	GroupTag   string
	AvatarURL  string
}

var GroupTags map[string]string = map[string]string{
	"Проверенный":          "verified",
	"Друг":                 "friend",
	"Модератор":            "moderator",
	"Всефрумный Модератор": "vmoderator",
	"Администратор":        "admin",
	"BH Team":              "bhteam",
}

var multiLanguageKeys map[string]string = map[string]string{
	"Messages":       "messages",
	"Solutions":      "solutions",
	"Threads":        "threads",
	"Reaction score": "reactions",
}

func GetProfileInfo(profileId string) (Profile, error) {
	var user Profile = Profile{
		ID:        profileId,
		Messages:  "0",
		Solutions: "0",
		Topics:    "0",
		Reactions: "0",
	}

	response, err := http.Get(fmt.Sprintf("https://www.blast.hk/members/%s/", profileId))
	if err != nil {
		return user, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return user, fmt.Errorf("HTTP_STATUS_%d", response.StatusCode)
	}

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return user, err
	}

	document.Find("h1").Each(func(i int, s *goquery.Selection) {
		user.Name = strings.TrimSpace(s.Text())
	})

	document.Find(".userBanner").Each(func(i int, s *goquery.Selection) {
		user.GroupTitle = s.Text()
		groupTag, tagExists := GroupTags[user.GroupTitle]
		if tagExists {
			user.GroupTag = groupTag
		}
	})

	document.Find(".avatar.avatar--l").Each(func(i int, s *goquery.Selection) {
		if len(user.AvatarURL) == 0 {
			user.AvatarURL = "https://www.blast.hk" + s.AttrOr("href", "nil")
		}
		fmt.Println("a", s.AttrOr("href", "nil"))
	})

	document.Find(".pairs.pairs--rows.pairs--rows--centered").Each(func(i int, s *goquery.Selection) {
		parts := strings.Fields(strings.TrimSpace(s.Text()))
		key, keyExists := multiLanguageKeys[parts[0]]
		if keyExists {
			if key == "messages" {
				user.Messages = parts[1]
			} else if key == "solutions" {
				user.Solutions = parts[1]
			} else if key == "threads" {
				user.Topics = parts[1]
			} else if key == "reactions" {
				user.Reactions = parts[1]
			}
		}
	})
	return user, nil
}
