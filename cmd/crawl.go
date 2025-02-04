package cmd

import (
	"fmt"
	"log"

	"github.com/Techeer-Hogwarts/crawling/cmd/blogs"
)

func CrawlBlog(targetURL, host, crawlingType string) (blogs.BlogResponse, error) {
	log.Println("Crawling blog:", targetURL)
	if crawlingType == "signUp_blog_fetch" || crawlingType == "blogs_daily_update" {
		switch host {
		case "medium.com":
			return blogs.ProcessMediumBlog(targetURL, 3)
		case "velog.io":
			return blogs.ProcessVelogBlog(targetURL, 3)
		case "tistory.com":
			return blogs.ProcessTistoryBlog(targetURL, 2)
		default:
			return blogs.BlogResponse{}, fmt.Errorf("unsupported host: %s", host)
		}
	} else if crawlingType == "shared_post_fetch" {
		switch host {
		case "medium.com":
			return blogs.ProcessSingleMediumBlog(targetURL)
		case "velog.io":
			return blogs.ProcessSingleVelogBlog(targetURL)
		case "tistory.com":
			return blogs.ProcessSingleTistoryBlog(targetURL)
		default:
			return blogs.BlogResponse{}, fmt.Errorf("unsupported host: %s", host)
		}
	} else {
		return blogs.BlogResponse{}, fmt.Errorf("unsupported category type: %s", crawlingType)
	}
}
