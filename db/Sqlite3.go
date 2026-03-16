package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init() error {
	var err error
	DB, err = sql.Open("sqlite", "./goship.db")
	if err != nil {
		return err
	}

	// create tables if they don't exist
	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            sub TEXT UNIQUE NOT NULL,
            name TEXT,
            email TEXT,
            slack_id TEXT,
            ysws_eligible BOOLEAN DEFAULT FALSE,
            tokens INTEGER DEFAULT 0,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );

        CREATE TABLE IF NOT EXISTS sessions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            session_token TEXT UNIQUE NOT NULL,
            expires_at DATETIME NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id)
        );
    `)
	return err
}

func SaveUser(userInfo map[string]interface{}) error {
	sub := userInfo["sub"].(string)
	name := userInfo["name"].(string)
	email := userInfo["email"].(string)
	slackId := userInfo["slack_id"].(string)
	yswsEligible := userInfo["verification_status"].(string) == "verified"

	_, err := DB.Exec(`
		INSERT OR IGNORE INTO users (sub, name, email, slack_id, ysws_eligible)
		VALUES (?, ?, ?, ?, ?)
	`, sub, name, email, slackId, yswsEligible)
	return err
}

// create new session, return cookie value
func NewSession(sub string) (string, error) {
	var userId int
	err := DB.QueryRow("SELECT id FROM users WHERE sub = ?", sub).Scan(&userId)
	if err != nil {
		fmt.Print("Error finding user for session: ", err)
		return "", err
	}

	// generate random session token
	tokenBytes := make([]byte, 32)
	_, err = rand.Read(tokenBytes)
	if err != nil {
		fmt.Print("Error generating session token: ", err)
		return "", err
	}
	sessionToken := hex.EncodeToString(tokenBytes)

	// insert session into db, expires in 30 days
	_, err = DB.Exec(`
		INSERT INTO sessions (user_id, session_token, expires_at)
		VALUES (?, ?, datetime('now', '+30 days'))
	`, userId, sessionToken)
	if err != nil {
		fmt.Print("Error saving session to db: ", err)
		return "", err
	}

	return sessionToken, nil
}
	
func IsLoggedIn(Cookie string) bool {
    var exists bool
    err := DB.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM sessions 
            WHERE session_token = ? 
            AND expires_at > datetime('now')
        )
    `, Cookie).Scan(&exists)
    
    return err == nil && exists
}

var cachedRsvpCount = -1
var lastRsvpCountTime = time.Time{}

func RsvpCount() (int, error) {
	var count int

	// check if cache exists yet
	if cachedRsvpCount != -1 && time.Since(lastRsvpCountTime) < 5*time.Minute {
		return cachedRsvpCount, nil
	}

	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err == nil {
		cachedRsvpCount = count
		lastRsvpCountTime = time.Now()
	}
	return count, err
}