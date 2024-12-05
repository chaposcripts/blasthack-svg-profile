package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Profile struct {
	ID           string
	Name         string
	Messages     string
	Solutions    string
	Topics       string
	Reactions    string
	GroupTitle   string
	GroupTag     string
	AvatarURL    string
	AvatarBase64 string
}

var GroupTags map[string]string = map[string]string{
	"Проверенный": "verified",
	"Друг":        "friend",
	"Модератор":   "moderator",
	"Всефорумный модератор": "vmoderator",
	"Администратор":         "admin",
	"BH Team":               "bhteam",
}

var multiLanguageKeys map[string]string = map[string]string{
	"Сообщения": "messages",
	"Решения":   "solutions",
	"Темы":      "threads",
	"Реакции":   "reactions",
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
			avatarPath := s.AttrOr("href", "")
			if len(avatarPath) > 0 {
				user.AvatarURL = "https://www.blast.hk" + avatarPath
				user.AvatarBase64, err = DownloadAvatar(user.AvatarURL)
				if err != nil {
					fmt.Println("Error loading user avatar as base64", err.Error())
				}
			}
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
		} else {
			fmt.Println("key", parts[0], "not found")
		}
	})
	return user, nil
}

func DownloadAvatar(avatarURL string) (string, error) {
	var result string
	response, err := http.Get(avatarURL)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	imageData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении данных изображения:", err)
		return result, err
	}

	result = "data:image/png;base64," + base64.StdEncoding.EncodeToString(imageData)
	// fmt.Println("User profile avatar base64:", result)
	return result, nil
}
