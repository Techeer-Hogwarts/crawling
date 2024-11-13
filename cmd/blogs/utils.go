package blogs

type Image struct {
	Src           string `json:"src"`
	Alt           string `json:"alt"`
	FetchPriority string `json:"fetchpriority"`
	Decoding      string `json:"decoding"`
	DataNimg      string `json:"data-nimg"`
	Style         string `json:"style"`
}

type Link struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

type Tag struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

type SubInfo struct {
	Date     string `json:"date"`
	Comments string `json:"comments"`
	Likes    string `json:"likes"`
}

type DivContent struct {
	HTML    string  `json:"html"`
	Links   []Link  `json:"links"`
	Images  []Image `json:"images"`
	Tags    []Tag   `json:"tags"`
	SubInfo SubInfo `json:"sub_info"`
	Text    string  `json:"text"`
}

type BlogResponse struct {
	UserID string       `json:"user_id"`
	Posts  []VelogPosts `json:"posts"`
}

type BlogRequest struct {
	Type string `json:"type"`
	Data string `json:"data"`
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
