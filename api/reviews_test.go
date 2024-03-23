package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerReviews(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/reviews?id=1666653815", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerReviews)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("HandlerReviews returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Assuming ReviewsResp is the expected type
	var result ReviewsResp
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			t.Errorf("JSON syntax error at byte offset %d: %s", syntaxErr.Offset, err)
		} else {
			t.Errorf("HandlerReviews returned a non-JSON body: %v", err)
		}
	}

	fmt.Printf("Reviews: %v\n", result.Reviews)

}

func TestHandlerTopApps(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/top-apps", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerTopApps)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("HandlerTopApps returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var result AppsResp
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			t.Errorf("JSON syntax error at byte offset %d: %s", syntaxErr.Offset, err)
		} else {
			t.Errorf("HandlerTopApps returned a non-JSON body: %v", err)
		}
	}

	fmt.Printf("Top Apps: %v\n", result.Apps)

	// assert that the response body is a JSON array
	if len(result.Apps) == 0 {
		t.Errorf("HandlerTopApps returned an empty list of apps")

	}
}
