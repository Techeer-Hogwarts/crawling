package blogs

// for general blos request and response
type BlogResponse struct {
	UserID  int     `json:"userID"`
	BlogURL string  `json:"blogURL"`
	Posts   []Posts `json:"posts"`
}

type Posts struct {
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Author      string   `json:"author"`
	AuthorImage string   `json:"authorImage"`
	Thumbnail   string   `json:"thumbnail"`
	Category    string   `json:"category"`
	Date        string   `json:"date"`
	Tags        []string `json:"tags"`
}

type BlogRequest struct {
	UserID   int    `json:"userID"`
	Data     string `json:"data"`
	Category string `json:"category"`
}

// for velog
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type VelogPosts struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	ShortDescription string    `json:"short_description"`
	Thumbnail        string    `json:"thumbnail"`
	URLSlug          string    `json:"url_slug"`
	ReleasedAt       string    `json:"released_at"`
	UpdatedAt        string    `json:"updated_at"`
	CommentsCount    int       `json:"comments_count"`
	Tags             []string  `json:"tags"`
	Likes            int       `json:"likes"`
	User             VelogUser `json:"user"`
}

type VelogUser struct {
	ID       string       `json:"id"`
	Username string       `json:"username"`
	Profile  VelogProfile `json:"profile"`
}

type VelogProfile struct {
	ID          string `json:"id"`
	Thumbnail   string `json:"thumbnail"`
	DisplayName string `json:"display_name"`
}

type VelogData struct {
	Posts []VelogPosts `json:"posts"`
}

type VelogResponse struct {
	Data VelogData `json:"data"`
}

// for medium
type MediumUserResultWrapper struct {
	Data MediumUserResult `json:"data"`
}

type MediumUserResult struct {
	UserResult MediumUserId `json:"userResult"`
}

type MediumUserId struct {
	ID string `json:"id"`
}

type MediumResponse struct {
	Payload MediumPayload `json:"payload"`
}

type MediumPayload struct {
	StreamItems []MediumStreamItems `json:"streamItems"`
	References  MediumReference     `json:"references"`
	UserInfo    MediumUserInfo      `json:"user"`
}

type MediumReference struct {
	Posts map[string]MediumPosts `json:"Post"`
}

type MediumStreamItems struct {
	ItemType    string            `json:"itemType"`
	PostPreview MediumPostPreview `json:"postPreview"`
}

type MediumPostPreview struct {
	PostID string `json:"postId"`
}

type MediumPosts struct {
	Title    string         `json:"title"`
	URL      string         `json:"uniqueSlug"`
	Date     int64          `json:"firstPublishedAt"`
	Virtuals MediumVirtuals `json:"virtuals"`
}

type MediumVirtuals struct {
	Tags         []MediumTags       `json:"tags"`
	PreviewImage MediumPreviewImage `json:"previewImage"`
}

type MediumTags struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type MediumPreviewImage struct {
	ImageID string `json:"imageId"`
}

type MediumUserInfo struct {
	Name    string `json:"name"`
	ImageID string `json:"imageId"`
}

// for tistory
type TistoryResponse struct {
	Channel TistoryChannel `xml:"channel" json:"channel"`
}

type TistoryChannel struct {
	Image TistoryImage  `xml:"image" json:"image"`
	Items []TistoryItem `xml:"item" json:"item"`
}

type TistoryImage struct {
	URL string `xml:"url" json:"url"`
}

type TistoryItem struct {
	Title   string `xml:"title" json:"title"`
	Link    string `xml:"link" json:"link"`
	PubDate string `xml:"pubDate" json:"pubDate"`
	Author  string `xml:"author" json:"author"`
}
