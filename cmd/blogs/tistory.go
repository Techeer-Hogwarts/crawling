package blogs

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ProcessTistoryBlog processes Tistory blog sitemap.xml and returns the blog posts
func ProcessTistoryBlog(url string, limit int) (BlogResponse, error) {
	var cleanURL string
	if strings.Contains(url, ".tistory.com") {
		cleanURL = strings.Split(url, ".tistory.com")[0] + ".tistory.com/rss"
	} else {
		return BlogResponse{}, fmt.Errorf("invalid tistory URL")
	}
	posts, err := getTistoryPosts(cleanURL, limit)
	posts.BlogURL = url
	if err != nil {
		return BlogResponse{}, err
	}
	// jsonResponse, err := json.MarshalIndent(posts, "", "  ")
	// if err != nil {
	// 	log.Fatalf("Error marshalling JSON: %v", err)
	// }
	// fmt.Println(string(jsonResponse))
	return posts, nil
}

func getTistoryPosts(url string, limit int) (BlogResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching sitemap: %v", err)
		return BlogResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected HTTP status: %s", resp.Status)
		return BlogResponse{}, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}
	var tistoryResponse TistoryResponse

	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&tistoryResponse)
	if err != nil {
		log.Printf("Error decoding XML: %v", err)
		return BlogResponse{}, err
	}
	authorProfileImage := tistoryResponse.Channel.Image.URL
	var tistoryBlogResponse BlogResponse
	for i, item := range tistoryResponse.Channel.Items {
		if i > limit {
			break
		}
		tistoryBlogResponse.Posts = append(tistoryBlogResponse.Posts, Posts{
			Title:       item.Title,
			URL:         item.Link,
			Date:        convertDateTimeTistory(item.PubDate),
			Author:      item.Author,
			AuthorImage: authorProfileImage,
			Category:    "techeer",
			Tags:        []string{},
		})
	}
	return tistoryBlogResponse, nil
}

func convertDateTimeTistory(dateString string) string {
	layout := "Mon, 2 Jan 2006 15:04:05 -0700"
	parsedTime, err := time.Parse(layout, dateString)
	if err != nil {
		log.Println("Error parsing date:", err)
		return "0000-00-00T00:00:00Z"
	}
	utcTime := parsedTime.UTC()
	return fmt.Sprint(utcTime.Format("2006-01-02T15:04:05Z07:00"))
}

func ProcessSingleTistoryBlog(blogURL string) (BlogResponse, error) {
	log.Printf("Processing single Tistory blog for URL: %s", blogURL)
	posts, err := ProcessTistoryBlog(blogURL, 40)
	if err != nil {
		return BlogResponse{}, err
	}
	originalURLDecoded, err := url.PathUnescape(blogURL)
	if err != nil {
		log.Printf("Error decoding URL: %v", err)
		return BlogResponse{}, err
	}
	newPosts := []Posts{}
	// Only leave one post with exact URL match (single post)
	for i, post := range posts.Posts {
		if post.URL == originalURLDecoded || post.URL == blogURL {
			categoryFixedPost := posts.Posts[i]
			categoryFixedPost.Category = "shared"
			newPosts = []Posts{categoryFixedPost}
			break
		}
	}
	posts.Posts = newPosts
	// jsonResponse, err := json.MarshalIndent(posts, "", "  ")
	// if err != nil {
	// 	log.Fatalf("Error marshalling JSON: %v", err)
	// }
	// fmt.Println(string(jsonResponse))
	return posts, nil
}
