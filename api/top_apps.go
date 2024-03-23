package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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

func HandlerTopApps(w http.ResponseWriter, r *http.Request) {
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

// Structs returned by the /api/top-apps endpoint
type AppsResp struct {
	Apps []AppData
}

type AppData struct {
	Id      string
	Name    string
	IconUrl string
}

// Structs for the JSON response from the iTunes API
// Define the top-level struct that matches the JSON structure
type AppFeedContainer struct {
	Feed AppFeed `json:"feed"`
}

// Feed contains all the details of the feed, including a slice of App details
type AppFeed struct {
	Title     string    `json:"title"`
	ID        string    `json:"id"`
	Author    AppAuthor `json:"author"`
	Links     []AppLink `json:"links"`
	Copyright string    `json:"copyright"`
	Country   string    `json:"country"`
	Icon      string    `json:"icon"`
	Updated   string    `json:"updated"`
	Results   []AppInfo `json:"results"`
}

// Author represents the author of the feed
type AppAuthor struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Link represents a link associated with the feed
type AppLink struct {
	Self string `json:"self"`
}

// AppInfo contains the details about each app in the results array
type AppInfo struct {
	ArtistName    string   `json:"artistName"`
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	ReleaseDate   string   `json:"releaseDate"`
	Kind          string   `json:"kind"`
	ArtworkUrl100 string   `json:"artworkUrl100"`
	Genres        []string `json:"genres"` // Assuming Genres would be a list of strings if populated
	URL           string   `json:"url"`
}
