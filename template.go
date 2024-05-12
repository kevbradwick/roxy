package roxy

import (
	"embed"
	"html/template"
)

//go:embed templates/*
var templatesFS embed.FS
var accessDeniedTemplate *template.Template

func init() {
	accessDeniedTemplate = template.Must(template.ParseFS(templatesFS, "templates/access-denied.html"))
}

type AccessDenied struct {
	EmailName    string
	Email        string
	RequestID    string
	ClientID     string
	ForwardedURL string
}
