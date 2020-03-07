package main

import (
	"encoding/json"
	"fmt"
	"net/url"
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
		log := log.With(zap.String("id", profile.Graphql.User.ID), zap.String("username", profile.Graphql.User.Username))

		log.Info("visted profile")

		err = c.Visit(profile.Graphql.User.ProfilePicURLHd)
		if err != nil {
			log.Error("failure visiting profile pic", zap.Error(err))
			continue
		}

		if profile.Graphql.User.EdgeOwnerToTimelineMedia.PageInfo.HasNextPage {
			params := profilePostsParams{
				ID:    profile.Graphql.User.ID,
				First: 12,
				After: "",
			}
			paramsText, err := json.Marshal(params)
			if err != nil {
				log.Error("failure marshalling profile posts params", zap.Error(err))
				continue
			}

			graphQLPosts := `https://www.instagram.com/graphql/query/?query_hash=e769aa130647d2354c40ea6a439bfc08&variables=%s`

			err = c.Visit(fmt.Sprintf(graphQLPosts, url.QueryEscape(string(paramsText))))
			if err != nil {
				log.Error("failure visiting profile pic", zap.Error(err))
				continue
			}
		}
	}
}

type profilePage struct {
	EntryData struct {
		ProfilePage []struct {
			Graphql struct {
				User struct {
					Biography                string    `json:"biography"`
					ExternalURL              string    `json:"external_url"`
					EdgeFollowedBy           countEdge `json:"edge_followed_by"`
					EdgeFollow               countEdge `json:"edge_follow"`
					FullName                 string    `json:"full_name"`
					ID                       string    `json:"id"`
					IsVerified               bool      `json:"is_verified"`
					ProfilePicURLHd          string    `json:"profile_pic_url_hd"`
					Username                 string    `json:"username"`
					EdgeOwnerToTimelineMedia struct {
						Count    int `json:"count"`
						PageInfo struct {
							HasNextPage bool   `json:"has_next_page"`
							EndCursor   string `json:"end_cursor"`
						} `json:"page_info"`
						Edges []struct {
							Node PostNode `json:"node"`
						} `json:"edges"`
					} `json:"edge_owner_to_timeline_media"`
				} `json:"user"`
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
