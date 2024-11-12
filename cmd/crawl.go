package cmd

import (
	"fmt"

	"github.com/Techeer-Hogwarts/crawling/cmd/blogs"
)

type Image struct {
	Src           string `json:"src"`
	Alt           string `json:"alt"`
	FetchPriority string `json:"fetchpriority"`
	Decoding      string `json:"decoding"`
	DataNimg      string `json:"data-nimg"`
	Style         string `json:"style"`
}

type Link struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

type Tag struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

type SubInfo struct {
	Date     string `json:"date"`
	Comments string `json:"comments"`
	Likes    string `json:"likes"`
}

type DivContent struct {
	HTML    string  `json:"html"`
	Links   []Link  `json:"links"`
	Images  []Image `json:"images"`
	Tags    []Tag   `json:"tags"`
	SubInfo SubInfo `json:"sub_info"`
	Text    string  `json:"text"`
}

type BlogPosts struct {
	Title     string   `json:"title"`
	Text      string   `json:"text"`
	Thumbnail string   `json:"thumbnail"`
	Link      string   `json:"link"`
	Tags      []string `json:"tags"`
	Date      string   `json:"date"`
}

type BlogResponse struct {
	UserID string      `json:"user_id"`
	Posts  []BlogPosts `json:"posts"`
}

type BlogRequest struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func CrawlBlog(targetURL, host string) (BlogResponse, error) {
	switch host {
	case "medium.com":
		return BlogResponse{}, blogs.ProcessMediumBlog(targetURL)
	case "velog.io":
		return BlogResponse{}, blogs.ProcessVelogBlog(targetURL)
	case "tistory.com":
		return BlogResponse{}, blogs.ProcessTistoryBlog(targetURL)
	default:
		return BlogResponse{}, fmt.Errorf("unsupported host: %s", host)
	}

	// return BlogResponse{}, nil
}
