package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleReviews(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/reviews?id=1666653815", nil)
	if err != nil {
		t.Fatal(err)
	}

	// print the req
	fmt.Print(req)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerReviews)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var result FeedContainer

	err = json.Unmarshal(rr.Body.Bytes(), &result)
	fmt.Print(result)

	if err != nil {
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			t.Errorf("JSON syntax error at byte offset %d: %s", syntaxErr.Offset, err)
		} else {
			t.Errorf("The handler returned non-JSON body: %v", err)
		}
	}
}

func TestHandleTopApps(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/topApps", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerTopApps)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var result AppFeedContainer

	err = json.Unmarshal(rr.Body.Bytes(), &result)
	fmt.Print(result)

	if err != nil {
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			t.Errorf("JSON syntax error at byte offset %d: %s", syntaxErr.Offset, err)
		} else {
			t.Errorf("The handler returned non-JSON body: %v", err)
		}
	}
}
