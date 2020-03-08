package main

import (
	"fmt"

	"go.uber.org/zap"
)

func crawlPost(post *PostNode) {
	log := log.With(
		zap.String("id", post.ID),
		zap.String("shortcode", post.Shortcode),
		zap.String("type", post.Typename),
	)

	switch post.Typename {
	case "GraphImage":
		crawlDisplayResources(post.DisplayResources)
	case "GraphVideo":
		if post.VideoURL != "" {
			visit(post.VideoURL)
		} else if post.Shortcode != "" {
			visit(fmt.Sprintf("https://www.instagram.com/p/%s/", post.Shortcode))
		} else {
			log.Warn("post without video URL nor shortcode, do not know how to handle")
		}
	case "GraphSidecar":
		for _, child := range post.EdgeSidecarToChildren.Edges {
			crawlPost(&child.Node)
		}
	default:
		log.Warn("unknown post type")
	}

	for _, taggedUser := range post.EdgeMediaToTaggedUser.Edges {
		crawlUser(&taggedUser.Node.User)
	}
}

type PostNode struct {
	Typename              string `json:"__typename"`
	ID                    string `json:"id"`
	Shortcode             string `json:"shortcode"`
	EdgeSidecarToChildren struct {
		Edges []struct {
			Node PostNode `json:"node"`
		} `json:"edges"`
	} `json:"edge_sidecar_to_children"`
	DisplayResources      displayResources `json:"display_resources"`
	VideoURL              string           `json:"video_url"`
	EdgeMediaToTaggedUser struct {
		Edges []struct {
			Node struct {
				User UserNode `json:"user"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"edge_media_to_tagged_user"`
}

func crawlDisplayResources(resources displayResources) {
	var lastSrc string
	var lastWidth, lastHeight int

	for _, resource := range resources {
		if lastSrc == "" ||
			resource.ConfigHeight > lastHeight ||
			resource.ConfigWidth > lastWidth {
			lastSrc = resource.Src
			lastWidth = resource.ConfigWidth
			lastHeight = resource.ConfigHeight
		}
	}

	visit(lastSrc)
}

type displayResources []struct {
	Src          string `json:"src"`
	ConfigWidth  int    `json:"config_width"`
	ConfigHeight int    `json:"config_height"`
}
