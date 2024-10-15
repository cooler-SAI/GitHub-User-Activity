package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchUserActivity_Success(t *testing.T) {
	mockResponse := []Event{
		{
			Type: "PushEvent",
			Repo: struct {
				Name string `json:"name"`
			}{
				Name: "test-repo",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/tester/events" {
			t.Fatalf("Expected to request '/users/testuser/events', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(mockResponse)
		if err != nil {
			return
		}
	}))
	defer ts.Close()

	events, err := fetchUserActivity(ts.URL, "tester")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(events) != 1 || events[0].Type != "PushEvent" || events[0].Repo.Name != "test-repo" {
		t.Errorf("Unexpected event data: %+v", events)
	}
}

func TestFetchUserActivity_Fail(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := fetchUserActivity(ts.URL, "tester")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}
