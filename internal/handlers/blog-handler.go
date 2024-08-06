package handlers

import "goblog/internal/configure"

type BlogHandler struct {
	// config is the configuration for the blog
	config configure.Config
}
