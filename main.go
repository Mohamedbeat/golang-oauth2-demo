package main

import (
	"net/http"
	"tst/auth"
)

func main() {
	oauthConfig := auth.NewOAuthConfig("http://localhost:8080")

	http.HandleFunc("/auth/github", oauthConfig.HandleGitHubLogin)
	http.HandleFunc("/auth/github/callback", oauthConfig.HandleGitHubCallback)

	http.ListenAndServe(":8080", nil)
}
