package blogs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ProcessVelogBlog processes a Velog blog and returns the blog posts
func ProcessVelogBlog(url string, limit int) (BlogResponse, error) {
	apiurl := "https://v3.velog.io/graphql"
	var username string
	urlParts := strings.Split(url, "/")
	for _, part := range urlParts {
		if len(part) > 0 && part[0] == '@' {
			username = part[1:]
			break
		}
	}
	if username == "" {
		return BlogResponse{}, fmt.Errorf("invalid velog URL")
	}
	log.Printf("Processing Velog blog for user: %s", username)

	query := `
		query velogPosts($input: GetPostsInput!) {
			posts(input: $input) {
				id
				title
				short_description
				thumbnail
				user {
					id
					username
					profile {
						id
						thumbnail
						display_name
					}
				}
				url_slug
				released_at
				updated_at
				comments_count
				tags
				is_private
				likes
			}
		}
	`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"username": username,
			"tag":      "",
			"cursor":   "",
			"limit":    limit,
		},
	}

	requestBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %v\n", err)
		return BlogResponse{}, err
	}

	req, err := http.NewRequest("POST", apiurl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return BlogResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v\n", err)
		return BlogResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return BlogResponse{}, err
	}
	var response VelogResponse
	if resp.StatusCode == http.StatusOK {
		if err := json.Unmarshal(body, &response); err != nil {
			log.Printf("Error unmarshalling response JSON: %v\n", err)
			return BlogResponse{}, err
		}
	} else {
		fmt.Printf("Query failed with status code %d: %s\n", resp.StatusCode, body)
	}

	posts := response.Data
	for i, post := range posts.Posts {
		posts.Posts[i].URLSlug = fmt.Sprintf("https://velog.io/@%s/%s", username, post.URLSlug)
	}
	var velogBlogResponse BlogResponse
	for _, post := range posts.Posts {
		velogBlogResponse.Posts = append(velogBlogResponse.Posts, Posts{
			Title:       post.Title,
			URL:         post.URLSlug,
			Author:      post.User.Username,
			AuthorImage: post.User.Profile.Thumbnail,
			Thumbnail:   post.Thumbnail,
			Date:        convertDateTimeVelog(post.ReleasedAt),
			Tags:        post.Tags,
			Category:    "techeer",
		})
	}
	velogBlogResponse.BlogURL = fmt.Sprintf("https://velog.io/@%s", username)
	return velogBlogResponse, nil
}

func convertDateTimeVelog(dt string) string {
	parsedTime, err := time.Parse("2006-01-02T15:04:05.999Z", dt)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return "0000-00-00T00:00:00Z"
	}
	return parsedTime.Format("2006-01-02T15:04:05Z")
}

func ProcessSingleVelogBlog(blogURL string) (BlogResponse, error) {
	log.Printf("Processing single Velog blog for URL: %s", blogURL)
	posts, err := ProcessVelogBlog(blogURL, 40)
	if err != nil {
		return BlogResponse{}, err
	}
	originalURLDecoded, err := url.PathUnescape(blogURL)
	if err != nil {
		log.Printf("Error decoding URL: %v", err)
		return BlogResponse{}, err
	}
	newPosts := []Posts{}
	for i, post := range posts.Posts {
		if post.URL == originalURLDecoded || post.URL == blogURL {
			categoryFixedPost := posts.Posts[i]
			categoryFixedPost.Category = "shared"
			newPosts = []Posts{categoryFixedPost}
			break
		}
	}
	posts.Posts = newPosts
	return posts, nil
}
