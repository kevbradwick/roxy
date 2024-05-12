package roxy

import "embed"

//go:embed templates/*
var templates embed.FS

type AccessDenied struct {
	EmailName    string
	Email        string
	RequestID    string
	ClientIP     string
	ForwardedURL string
}
