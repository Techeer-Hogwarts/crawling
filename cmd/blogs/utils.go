package blogs

type BlogResponse struct {
	UserID string       `json:"user_id"`
	Posts  []VelogPosts `json:"posts"`
}

type BlogRequest struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type VelogPosts struct {
	ID               string   `json:"id"`
	Title            string   `json:"title"`
	ShortDescription string   `json:"short_description"`
	Thumbnail        string   `json:"thumbnail"`
	URLSlug          string   `json:"url_slug"`
	ReleasedAt       string   `json:"released_at"`
	UpdatedAt        string   `json:"updated_at"`
	CommentsCount    int      `json:"comments_count"`
	Tags             []string `json:"tags"`
	Likes            int      `json:"likes"`
}

type VelogData struct {
	Posts []VelogPosts `json:"posts"`
}

type VelogResponse struct {
	Data VelogData `json:"data"`
}

type MediumUserResultWrapper struct {
	Data MediumUserResult `json:"data"`
}

type MediumUserResult struct {
	UserResult MediumUserId `json:"userResult"`
}

type MediumUserId struct {
	ID string `json:"id"`
}
