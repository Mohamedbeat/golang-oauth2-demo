package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type OAuthConfig struct {
	GitHub *oauth2.Config
}

type UserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Provider  string `json:"provider"`
}

func NewOAuthConfig(redirectURL string) *OAuthConfig {
	return &OAuthConfig{
		GitHub: &oauth2.Config{
			ClientID:     "client_id_1234567890abcdef",
			ClientSecret: "client_secret_abcdef1234567890",
			RedirectURL:  redirectURL + "/auth/github/callback",
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		},
	}
}

// Generate random state for OAuth flow
func generateStateOauthCookie(w http.ResponseWriter) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	http.SetCookie(w, &http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true, // Set to false for local development
	})

	return state
}

// GitHub OAuth handlers
func (oc *OAuthConfig) HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	state := generateStateOauthCookie(w)
	url := oc.GitHub.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (oc *OAuthConfig) HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	// Validate state parameter
	state := r.URL.Query().Get("state")
	cookie, err := r.Cookie("oauthstate")
	if err != nil || cookie.Value != state {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := oc.GitHub.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := oc.getGitHubUserInfo(token)
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("user info")
	fmt.Println(userInfo)

	oc.handleSuccessfulAuth(w, r, userInfo)
}
func (oc *OAuthConfig) getGitHubUserInfo(token *oauth2.Token) (*UserInfo, error) {
	client := oc.GitHub.Client(context.Background(), token)

	// Get basic user info
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user struct {
		ID        int    `json:"id"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
		Login     string `json:"login"`
	}

	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}

	// If email is not public, fetch it from the emails endpoint
	if user.Email == "" {
		email, err := oc.getGitHubPrimaryEmail(token)
		if err != nil {
			return nil, err
		}
		user.Email = email
	}

	return &UserInfo{
		ID:        fmt.Sprintf("%d", user.ID),
		Email:     user.Email,
		Name:      user.Name,
		AvatarURL: user.AvatarURL,
		Provider:  "github",
	}, nil
}

func (oc *OAuthConfig) getGitHubPrimaryEmail(token *oauth2.Token) (string, error) {
	client := oc.GitHub.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.Unmarshal(data, &emails); err != nil {
		return "", err
	}

	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	return "", fmt.Errorf("no primary email found")
}

func (oc *OAuthConfig) handleSuccessfulAuth(w http.ResponseWriter, r *http.Request, userInfo *UserInfo) {
	// Generate JWT token
	tokenString, err := generateJWTToken(userInfo)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Redirect to frontend with token as URL parameter
	frontendURL := "http://localhost:5500?token=" + tokenString
	http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)
	// // Generate JWT token
	// tokenString, err := generateJWTToken(userInfo)
	// if err != nil {
	// 	http.Error(w, "Failed to generate token", http.StatusInternalServerError)
	// 	return
	// }

	// // Return token to client
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"token":   tokenString,
	// 	"user":    userInfo,
	// 	"message": "GitHub login successful",
	// })
}

func generateJWTToken(userInfo *UserInfo) (string, error) {
	// Implement JWT token generation
	// You can use packages like github.com/golang-jwt/jwt
	return "your-jwt-token", nil
}
