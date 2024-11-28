package blogs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const graphQLURL = "https://medium.com/_/graphql"
const mediumURL = "https://medium.com/_/api/users/%v/profile/stream"

// ProcessMediumBlog processes a Medium blog and returns the blog posts
func ProcessMediumBlog(url string) (BlogResponse, error) {
	var username string
	urlParts := strings.Split(url, "/")
	for _, part := range urlParts {
		if len(part) > 0 && part[0] == '@' {
			username = part[1:]
			break
		}
	}
	if username == "" {
		return BlogResponse{}, fmt.Errorf("username not found in URL")
	}
	log.Printf("Processing Medium blog for user: %s", username)
	userID, err := getUserIdMedium(username)
	if err != nil {
		return BlogResponse{}, err
	}
	log.Printf("User ID: %s", userID)
	posts, err := getUserPostsMedium(userID)
	if err != nil {
		return BlogResponse{}, err
	}
	posts.BlogURL = fmt.Sprintf("https://medium.com/@%s", username)
	return posts, nil
}

func getUserIdMedium(username string) (string, error) {
	query := `
	query GetUserId($username: ID) {
	  userResult(username: $username) {
		... on User {
		  id
		}
	  }
	}`
	variables := map[string]interface{}{
		"username": username,
	}
	requestBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %v\n", err)
		return "", err
	}
	req, err := http.NewRequest("POST", graphQLURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return "", err
	}
	var response MediumUserResultWrapper
	if resp.StatusCode == http.StatusOK {
		// Parse and print the response JSON
		if err := json.Unmarshal(body, &response); err != nil {
			log.Printf("Error unmarshalling response JSON: %v\n", err)
			return "", err
		}
		// responseJSON, _ := json.MarshalIndent(response, "", "  ")
		// fmt.Println(string(responseJSON))
	} else {
		// Print the error
		fmt.Printf("Query failed with status code %d: %s\n", resp.StatusCode, body)
	}
	userID := response.Data.UserResult.ID
	return userID, nil
}

func getUserPostsMedium(userID string) (BlogResponse, error) {
	log.Printf("Getting posts for user ID: %s", userID)
	url := fmt.Sprintf(mediumURL, userID)
	req, err := http.NewRequest("GET", url, nil)
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return BlogResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Query failed with status code %d: %s\n", resp.StatusCode, body)
		return BlogResponse{}, fmt.Errorf("query failed with status code %d", resp.StatusCode)
	}
	// Parse and print the response JSON
	cleanedBody := removeUnwantedPrefix(string(body))
	var response MediumResponse
	if err := json.Unmarshal([]byte(cleanedBody), &response); err != nil {
		log.Printf("Error unmarshalling response JSON: %v\n", err)
		return BlogResponse{}, err
	}
	postIds := getMediumPostIds(response.Payload.StreamItems)
	mediumPostsFromResponse := response.Payload.References.Posts
	var finalPostResponse BlogResponse
	for _, postID := range postIds {
		post := mediumPostsFromResponse[postID]
		finalPostResponse.Posts = append(finalPostResponse.Posts, Posts{
			Title:       post.Title,
			URL:         fmt.Sprintf("https://medium.com/p/%s", post.URL),
			Date:        convertDateTimeMedium(post.Date),
			Author:      response.Payload.UserInfo.Name,
			AuthorImage: addMediumUserImage(response.Payload.UserInfo.ImageID),
			Thumbnail:   addMediumPostImage(post.Virtuals.PreviewImage.ImageID),
			Category:    "techeer",
			Tags:        processMediumTags(post.Virtuals.Tags),
		})
	}
	// responseJSON, _ := json.MarshalIndent(response, "", "  ")
	// fmt.Println(string(responseJSON))
	return finalPostResponse, nil
}

func removeUnwantedPrefix(body string) string {
	prefix := `])}while(1);</x>`

	if strings.HasPrefix(body, prefix) {
		return body[len(prefix):]
	}
	return body
}

func getMediumPostIds(streamItems []MediumStreamItems) []string {
	var postIds []string
	for _, item := range streamItems {
		if item.ItemType == "postPreview" {
			postIds = append(postIds, item.PostPreview.PostID)
		}
		if len(postIds) > 3 {
			break
		}
	}
	return postIds
}

func convertDateTimeMedium(date int64) string {
	seconds := date / 1000
	t := time.Unix(seconds, 0).UTC()
	iso8601Time := fmt.Sprint(t.Format("2006-01-02T15:04:05Z07:00"))
	return iso8601Time
}

func addMediumPostImage(imageID string) string {
	if imageID == "" {
		return ""
	}
	return fmt.Sprintf("https://miro.medium.com/v2/%s", imageID)
}

func addMediumUserImage(imageID string) string {
	if imageID == "" {
		return ""
	}
	return fmt.Sprintf("https://miro.medium.com/v2/%s", imageID)
}

func processMediumTags(tags []MediumTags) []string {
	var tagNames []string
	for _, tag := range tags {
		if tag.Type == "Tag" {
			tagNames = append(tagNames, tag.Name)
		}
	}
	if tagNames == nil {
		tagNames = []string{}
	}
	return tagNames
}
