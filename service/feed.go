package main

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
