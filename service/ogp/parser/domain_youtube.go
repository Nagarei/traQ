package parser

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/dyatlov/go-opengraph/opengraph"
	"github.com/dyatlov/go-opengraph/opengraph/types/image"
	"github.com/dyatlov/go-opengraph/opengraph/types/video"
	jsonIter "github.com/json-iterator/go"
)

// https://oembed.com/
type OEmbedResponse struct {
	Type            string `json:"type"`
	Title           string `json:"title"`
	AuthorName      string `json:"author_name"`
	AuthorURL       string `json:"author_url"`
	ProviderName    string `json:"provider_name"`
	ProviderURL     string `json:"provider_url"`
	CacheAge        string `json:"cache_age"`
	ThumbnailURL    string `json:"thumbnail_url"`
	ThumbnailWidth  uint64 `json:"thumbnail_width"`
	ThumbnailHeight uint64 `json:"thumbnail_height"`
	URL             string `json:"url"`  // photo
	HTML            string `json:"html"` // video
	Width           uint64 `json:"width"`
	Height          uint64 `json:"height"`
}

func FetchYoutubeInfo(u *url.URL) (*opengraph.OpenGraph, *DefaultPageMeta, error) {
	requestURL := fmt.Sprintf("https://www.youtube.com/oembed?url=%s", url.PathEscape(u.String()))
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, ErrNetwork
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 500 {
		return nil, nil, ErrServer
	} else if resp.StatusCode >= 400 {
		return nil, nil, ErrClient
	}

	var data OEmbedResponse
	if err = jsonIter.ConfigFastest.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, nil, err
	}

	rex := regexp.MustCompile("src=\"[^\"]+\"")
	videoUrl := rex.FindString(data.HTML)

	og := opengraph.OpenGraph{
		Type:        "video.other",
		URL:         u.String(),
		Title:       data.Title,
		Description: data.AuthorName,
		Images: []*image.Image{{
			URL:    data.ThumbnailURL,
			Width:  data.ThumbnailWidth,
			Height: data.ThumbnailHeight,
		}},
		Videos: []*video.Video{{
			URL:    videoUrl,
			Type:   "text/html",
			Width:  data.Width,
			Height: data.Height,
		}},
	}
	meta := DefaultPageMeta{}
	return &og, &meta, nil
}
