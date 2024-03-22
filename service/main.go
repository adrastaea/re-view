package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type ReviewContainer struct {
	Reviews []Review
}

type Review struct {
	Id      string
	Date    string
	Author  string
	Score   string
	Content string
}

func convertEntryToReview(entry Entry) Review {
	return Review{
		Id:      entry.Id.Label,
		Date:    entry.Updated.Label,
		Author:  entry.Author.Name.Label,
		Score:   entry.ImRating.Label,
		Content: entry.Content.Label,
	}
}

func handleReviews(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get("https://itunes.apple.com/us/rss/customerreviews/id=595068606/sortBy=mostRecent/page=1/json")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// decode the response body into a Feed struct
	var feed FeedContainer
	if err := json.NewDecoder(resp.Body).Decode(&feed); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Feed: %v", feed)

	// Open a file for writing
	file, err := os.OpenFile("reviews.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if file is empty
	fi, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if fi.Size() == 0 {
		// Write the header
		if _, err := file.Write([]byte("Id,Date,Author,Score,Content\n")); err != nil {
			file.Close()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	now := time.Now()
	var reviews ReviewContainer

	// Put all entries in the feed into a ReviewContainer
	for _, entry := range feed.Feed.Entry {
		// Check if Update.Label date is within the last 14 days
		timestamp, err := time.Parse(time.RFC3339, entry.Updated.Label)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if now.Sub(timestamp) <= 14*24*time.Hour {
			review := convertEntryToReview(entry)
			reviews.Reviews = append(reviews.Reviews, review)
			if _, err := file.Write([]byte(review.Id + "," + review.Date + "," + review.Author + "," + review.Score + "," + review.Content + "\n")); err != nil {
				log.Printf("Failed to write entry to file: %v", err)
				file.Close() // ignore error; Write error takes precedence
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	// Close the file
	if err := file.Close(); err != nil {
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

func main() {
	// Open a file for logging
	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Set up routers
	http.HandleFunc("/api/reviews", handleReviews)

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
