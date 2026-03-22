package handlers

import (
	"bananajeanss/go-ship/StartTime"
	"bananajeanss/go-ship/db"
	"fmt"
	"net/http"
	"time"
)

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	// get rsvp count from db
	rsvpCount, err := db.RsvpCount()
	if err != nil {
		http.Error(w, "500 internal server error", 500)
		return
	}

	// get uptime of this process
	uptime := time.Now().Unix() - StartTime.GetStartTime()

	// return as json
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"rsvp_count": ` + fmt.Sprint(rsvpCount) + `, "uptime": ` + fmt.Sprint(uptime) + `}`))
}
