package main

import "fmt"
import "net/http"
import "html/template"
import "os"
import "os/signal"
import "path/filepath"
import "strings"
import "github.com/joho/godotenv"
import "bananajeanss/go-ship/handlers"
import "bananajeanss/go-ship/db"

func notFoundHandler(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles("./templates/404.html")
	// if no template, just return 404
	if err != nil {
		http.Error(w, "404 not found :(", 404)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	tmpl.Execute(w, nil)
}

func dynamicHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	var templatePath string
	if path == "/" {
		handlers.HomeHandler(w, r)
		return
	} else {
		templatePath = fmt.Sprintf("./templates%s/index.html", path)
	}

	// no path traversal
	cleanPath := filepath.Clean(templatePath)
	if !strings.HasPrefix(cleanPath, filepath.Clean("./templates")) {
		notFoundHandler(w)
		return
	}

	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		notFoundHandler(w)
		return
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "500 internal server error", 500)
		return
	}

	tmpl.Execute(w, nil)
}

func init() {
	// env vars
	godotenv.Load()

	// handle ctrl+c gracefully to not get an error
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("\nShutting down...")
		os.Exit(0)
	}()
}

func main() {
	// init db
	if err := db.Init(); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		return
	}

	// serve public static files
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// redirect site.webmanifest and favicon.ico to the public folder
	http.HandleFunc("/site.webmanifest", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/site.webmanifest")
	})
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/favicon.ico")
	})

	// specific handlers
	http.HandleFunc("/auth/callback", handlers.AuthCallbackHandler)

	// catch-all
	http.HandleFunc("/", dynamicHandler)

	// listen and serve
	fmt.Printf("Listening & Serving on http://localhost:3000\n")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
