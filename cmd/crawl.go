package cmd

import (
	"fmt"

	"github.com/Techeer-Hogwarts/crawling/cmd/blogs"
)

func CrawlBlog(targetURL, host string) (blogs.BlogResponse, error) {
	switch host {
	case "medium.com":
		return blogs.ProcessMediumBlog(targetURL)
	case "velog.io":
		return blogs.ProcessVelogBlog(targetURL)
	case "tistory.com":
		return blogs.BlogResponse{}, blogs.ProcessTistoryBlog(targetURL)
	default:
		return blogs.BlogResponse{}, fmt.Errorf("unsupported host: %s", host)
	}
}
