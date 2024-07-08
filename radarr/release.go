package radarr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"golift.io/starr"
)

const bpRelease = APIver + "/release"

// Release is the output from the Radarr release endpoint.
type Release struct {
	ID                  int64          `json:"id"`
	GUID                string         `json:"guid,omitempty"`
	Quality             starr.Quality  `json:"quality"`                 // QualityModel
	CustomFormats       []any          `json:"customFormats,omitempty"` // CustomFormatResource
	CustomFormatScore   int64          `json:"customFormatScore"`
	QualityWeight       int64          `json:"qualityWeight"`
	Age                 int64          `json:"age"`
	AgeHours            float64        `json:"ageHours"`
	AgeMinutes          float64        `json:"ageMinutes"`
	Size                int64          `json:"size"`
	IndexerID           int64          `json:"indexerId"`
	Indexer             string         `json:"indexer,omitempty"`
	ReleaseGroup        string         `json:"releaseGroup,omitempty"`
	SubGroup            string         `json:"subGroup,omitempty"`
	ReleaseHash         string         `json:"releaseHash,omitempty"`
	Title               string         `json:"title,omitempty"`
	SceneSource         bool           `json:"sceneSource"`
	MovieTitles         []*string      `json:"movieTitles,omitempty"`
	Languages           []*starr.Value `json:"languages,omitempty"` // Language
	MappedMovieID       *int64         `json:"mappedMovieId,omitempty"`
	Approved            bool           `json:"approved"`
	TemporarilyRejected bool           `json:"temporarilyRejected"`
	Rejected            bool           `json:"rejected"`
	TmdbID              int64          `json:"tmdbId"`
	ImdbID              int64          `json:"imdbId"`
	Rejections          []*string      `json:"rejections,omitempty"`
	PublishDate         time.Time      `json:"publishDate"`
	CommentUrl          string         `json:"commentUrl,omitempty"`
	DownloadUrl         string         `json:"downloadUrl,omitempty"`
	InfoUrl             string         `json:"infoUrl,omitempty"`
	DownloadAllowed     bool           `json:"downloadAllowed"`
	ReleaseWeight       int64          `json:"releaseWeight"`
	Edition             string         `json:"edition,omitempty"`
	MagnetUrl           string         `json:"magnetUrl,omitempty"`
	InfoHash            string         `json:"infoHash,omitempty"`
	Seeders             *int64         `json:"seeders,omitempty"`
	Leechers            *int64         `json:"leechers,omitempty"`
	Protocol            starr.Protocol `json:"protocol"` // DownloadProtocol
	MovieID             *int64         `json:"movieId,omitempty"`
	DownloadClientID    *int64         `json:"downloadClientId,omitempty"`
	DownloadClient      string         `json:"downloadClient,omitempty"`
	ShouldOverride      *bool          `json:"shouldOverride,omitempty"`
}

// SearchRelease searches for and returns a list releases available for download.
func (r *Radarr) SearchRelease(movieID int64) ([]*Release, error) {
	return r.SearchReleaseContext(context.Background(), movieID)
}

// SearchReleaseContext searches for and returns a list releases available for download.
func (r *Radarr) SearchReleaseContext(ctx context.Context, movieID int64) ([]*Release, error) {
	req := starr.Request{URI: bpRelease, Query: make(url.Values)}
	req.Query.Set("movieId", starr.Str(movieID))

	var output []*Release
	if err := r.GetInto(ctx, req, &output); err != nil {
		return nil, fmt.Errorf("api.Get(%s): %w", &req, err)
	}

	return output, nil
}

// Grab is the output from the Grab* methods.
type Grab struct {
	GUID                string         `json:"guid"`
	Quality             *starr.Quality `json:"quality"`
	CustomFormatScore   int64          `json:"customFormatScore"`
	QualityWeight       int64          `json:"qualityWeight"`
	Age                 int64          `json:"age"`
	AgeHours            int            `json:"ageHours"`
	AgeMinutes          int            `json:"ageMinutes"`
	Size                int64          `json:"size"`
	IndexerID           int64          `json:"indexerId"`
	SceneSource         bool           `json:"sceneSource"`
	Languages           []*starr.Value `json:"languages"`
	Approved            bool           `json:"approved"`
	TemporarilyRejected bool           `json:"temporarilyRejected"`
	Rejected            bool           `json:"rejected"`
	TmdbID              int64          `json:"tmdbId"`
	ImdbID              int64          `json:"imdbId"`
	PublishDate         time.Time      `json:"publishDate"`
	DownloadAllowed     bool           `json:"downloadAllowed"`
	ReleaseWeight       int64          `json:"releaseWeight"`
	Protocol            string         `json:"protocol"`
	MovieID             int64          `json:"movieId"`
	ShouldOverride      bool           `json:"shouldOverride"`
}

// GrabRelease attempts to download a release for a movie from a search.
// Pass the release for the item from the SearchRelease output, and the movie ID you want the grab associated with.
// If the movieID is 0 then the MappedMovieID in the release is used, but that is not always set.
func (r *Radarr) GrabRelease(release *Release, movieID int64) (*Grab, error) {
	return r.GrabReleaseContext(context.Background(), release)
}

// GrabReleaseContext attempts to download a release for a movie from a search.
// Pass the release for the item from the SearchRelease output, and the movie ID you want the grab associated with.
// If the movieID is 0 then the MappedMovieID in the release is used, but that is not always set.
func (r *Radarr) GrabReleaseContext(ctx context.Context, release *Release) (*Grab, error) {
	grab := struct { // These are the required fields on the Radarr POST /release endpoint.
		G string `json:"guid"`
		I int64  `json:"indexerId"`
	}{G: release.GUID, I: release.IndexerID}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&grab); err != nil {
		return nil, fmt.Errorf("json.Marshal(%s): %w", bpRelease, err)
	}

	var output Grab

	req := starr.Request{URI: bpRelease, Body: &body}
	if err := r.PostInto(ctx, req, &output); err != nil {
		return nil, fmt.Errorf("api.Post(%s): %w", &req, err)
	}

	return &output, nil
}
