package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	"go.uber.org/zap"
)

func crawlProfilePostsGraphQL(sourceURL *url.URL, data []byte) {
	if sourceURL == nil {
		return
	}

	var params profilePostsParams
	err := json.Unmarshal([]byte(sourceURL.Query().Get("variables")), &params)
	if err != nil {
		log.Error("failure unmarshalling source URL params", zap.Error(err))
		return
	}

	var posts graphQLPosts
	err = json.Unmarshal(data, &posts)
	if err != nil {
		log.Error("failure unmarshalling graphql posts", zap.Error(err))
		return
	}

	for _, post := range posts.Data.User.EdgeOwnerToTimelineMedia.Edges {
		crawlPost(&post.PostNode)
	}

	if posts.Data.User.EdgeOwnerToTimelineMedia.PageInfo.HasNextPage {
		params.After = posts.Data.User.EdgeOwnerToTimelineMedia.PageInfo.EndCursor
		paramsText, err := json.Marshal(params)
		if err != nil {
			log.Error("failure marshalling profile posts params", zap.Error(err))
			return
		}

		graphQLPosts := `https://www.instagram.com/graphql/query/?query_hash=e769aa130647d2354c40ea6a439bfc08&variables=%s`

		err = c.Visit(fmt.Sprintf(graphQLPosts, url.QueryEscape(string(paramsText))))
		if err != nil {
			log.Error("failure visiting profile pic", zap.Error(err))
			return
		}
	}
}

type graphQLPosts struct {
	Data struct {
		User struct {
			EdgeOwnerToTimelineMedia struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool   `json:"has_next_page"`
					EndCursor   string `json:"end_cursor"`
				} `json:"page_info"`
				Edges []struct {
					PostNode `json:"node"`
				} `json:"edges"`
			} `json:"edge_owner_to_timeline_media"`
		} `json:"user"`
	} `json:"data"`
	Status string `json:"status"`
}
