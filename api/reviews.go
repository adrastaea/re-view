package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func convertEntryToReview(entry Entry) (ReviewData, error) {
	if entry.Id.Label == "" || entry.Updated.Label == "" || entry.Author.Name.Label == "" || entry.ImRating.Label == "" || entry.Content.Label == "" {
		return ReviewData{}, errors.New("missing required fields in entry")
	}

	return ReviewData{
		Id:      entry.Id.Label,
		Date:    entry.Updated.Label,
		Author:  entry.Author.Name.Label,
		Score:   entry.ImRating.Label,
		Content: entry.Content.Label,
	}, nil
}

// parseReviewsResp parses the response from the iTunes API and returns a ReviewsResp struct.
// It also saves the reviews to a file if within the specified time period.
func parseReviewsResp(feed FeedContainer, timePeriod time.Duration) (ReviewsResp, error) {
	var reviews ReviewsResp
	now := time.Now()
	for _, entry := range feed.Feed.Entry {
		reviewTime, err := time.Parse(time.RFC3339, entry.Updated.Label)
		if err != nil {
			continue // Skip entries with invalid dates
		}
		if now.Sub(reviewTime) <= timePeriod {
			review, err := convertEntryToReview(entry)
			if err != nil {
				continue // Skip entries with missing required fields
			}

			// Append the review to the response and write it to the file
			reviews.Reviews = append(reviews.Reviews, review)
		}
	}
	return reviews, nil
}

// HandlerReviews serves as the HTTP handler for fetching and returning app reviews.
func HandlerReviews(w http.ResponseWriter, r *http.Request) {

	// add app_id as a query parameter
	// e.x. http://localhost:8080/api/reviews?id=123456
	app_id := r.URL.Query().Get("id")
	if app_id == "" {
		http.Error(w, "Missing id query parameter", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("https://itunes.apple.com/us/rss/customerreviews/id=%s/sortBy=mostRecent/page=1/json", app_id)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "failed fetching reviews", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// decode the response body into a Feed struct
	var feed FeedContainer
	if err := json.NewDecoder(resp.Body).Decode(&feed); err != nil {
		http.Error(w, "failed decoding response body", http.StatusInternalServerError)
		return
	}

	reviews, err := parseReviewsResp(feed, 2*24*time.Hour)
	if err != nil {
		http.Error(w, "failed parsing resp from feed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reviews); err != nil {
		http.Error(w, "failed encoding json", http.StatusInternalServerError)
		return
	}
}

// Structs returned by the /api/reviews endpoint
type ReviewsResp struct {
	Reviews []ReviewData
}

type ReviewData struct {
	Id      string
	Date    string
	Author  string
	Score   string
	Content string
}

// Structs for the JSON response from the iTunes API
type FeedContainer struct {
	Feed Feed `json:"feed"`
}
type Feed struct {
	Author  Author  `json:"author"`
	Entry   []Entry `json:"entry"`
	Updated Label   `json:"updated"`
	Rights  Label   `json:"rights"`
	Title   Label   `json:"title"`
	Icon    Label   `json:"icon"`
	Link    []Link  `json:"link"`
	Id      Label   `json:"id"`
}

type Author struct {
	Name Label `json:"name"`
	Uri  Label `json:"uri"`
}

type Entry struct {
	Author        Author      `json:"author"`
	Updated       Label       `json:"updated"`
	ImRating      Label       `json:"im:rating"`
	ImVersion     Label       `json:"im:version"`
	Id            Label       `json:"id"`
	Title         Label       `json:"title"`
	Content       Content     `json:"content"`
	Link          EntryLink   `json:"link"`
	ImVoteSum     Label       `json:"im:voteSum"`
	ImContentType ContentType `json:"im:contentType"`
	ImVoteCount   Label       `json:"im:voteCount"`
}

type Label struct {
	Label string `json:"label"`
}

type Content struct {
	Label      string     `json:"label"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	Type  string `json:"type,omitempty"`
	Rel   string `json:"rel,omitempty"`
	Href  string `json:"href,omitempty"`
	Term  string `json:"term,omitempty"`
	Label string `json:"label,omitempty"`
}

type Link struct {
	Attributes Attributes `json:"attributes"`
}

type EntryLink struct {
	Attributes Attributes `json:"attributes"`
}

type ContentType struct {
	Attributes Attributes `json:"attributes"`
}
