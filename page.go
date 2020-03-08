package main

import (
	"encoding/json"
	"strings"

	"github.com/gocolly/colly"
	"go.uber.org/zap"
)

func crawlPage(e *colly.HTMLElement) {
	jsonText := e.Text[strings.Index(e.Text, "{") : strings.LastIndex(e.Text, "}")+1]

	var profile profilePage
	err := json.Unmarshal([]byte(jsonText), &profile)
	if err != nil {
		log.Error("failure unmarshalling profile page", zap.Error(err))
		return
	}

	for _, profile := range profile.EntryData.ProfilePage {
		crawlUser(&profile.Graphql.User)
	}
}

type profilePage struct {
	EntryData struct {
		ProfilePage []struct {
			Graphql struct {
				User UserNode `json:"user"`
			} `json:"graphql"`
		} `json:"ProfilePage"`
	} `json:"entry_data"`
}

type profilePostsParams struct {
	ID    string `json:"id"`
	First int    `json:"first"`
	After string `json:"after"`
}

type countEdge struct {
	Count int `json:"count"`
}
