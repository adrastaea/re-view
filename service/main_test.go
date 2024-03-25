package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleReviews(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/reviews?id=595068606", nil)
	if err != nil {
		t.Fatal(err)
	}

	// print the req
	fmt.Print(req)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleReviews)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var result FeedContainer

	err = json.Unmarshal(rr.Body.Bytes(), &result)

	if err != nil {
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			t.Errorf("JSON syntax error at byte offset %d: %s", syntaxErr.Offset, err)
		} else {
			t.Errorf("The handler returned non-JSON body: %v", err)
		}
	}
}

func TestHandleTopApps(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/top-apps", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleTopApps)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var result AppFeedContainer

	err = json.Unmarshal(rr.Body.Bytes(), &result)

	if err != nil {
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			t.Errorf("JSON syntax error at byte offset %d: %s", syntaxErr.Offset, err)
		} else {
			t.Errorf("The handler returned non-JSON body: %v", err)
		}
	}
}
