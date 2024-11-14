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

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type GraphQLResponse struct {
	Data   interface{} `json:"data"`
	Errors interface{} `json:"errors"`
}

func ProcessVelogBlog(url string) BlogResponse {
	apiurl := "https://v3.velog.io/graphql"
	var username string
	urlParts := strings.Split(url, "/")
	for _, part := range urlParts {
		if len(part) > 0 && part[0] == '@' {
			username = part[1:]
			break
		}
	}
	log.Printf("Processing Velog blog for user: %s", username)

	// Define the GraphQL query
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

	// Define the variables to be passed with the query
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"username": username,
			"tag":      "",
			"cursor":   "",
			"limit":    10,
		},
	}

	// Create the GraphQL request payload
	requestBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	// Convert the request payload to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
		return BlogResponse{}
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", apiurl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return BlogResponse{}
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return BlogResponse{}
	}
	defer resp.Body.Close()

	// Check if the request was successful
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return BlogResponse{}
	}
	var response VelogResponse
	if resp.StatusCode == http.StatusOK {
		// Parse and print the response JSON
		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Printf("Error unmarshalling response JSON: %v\n", err)
			return BlogResponse{}
		}
		// responseJSON, _ := json.MarshalIndent(response, "", "  ")
		// fmt.Println(string(responseJSON))
	} else {
		// Print the error
		fmt.Printf("Query failed with status code %d: %s\n", resp.StatusCode, body)
	}

	posts := response.Data
	return BlogResponse{
		UserID: username,
		Posts:  posts.Posts,
	}
}
