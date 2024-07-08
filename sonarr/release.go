package sonarr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"golift.io/starr"
)

const bpRelease = APIver + "/release"

// Release is the output from the Sonarr release endpoint.
type Release struct {
	ID                           int64                 `json:"id"`
	GUID                         string                `json:"guid,omitempty"`
	Quality                      starr.Quality         `json:"quality"` // QualityModel
	QualityWeight                int64                 `json:"qualityWeight"`
	Age                          int64                 `json:"age"`
	AgeHours                     float64               `json:"ageHours"`
	AgeMinutes                   float64               `json:"ageMinutes"`
	Size                         int64                 `json:"size"`
	IndexerID                    int64                 `json:"indexerId"`
	Indexer                      string                `json:"indexer,omitempty"`
	ReleaseGroup                 string                `json:"releaseGroup,omitempty"`
	SubGroup                     string                `json:"subGroup,omitempty"`
	ReleaseHash                  string                `json:"releaseHash,omitempty"`
	Title                        string                `json:"title,omitempty"`
	FullSeason                   bool                  `json:"fullSeason"`
	SceneSource                  bool                  `json:"sceneSource"`
	SeasonNumber                 int64                 `json:"seasonNumber"`
	Languages                    []*starr.Value        `json:"languages,omitempty"` // Language
	LanguageWeight               int64                 `json:"languageWeight"`
	AirDate                      string                `json:"airDate,omitempty"`
	SeriesTitle                  string                `json:"seriesTitle,omitempty"`
	EpisodeNumbers               []*int64              `json:"episodeNumbers,omitempty"`
	AbsoluteEpisodeNumbers       []*int64              `json:"absoluteEpisodeNumbers,omitempty"`
	MappedSeasonNumber           *int64                `json:"mappedSeasonNumber,omitempty"`
	MappedEpisodeNumbers         []*int64              `json:"mappedEpisodeNumbers,omitempty"`
	MappedAbsoluteEpisodeNumbers []*int64              `json:"mappedAbsoluteEpisodeNumbers,omitempty"`
	MappedSeriesID               *int64                `json:"mappedSeriesId,omitempty"`
	MappedEpisodeInfo            []*ReleaseEpisodeInfo `json:"mappedEpisodeInfo,omitempty"` // ReleaseEpisodeResource
	Approved                     bool                  `json:"approved"`
	TemporarilyRejected          bool                  `json:"temporarilyRejected"`
	Rejected                     bool                  `json:"rejected"`
	TvdbID                       int64                 `json:"tvdbId"`
	TvRageID                     int64                 `json:"tvRageId"`
	Rejections                   []*string             `json:"rejections,omitempty"`
	PublishDate                  time.Time             `json:"publishDate"`
	CommentUrl                   string                `json:"commentUrl,omitempty"`
	DownloadUrl                  string                `json:"downloadUrl,omitempty"`
	InfoUrl                      string                `json:"infoUrl,omitempty"`
	EpisodeRequested             bool                  `json:"episodeRequested"`
	DownloadAllowed              bool                  `json:"downloadAllowed"`
	ReleaseWeight                int64                 `json:"releaseWeight"`
	CustomFormats                []*CustomFormatOutput `json:"customFormats,omitempty"` // CustomFormatResource
	CustomFormatScore            int64                 `json:"customFormatScore"`
	SceneMapping                 ReleaseSceneMapping   `json:"sceneMapping"` // AlternateTitleResource
	MagnetUrl                    string                `json:"magnetUrl,omitempty"`
	InfoHash                     string                `json:"infoHash,omitempty"`
	Seeders                      *int64                `json:"seeders,omitempty"`
	Leechers                     *int64                `json:"leechers,omitempty"`
	Protocol                     starr.Protocol        `json:"protocol"` // DownloadProtocol
	// IndexerFlags                 int64                 `json:"indexerFlags"`
	IsDaily                  bool     `json:"isDaily"`
	IsAbsoluteNumbering      bool     `json:"isAbsoluteNumbering"`
	IsPossibleSpecialEpisode bool     `json:"isPossibleSpecialEpisode"`
	Special                  bool     `json:"special"`
	SeriesID                 *int64   `json:"seriesId,omitempty"`
	EpisodeID                *int64   `json:"episodeId,omitempty"`
	EpisodeIds               []*int64 `json:"episodeIds,omitempty"`
	DownloadClientID         *int64   `json:"downloadClientId,omitempty"`
	DownloadClient           string   `json:"downloadClient,omitempty"`
	ShouldOverride           *bool    `json:"shouldOverride,omitempty"`
}

// ReleaseSceneMapping is part of a release.
type ReleaseSceneMapping struct {
	Title             string `json:"title"`
	SeasonNumber      int    `json:"seasonNumber"`
	SceneSeasonNumber int    `json:"sceneSeasonNumber"`
	SceneOrigin       string `json:"sceneOrigin"`
	Comment           string `json:"comment"`
}

