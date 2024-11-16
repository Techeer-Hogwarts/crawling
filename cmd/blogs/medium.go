package blogs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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
	log.Printf("Found posts: %v", posts)
	return BlogResponse{}, nil
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

func getUserPostsMedium(userID string) ([]BlogResponse, error) {
	log.Printf("Getting posts for user ID: %s", userID)
	url := fmt.Sprintf(mediumURL, userID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v\n", err)
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Query failed with status code %d: %s\n", resp.StatusCode, body)
		return nil, fmt.Errorf("query failed with status code %d", resp.StatusCode)
	}
	// Parse and print the response JSON
	cleanedBody := removeUnwantedPrefix(string(body))
	var response interface{}
	if err := json.Unmarshal([]byte(cleanedBody), &response); err != nil {
		log.Printf("Error unmarshalling response JSON: %v\n", err)
		return nil, err
	}
	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(responseJSON))
	return nil, nil
}

func removeUnwantedPrefix(body string) string {
	prefix := `])}while(1);</x>`

	if strings.HasPrefix(body, prefix) {
		return body[len(prefix):]
	}
	return body
}
