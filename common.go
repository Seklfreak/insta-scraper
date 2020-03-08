package main

import (
	"strings"

	"go.uber.org/zap"
)

func visit(page string) {
	if !cfg.DownloadMedia {
		if strings.Contains(page, ".jpg") || strings.Contains(page, ".mp4") {
			return
		}
	}

	err := c.Visit(page)
	if err != nil {
		log.Error("failure visiting page", zap.String("page", page), zap.Error(err))
	}
}
