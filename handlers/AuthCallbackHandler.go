package handlers

import (
	"bananajeanss/go-ship/db"
	"encoding/json"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"net/url"
	"os"
)

type AuthCallbackData struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	IdToken     string `json:"id_token"`
}

func AuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "400 bad request: missing code", 400)
		return
	}

	// exchange code for tokens
	response, err := http.PostForm("https://auth.hackclub.com/oauth/token", url.Values{
		"client_id":     {os.Getenv("HCA_CLIENT_ID")},
		"client_secret": {os.Getenv("HCA_CLIENT_SECRET")},
		"redirect_uri":  {os.Getenv("HCA_REDIRECT_URI")},
		"code":          {code},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		http.Error(w, "500 internal server error: failed to exchange code", 500)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		http.Error(w, "500 internal server error: token exchange failed", 500)
		return
	}

	tokenData := AuthCallbackData{}
	json.NewDecoder(response.Body).Decode(&tokenData)

	// fetch HCA public keys from JWKS endpoint
	jwks, err := keyfunc.NewDefault([]string{"https://auth.hackclub.com/oauth/discovery/keys"})
	if err != nil {
		http.Error(w, "500 internal server error: failed to fetch JWKS", 500)
		return
	}

	// validate the ID token with RS256 public key
	token, err := jwt.Parse(tokenData.IdToken, jwks.Keyfunc)
	if err != nil || !token.Valid {
		http.Error(w, "500 internal server error: failed to validate ID token", 500)
		return
	}

	// after validating the token, call userinfo endpoint
	req, _ := http.NewRequest("GET", "https://auth.hackclub.com/oauth/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+tokenData.AccessToken)

	client := &http.Client{}
	meResponse, err := client.Do(req)
	if err != nil {
		http.Error(w, "500 internal server error: failed to get user info", 500)
		return
	}
	defer meResponse.Body.Close()

	var userInfo map[string]interface{}
	json.NewDecoder(meResponse.Body).Decode(&userInfo)

	// save user info to db, and create a session
	err = db.SaveUser(userInfo)
	if err != nil {
		http.Error(w, "500 internal server error: failed to save user", 500)
		return
	}

	// if save success, create a session, set token.
	sessiontoken, err := db.NewSession(userInfo["sub"].(string))
	if err != nil {
		http.Error(w, "500 internal server error: failed to create session", 500)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "goship_session",
		Value:    sessiontoken,
		HttpOnly: true,
		Path:     "/",                  // allow cookie on all paths
		MaxAge:   60 * 60 * 24 * 30,    // 30 days
		SameSite: http.SameSiteLaxMode, // adios CSRF
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
