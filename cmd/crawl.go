package cmd

import (
	"fmt"

	"github.com/Techeer-Hogwarts/crawling/cmd/blogs"
)

func CrawlBlog(targetURL, host, crawlingType string) (blogs.BlogResponse, error) {
	if crawlingType == "signUp_blog_fetch" || crawlingType == "blogs_daily_update" {
		switch host {
		case "medium.com":
			return blogs.ProcessMediumBlog(targetURL)
		case "velog.io":
			return blogs.ProcessVelogBlog(targetURL)
		case "tistory.com":
			return blogs.ProcessTistoryBlog(targetURL)
		default:
			return blogs.BlogResponse{}, fmt.Errorf("unsupported host: %s", host)
		}
	} else if crawlingType == "shared_post_fetch" {
		switch host {
		case "medium.com":
			return blogs.ProcessMediumBlog(targetURL)
		case "velog.io":
			return blogs.ProcessVelogBlog(targetURL)
		case "tistory.com":
			return blogs.ProcessTistoryBlog(targetURL)
		default:
			return blogs.BlogResponse{}, fmt.Errorf("unsupported host: %s", host)
		}
	} else {
		return blogs.BlogResponse{}, fmt.Errorf("unsupported crawling type: %s", crawlingType)
	}
}
