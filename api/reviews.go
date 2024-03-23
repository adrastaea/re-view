package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
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

// parseReviewsResp parses the response within a given timePeriod from the iTunes API
// and returns a ReviewsResp struct containing the reviews from the feed.
// It also saves the reviews to a file.
func parseReviewsResp(feed FeedContainer, appId string, saveFile string, timePeriod time.Duration) (ReviewsResp, error) {
	// Open a file for writing
	file, err := os.OpenFile(saveFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return ReviewsResp{}, err
	}
	defer file.Close()

	// Check if file is empty
	fi, err := file.Stat()
	if err != nil {
		return ReviewsResp{}, err
	}
	if fi.Size() == 0 {
		// Write the header
		if _, err := file.Write([]byte("AppId,ReviewId,Date,Author,Score,Content\n")); err != nil {
			return ReviewsResp{}, err
		}
	}

	now := time.Now()
	var reviews ReviewsResp

	// Put all entries in the feed into a ReviewContainer
	for _, entry := range feed.Feed.Entry {
		timestamp, err := time.Parse(time.RFC3339, entry.Updated.Label)
		if err != nil {
			return ReviewsResp{}, err
		}
		if now.Sub(timestamp) <= timePeriod {
			review, err := convertEntryToReview(entry)
			if err != nil {
				continue
			}
			// Append the review to the return list and add to persistent storage
			reviews.Reviews = append(reviews.Reviews, review)
			if _, err := file.Write([]byte(appId + "," + review.Id + "," + review.Date + "," + review.Author + "," + review.Score + "," + strings.ReplaceAll(review.Content, "\n", "\\n") + "\n")); err != nil {
				continue
			}
		}
	}
	return reviews, nil
}

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

	reviews, err := parseReviewsResp(feed, app_id, "reviews.txt", 2*24*time.Hour)
	if err != nil {
		http.Error(w, "failed parsing resp from feed", http.StatusInternalServerError)
		return
	}

	// encode the Feed struct into a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusOK)
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
