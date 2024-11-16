import requests

# Define the endpoint
url = "https://medium.com/_/graphql"  # Replace with the actual URL

# Replace <USERNAME> with the actual username you want to query
username = "printSAN0"

# Define the query payload
payload = {
    "operationName": "GetUserId",
    "variables": {
        "username": username
    },
    "query": """
    query GetUserId($username: ID) {
      userResult(username: $username) {
        ... on User {
          id
        }
      }
    }
    """
}

# Make the request
headers = {
    "Content-Type": "application/json"
}

response = requests.post(url, json=payload, headers=headers)

# Print the response
print(response.json())
