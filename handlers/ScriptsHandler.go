package handlers

import (
	"bytes"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func ScriptsHandler(w http.ResponseWriter, r *http.Request) {
	// returns available go scripts and contents of public/scripts as json

	// only allow GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Failed to get current working directory", http.StatusInternalServerError)
		return
	}

	pathTo := cwd + "/public/scripts/"

	files, err := os.ReadDir(pathTo)
	if err != nil {
		http.Error(w, "Failed to read scripts directory", http.StatusInternalServerError)
		return
	}

	var scriptNames []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
			scriptNames = append(scriptNames, file.Name())
		}
	}

	var scriptContents []string
	for _, name := range scriptNames {
		content, err := os.ReadFile(pathTo + name)
		if err != nil {
			http.Error(w, "Failed to read script file: "+name, http.StatusInternalServerError)
			return
		}
		scriptContents = append(scriptContents, string(content))
	}

	var buffer bytes.Buffer
	buffer.WriteString("{")
	for i, name := range scriptNames {
		escapedContent := url.QueryEscape(scriptContents[i])
		buffer.WriteString(`"` + name + `":"` + escapedContent + `"`)
		if i < len(scriptNames)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("}")

	w.Header().Set("Content-Type", "application/json")
	w.Write(buffer.Bytes())
}