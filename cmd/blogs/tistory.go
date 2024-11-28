package blogs

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// ProcessTistoryBlog processes Tistory blog sitemap.xml and returns the blog posts
func ProcessTistoryBlog(url string) (BlogResponse, error) {
	var cleanURL string
	if strings.Contains(url, "tistory.com") {
		cleanURL = strings.Split(url, ".tistory.com")[0] + ".tistory.com/rss"
	} else {
		return BlogResponse{}, fmt.Errorf("invalid tistory URL")
	}
	posts, err := getTistoryPosts(cleanURL)
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

func getTistoryPosts(url string) (BlogResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching sitemap: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected HTTP status: %s", resp.Status)
	}
	var tistoryResponse TistoryResponse

	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&tistoryResponse)
	if err != nil {
		log.Fatalf("Error decoding XML: %v", err)
	}
	authorProfileImage := tistoryResponse.Channel.Image.URL
	var tistoryBlogResponse BlogResponse
	for i, item := range tistoryResponse.Channel.Items {
		if i > 2 {
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
