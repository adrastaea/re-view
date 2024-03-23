package main

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
