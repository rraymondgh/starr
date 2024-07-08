package prowlarr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"golift.io/starr"
)

const bpSearch = APIver + "/search"

// Search is the output from the Prowlarr search endpoint.
type Search struct {
	ID          int64     `json:"id"`
	GUID        string    `json:"guid,omitempty"`
	Age         int64     `json:"age"`
	AgeHours    float64   `json:"ageHours"`
	AgeMinutes  float64   `json:"ageMinutes"`
	Size        int64     `json:"size"`
	Files       *int64    `json:"files,omitempty"`
	Grabs       *int64    `json:"grabs,omitempty"`
	IndexerID   int64     `json:"indexerId"`
	Indexer     string    `json:"indexer,omitempty"`
	SubGroup    string    `json:"subGroup,omitempty"`
	ReleaseHash string    `json:"releaseHash,omitempty"`
	Title       string    `json:"title,omitempty"`
	SortTitle   string    `json:"sortTitle,omitempty"`
	ImdbID      int64     `json:"imdbId"`
	TmdbID      int64     `json:"tmdbId"`
	TvdbID      int64     `json:"tvdbId"`
	TvMazeID    int64     `json:"tvMazeId"`
	PublishDate time.Time `json:"publishDate"`
	CommentUrl  string    `json:"commentUrl,omitempty"`
	DownloadUrl string    `json:"downloadUrl,omitempty"`
	InfoUrl     string    `json:"infoUrl,omitempty"`
	PosterUrl   string    `json:"posterUrl,omitempty"`
	// IndexerFlags     []*string      `json:"indexerFlags,omitempty"`
	Categories       []*Category    `json:"categories,omitempty"` // IndexerCategory
	MagnetUrl        string         `json:"magnetUrl,omitempty"`
	InfoHash         string         `json:"infoHash,omitempty"`
	Seeders          *int64         `json:"seeders,omitempty"`
	Leechers         *int64         `json:"leechers,omitempty"`
	Protocol         starr.Protocol `json:"protocol"` // DownloadProtocol
	FileName         string         `json:"fileName,omitempty"`
	DownloadClientID *int64         `json:"downloadClientId,omitempty"`
}

// Category is part of the Search output.
type Category struct {
	ID            int64       `json:"id"`
	Name          string      `json:"name"`
	SubCategories []*Category `json:"subCategories"`
}

// SearchInput is the input to the search endpoint.
type SearchInput struct {
	Query      string  `json:"query"` // Query is required. Fill it in.
	Type       string  `json:"type"`  // defaults to "search" if left empty
	IndexerIDs []int64 `json:"indexerIds"`
	Categories []int64 `json:"categories"`
	Limit      int     `json:"limit"`  // Defaults to 100 if left empty or less than 1.
	Offset     int     `json:"offset"` // Skip this many records.
}

// Search the Prowlarr indexers for media and content. Must provide a Query in the SearchInput.
func (p *Prowlarr) Search(search SearchInput) ([]*Search, error) {
	return p.SearchContext(context.Background(), search)
}

// SearchContext searches the Prowlarr indexers for media and content.
func (p *Prowlarr) SearchContext(ctx context.Context, search SearchInput) ([]*Search, error) {
	const defaultSearchLimit = 100

	if search.Type == "" {
		search.Type = "search"
	}

	if search.Limit < 1 {
		search.Limit = defaultSearchLimit
	}

	if search.Limit < 0 {
		search.Limit = 0
	}

	req := starr.Request{URI: bpSearch, Query: make(url.Values)}
	req.Query.Set("query", search.Query)
	req.Query.Set("type", search.Type)
	req.Query.Set("limit", starr.Str(int64(search.Limit)))
	req.Query.Set("offset", starr.Str(int64(search.Offset)))

	for _, val := range search.Categories {
		req.Query.Add("categories", starr.Str(val))
	}

	for _, val := range search.IndexerIDs {
		req.Query.Add("indexerIds", starr.Str(val))
	}

	var output []*Search
	if err := p.GetInto(ctx, req, &output); err != nil {
		return nil, fmt.Errorf("api.Get(%s): %w", &req, err)
	}

	return output, nil
}

// Grab is the output from the /search POST endpoint.
// GUID may be the only member with a value set.
type Grab struct {
	GUID        string         `json:"guid"`
	Age         int64          `json:"age"`
	AgeHours    int            `json:"ageHours"`
	AgeMinutes  int            `json:"ageMinutes"`
	Size        int            `json:"size"`
	IndexerID   int64          `json:"indexerId"`
	ImdbID      int64          `json:"imdbId"`
	TmdbID      int64          `json:"tmdbId"`
	TvdbID      int64          `json:"tvdbId"`
	TvMazeID    int64          `json:"tvMazeId"`
	PublishDate time.Time      `json:"publishDate"`
	Protocol    starr.Protocol `json:"protocol"`
	FileName    string         `json:"fileName"`
}

// Grab attempts to download a searched item by guid. Use this with Pr*wlarr search output.
func (p *Prowlarr) Grab(guid string, indexerID int64) (*Grab, error) {
	return p.GrabContext(context.Background(), guid, indexerID)
}

// GrabContext attempts to download a searched item by guid. Use this with Pr*wlarr search output.
func (p *Prowlarr) GrabContext(ctx context.Context, guid string, indexerID int64) (*Grab, error) {
	return p.GrabSearchContext(ctx, &Search{IndexerID: indexerID, GUID: guid})
}

// GrabSearch attempts to download an item returned from a search.
// Pass the search for the item from the Search() output.
func (p *Prowlarr) GrabSearch(search *Search) (*Grab, error) {
	return p.GrabSearchContext(context.Background(), search)
}

// GrabSearchContext attempts to download an item returned from a search.
// Pass the search for the item from the Search() output.
func (p *Prowlarr) GrabSearchContext(ctx context.Context, search *Search) (*Grab, error) {
	grab := struct { // We only use/need the guid and indexerID from the search.
		G string `json:"guid"`
		I int64  `json:"indexerId"`
	}{G: search.GUID, I: search.IndexerID}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&grab); err != nil {
		return nil, fmt.Errorf("json.Marshal(%s): %w", bpSearch, err)
	}

	var output Grab

	req := starr.Request{URI: bpSearch, Body: &body}
	if err := p.PostInto(ctx, req, &output); err != nil {
		return nil, fmt.Errorf("api.Post(%s): %w", &req, err)
	}

	return &output, nil
}
