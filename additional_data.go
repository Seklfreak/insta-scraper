package main

import (
	"encoding/json"
	"strings"

	"github.com/gocolly/colly"
	"go.uber.org/zap"
)

func crawlAdditionalData(e *colly.HTMLElement) {
	jsonText := e.Text[strings.Index(e.Text, "{") : strings.LastIndex(e.Text, "}")+1]

	var data additionalData
	err := json.Unmarshal([]byte(jsonText), &data)
	if err != nil {
		log.Error("failure unmarshalling profile page", zap.Error(err))
		return
	}

	crawlPost(&data.Graphql.ShortcodeMedia)
}

type additionalData struct {
	Graphql struct {
		ShortcodeMedia PostNode `json:"shortcode_media"`
	} `json:"graphql"`
}
