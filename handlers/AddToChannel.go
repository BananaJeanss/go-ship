package handlers

import (
	"bananajeanss/go-ship/db"
	"bananajeanss/go-ship/slack"
	"net/http"
	"strings"
)

// adds the user to the slack channel
func PostAddToChannelHandler(w http.ResponseWriter, r *http.Request) {
	// r they logged in
	cookie, err := r.Cookie("goship_session")
	if err != nil || !db.IsLoggedIn(cookie.Value) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// get user info from db
	userInfo, err := db.GetUserInfoBySessionToken(cookie.Value)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	// add them to the channel
	err = slack.AddPersonToChannel(userInfo["slack_id"].(string))
	if err != nil {
		if strings.Contains(err.Error(), "already_in_channel") {
			http.Error(w, "User is already in the channel", 418) // teapot hehe
			return
		}
		http.Error(w, "Failed to add user to channel: "+err.Error(), http.StatusInternalServerError)
		return
	}
}