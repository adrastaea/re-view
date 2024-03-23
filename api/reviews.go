package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
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

func addReviewsToDB(appId string, reviews []ReviewData) error {
	databaseURL := os.Getenv("POSTGRES_URL")
	if databaseURL == "" {
		return fmt.Errorf("database URL is not set")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	defer db.Close()

	// Ensure the reviews table exists
	createTableSQL := `
CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    review_id TEXT NOT NULL,
    app_id TEXT NOT NULL,
    date TEXT,
    author TEXT,
    score TEXT,
    content TEXT,
    UNIQUE (review_id, app_id)
);
;`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to ensure reviews table exists: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %v", err)
	}

	stmt, err := tx.Prepare("INSERT INTO reviews (review_id, app_id, date, author, score, content) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (review_id, app_id) DO NOTHING;")
	if err != nil {
		tx.Rollback() // Rollback in case of statement preparation errors
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	for _, review := range reviews {
		_, err := stmt.Exec(review.Id, appId, review.Date, review.Author, review.Score, review.Content)
		if err != nil {
			tx.Rollback() // Rollback the transaction in case of execution errors
			return fmt.Errorf("failed to insert review: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %v", err)
	}

	return nil
}

// HandlerReviews serves as the HTTP handler for fetching and returning app reviews.
func HandlerReviews(w http.ResponseWriter, r *http.Request) {

	// add appId as a query parameter
	// e.x. http://localhost:8080/api/reviews?id=123456
	appId := r.URL.Query().Get("id")
	if appId == "" {
		http.Error(w, "Missing id query parameter", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("https://itunes.apple.com/us/rss/customerreviews/id=%s/sortBy=mostRecent/page=1/json", appId)
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

	if err := addReviewsToDB(appId, reviews.Reviews); err != nil {
		log.Print(err)
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
