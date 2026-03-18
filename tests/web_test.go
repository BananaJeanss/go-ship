package tests

import (
	"bananajeanss/go-ship/handlers"
	"bananajeanss/go-ship/db"
	"net/http"
	"net/http/httptest"
	"testing"
	"os"
)

func TestMain(m *testing.M) {
	// setup cwd
    os.Chdir("..")
	// init test database
	db.Init()
	// run tests
    os.Exit(m.Run())
}

func TestHomeHandler(t *testing.T) {
	// hit the index/home route and check for 200 response
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HomeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Logf("body: %s, cwd: %s", rr.Body.String(), os.Getenv("PWD"))
		
	}
}
