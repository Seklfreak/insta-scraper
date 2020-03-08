package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	"go.uber.org/zap"
)

func crawlUser(user *UserNode) {
	log := log.With(zap.String("id", user.ID), zap.String("username", user.Username))

	if !visitUser(user.ID) {
		return
	}

	log.Info("visted profile")

	if user.ProfilePicURLHd == "" {
		visit(fmt.Sprintf("https://www.instagram.com/%s/", user.Username))
		return
	}

	visit(user.ProfilePicURLHd)

	if user.EdgeOwnerToTimelineMedia.PageInfo.HasNextPage {
		params := profilePostsParams{
			ID:    user.ID,
			First: 12,
			After: "",
		}
		paramsText, err := json.Marshal(params)
		if err != nil {
			log.Error("failure marshalling profile posts params", zap.Error(err))
			return
		}

		graphQLPosts := `https://www.instagram.com/graphql/query/?query_hash=e769aa130647d2354c40ea6a439bfc08&variables=%s`

		visit(fmt.Sprintf(graphQLPosts, url.QueryEscape(string(paramsText))))
	}
}

type UserNode struct {
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
}
