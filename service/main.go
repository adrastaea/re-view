package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
func parseReviewsResp(feed FeedContainer, saveFile string, timePeriod time.Duration) (ReviewsResp, error) {
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
		if _, err := file.Write([]byte("App,Id,Date,Author,Score,Content\n")); err != nil {
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
				log.Printf("Failed to convert entry to review: %v", err)
				continue
			}
			// Append the review to the return list and add to persistent storage
			reviews.Reviews = append(reviews.Reviews, review)
			if _, err := file.Write([]byte(feed.Feed.Title.Label + "," + review.Id + "," + review.Date + "," + review.Author + "," + review.Score + "," + strings.ReplaceAll(review.Content, "\n", "\\n") + "\n")); err != nil {
				log.Printf("Failed to write entry to file: %v", err)
				continue
			}
		}
	}
	return reviews, nil
}

func convertResultsToAppData(app AppInfo) AppData {
	return AppData{
		Id:      app.ID,
		Name:    app.Name,
		IconUrl: app.ArtworkUrl100,
	}
}

func convertAppFeedContainertoAppsResp(feed AppFeedContainer) AppsResp {
	var apps AppsResp
	for _, app := range feed.Feed.Results {
		apps.Apps = append(apps.Apps, convertResultsToAppData(app))
	}
	return apps
}

func handleReviews(w http.ResponseWriter, r *http.Request) {

	// add app_id as a query parameter
	// e.x. http://localhost:8080/api/reviews?id=123456
	app_id := r.URL.Query().Get("id")
	if app_id == "" {
		http.Error(w, "Missing id query parameter", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("https://itunes.apple.com/us/rss/customerreviews/id=%s/sortBy=mostRecent/page=1/json", app_id)
	log.Printf("Fetching reviews from %s", url)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// decode the response body into a Feed struct
	var feed FeedContainer
	if err := json.NewDecoder(resp.Body).Decode(&feed); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reviews, err := parseReviewsResp(feed, "reviews.txt", 2*24*time.Hour)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// encode the Feed struct into a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(reviews); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleTopApps(w http.ResponseWriter, r *http.Request) {
	TOP_APP_URL := "https://rss.applemarketingtools.com/api/v2/us/apps/top-free/10/apps.json"
	resp, err := http.Get(TOP_APP_URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	log.Printf("Fetching top apps from %s", TOP_APP_URL)

	// decode the response body into a AppStoreFeed struct
	var feed AppFeedContainer
	if err := json.NewDecoder(resp.Body).Decode(&feed); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apps := convertAppFeedContainertoAppsResp(feed)
	// encode the AppsResp struct into a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(apps); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	// Open a file for logging
	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	// defer logFile.Close()
	log.SetOutput(logFile)

	// Set up routers
	http.HandleFunc("/api/reviews", handleReviews)
	http.HandleFunc("/api/top-apps", handleTopApps)

	// CORS configuration for development purposes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	})

	println("Server is running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Printf("Failed to start server: %s", err)
		panic(err)
	}
	log.Println("Server is running at http://localhost:8080")
}