// ReleaseEpisodeInfo is part of a release.
type ReleaseEpisodeInfo struct {
	ID                    int64  `json:"id"`
	SeasonNumber          int    `json:"seasonNumber"`
	EpisodeNumber         int    `json:"episodeNumber"`
	AbsoluteEpisodeNumber int    `json:"absoluteEpisodeNumber"`
	Title                 string `json:"title"`
}

// SearchRelease is the input needed to search for releases through Sonarr.
type SearchRelease struct {
	SeriesID     int64 `json:"seriesId"`
	EpisodeID    int64 `json:"episodeId"`
	SeasonNumber int   `json:"seasonNumber"`
}

// SearchRelease searches for and returns a list releases available for download.
func (s *Sonarr) SearchRelease(input *SearchRelease) ([]*Release, error) {
	return s.SearchReleaseContext(context.Background(), input)
}

// SearchReleaseContext searches for and returns a list releases available for download.
func (s *Sonarr) SearchReleaseContext(ctx context.Context, input *SearchRelease) ([]*Release, error) {
	req := starr.Request{URI: bpRelease, Query: make(url.Values)}
	req.Query.Set("seriesId", starr.Str(input.SeriesID))
	req.Query.Set("episodeId", starr.Str(input.EpisodeID))
	req.Query.Set("seasonNumber", strconv.Itoa(input.SeasonNumber))

	var output []*Release
	if err := s.GetInto(ctx, req, &output); err != nil {
		return nil, fmt.Errorf("api.Get(%s): %w", &req, err)
	}

	return output, nil
}

// Grab is the output from the Grab* methods.
type Grab struct {
	Approved                 bool           `json:"approved"`
	DownloadAllowed          bool           `json:"downloadAllowed"`
	EpisodeRequested         bool           `json:"episodeRequested"`
	FullSeason               bool           `json:"fullSeason"`
	Special                  bool           `json:"special"`
	TemporarilyRejected      bool           `json:"temporarilyRejected"`
	IsAbsoluteNumbering      bool           `json:"isAbsoluteNumbering"`
	IsDaily                  bool           `json:"isDaily"`
	IsPossibleSpecialEpisode bool           `json:"isPossibleSpecialEpisode"`
	Rejected                 bool           `json:"rejected"`
	SceneSource              bool           `json:"sceneSource"`
	AgeHours                 int            `json:"ageHours"`
	AgeMinutes               int            `json:"ageMinutes"`
	SeasonNumber             int            `json:"seasonNumber"`
	Size                     int            `json:"size"`
	Age                      int64          `json:"age"`
	CustomFormatScore        int64          `json:"customFormatScore"`
	IndexerFlags             int64          `json:"indexerFlags"`
	IndexerID                int64          `json:"indexerId"`
	LanguageWeight           int64          `json:"languageWeight"`
	QualityWeight            int64          `json:"qualityWeight"`
	ReleaseWeight            int64          `json:"releaseWeight"`
	TvRageID                 int64          `json:"tvRageId"`
	TvdbID                   int64          `json:"tvdbId"`
	PublishDate              time.Time      `json:"publishDate"`
	GUID                     string         `json:"guid"`
	Protocol                 starr.Protocol `json:"protocol"`
}

// Grab adds a release and attempts to download it. Use this with Pr*wlarr search output.
func (s *Sonarr) Grab(guid string, indexerID int64) (*Grab, error) {
	return s.GrabContext(context.Background(), guid, indexerID)
}

// GrabContext adds a release and attempts to download it. Use this with Pr*wlarr search output.
func (s *Sonarr) GrabContext(ctx context.Context, guid string, indexerID int64) (*Grab, error) {
	return s.GrabReleaseContext(ctx, &Release{IndexerID: indexerID, GUID: guid})
}

// GrabRelease adds a release and attempts to download it.
// Pass the release for the item from the SearchRelease output.
func (s *Sonarr) GrabRelease(release *Release) (*Grab, error) {
	return s.GrabReleaseContext(context.Background(), release)
}

// GrabReleaseContext adds a release and attempts to download it.
// Pass the release for the item from the SearchRelease output.
func (s *Sonarr) GrabReleaseContext(ctx context.Context, release *Release) (*Grab, error) {
	grab := struct { // We only use/need the guid and indexerID from the release.
		G string `json:"guid"`
		I int64  `json:"indexerId"`
	}{G: release.GUID, I: release.IndexerID}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&grab); err != nil {
		return nil, fmt.Errorf("json.Marshal(%s): %w", bpRelease, err)
	}

	var output Grab

	req := starr.Request{URI: bpRelease, Body: &body}
	if err := s.PostInto(ctx, req, &output); err != nil {
		return nil, fmt.Errorf("api.Post(%s): %w", &req, err)
	}

	return &output, nil
}
