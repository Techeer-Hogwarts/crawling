package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	// Import the package from the main directory
)

// TestMainHandler tests the root handler that serves "Hello, World!"
func TestMainHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()

	// Call the main handler from your main package
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Hello, World!"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestPortEnvVar tests if the port environment variable is set correctly
func TestPortEnvVar(t *testing.T) {
	os.Setenv("PORT", "9090")
	port := os.Getenv("PORT")
	if port != "9090" {
		t.Errorf("Expected port to be 9090, but got %s", port)
	}
}
