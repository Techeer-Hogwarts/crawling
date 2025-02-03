package cmd

import (
	"fmt"
	"log"

	"github.com/Techeer-Hogwarts/crawling/cmd/blogs"
)

func CrawlBlog(targetURL, host, category string) (blogs.BlogResponse, error) {
	log.Println("Crawling blog:", targetURL)
	if category == "TECHEER" {
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
	} else if category == "SHARED" {
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
		return blogs.BlogResponse{}, fmt.Errorf("unsupported category type: %s", category)
	}
}
