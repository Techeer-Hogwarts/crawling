import requests
import json

# Define the GraphQL endpoint
url = 'https://v3.velog.io/graphql'

# Define the GraphQL query
query = """
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
"""

# Define the variables to be passed with the query
variables = {
    "input": {
        "username": "hahnwoong",
        "tag": "",
        "cursor": "",
        "limit": 10
    }
}

payload = {
    "query": query,
    "variables": variables
}

# Send the request to the GraphQL API
response = requests.post(url, json=payload)

# Check if the request was successful
if response.status_code == 200:
    # Print the response JSON
    data = response.json()
    data = json.dumps(data, indent=2)
    print(data)
else:
    # Print the error
    print(f"Query failed with status code {response.status_code}: {response.text}")
