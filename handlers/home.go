package handlers

import "net/http"
import "html/template"
import "os"
import "bananajeanss/go-ship/db"

type PageData struct {
	HCAAuthURL string
	IsAuthed   bool
	RSVPCount  int
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		http.Error(w, "500 internal server error", 500)
		return
	}

	clientId := os.Getenv("HCA_CLIENT_ID")
	redirectURI := os.Getenv("HCA_REDIRECT_URI")

	loggedIn := false
	cookie, err := r.Cookie("goship_session")
	if err == nil {
		loggedIn = db.IsLoggedIn(cookie.Value)
	}

	rsvpCount, err := db.RsvpCount()
	if err != nil {
		rsvpCount = 9999 // 9999 cause it's a unrealistic expectation and i can use that clientside to show sumn went wrong
	}

	data := PageData{
		HCAAuthURL: "https://auth.hackclub.com/oauth/authorize?client_id=" + clientId + "&redirect_uri=" + redirectURI + "&response_type=code&scope=openid+profile+email+name+profile+slack_id+verification_status&prompt=consent",
		IsAuthed:   loggedIn,
		RSVPCount: rsvpCount,
	}

	tmpl.Execute(w, data)
}
