package cmd

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

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

type BlogPosts struct {
	Title     string   `json:"title"`
	Text      string   `json:"text"`
	Thumbnail string   `json:"thumbnail"`
	Link      string   `json:"link"`
	Tags      []string `json:"tags"`
	Date      string   `json:"date"`
}

type BlogResponse struct {
	UserID string      `json:"user_id"`
	Posts  []BlogPosts `json:"posts"`
}

type BlogRequest struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

func CrawlBlog(targetURL string) (BlogResponse, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
	)
	allocatorCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	if cancel != nil {
		defer cancel()
	}
	ctx, cancel := chromedp.NewContext(allocatorCtx)
	if cancel != nil {
		defer cancel()
	}
	var divContents []DivContent
	err := chromedp.Run(ctx,
		chromedp.Navigate(targetURL),
		chromedp.WaitReady("body"),
	)
	if err != nil {
		log.Fatal(err)
	}
	var prevHeight, currHeight int
	for i := 0; i < 4; i++ { // 최대 40 언저리
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`document.body.scrollHeight`, &prevHeight),
		)
		if err != nil {
			log.Fatal(err)
			return BlogResponse{}, err
		}

		err = chromedp.Run(ctx,
			chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight)`, nil),
			chromedp.Sleep(2*time.Second),
		)
		if err != nil {
			log.Fatal(err)
			return BlogResponse{}, err
		}

		err = chromedp.Run(ctx,
			chromedp.Evaluate(`document.body.scrollHeight`, &currHeight),
		)
		if err != nil {
			log.Fatal(err)
			return BlogResponse{}, err
		}

		if prevHeight == currHeight {
			log.Println("No more new content loaded. Stopping scroll.")
			break
		}
	}
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('div.FlatPostCard_block__a1qM7')).map(div => {
				const getText = (selector) => {
					const el = div.querySelector(selector);
					return el ? el.textContent : "";
				};

				const getAttribute = (selector, attribute) => {
					const el = div.querySelector(selector);
					return el ? el.getAttribute(attribute) : "";
				};

				return {
					html: div.outerHTML,
					links: Array.from(div.querySelectorAll('a')).map(a => ({
						href: a.getAttribute('href') || "",
						text: a.textContent || ""
					})),
					images: Array.from(div.querySelectorAll('img')).map(img => ({
						src: img.getAttribute('src') || "",
						alt: img.getAttribute('alt') || "",
						fetchpriority: img.getAttribute('fetchpriority') || "",
						decoding: img.getAttribute('decoding') || "",
						data_nimg: img.getAttribute('data-nimg') || "",
						style: img.getAttribute('style') || ""
					})),
					tags: Array.from(div.querySelectorAll('div.FlatPostCard_tagsWrapper__iNQR3 a')).map(a => ({
						href: a.getAttribute('href') || "",
						text: a.textContent || ""
					})),
					sub_info: {
						date: getText('div.FlatPostCard_subInfo__cT3J6 span:nth-child(1)'),
						comments: getText('div.FlatPostCard_subInfo__cT3J6 span:nth-child(4)'),
						likes: getText('div.FlatPostCard_subInfo__cT3J6 span.FlatPostCard_likes__TtpEU')
					},
					text: getText('p')
				}
			})
		`, &divContents),
	)
	if err != nil {
		log.Fatal(err)
	}

	var blogPosts []BlogPosts

	for _, divContent := range divContents {
		var blogs BlogPosts
		blogs.Text = divContent.Text
		if len(divContent.Images) > 0 {
			blogs.Thumbnail = divContent.Images[0].Src
		}

		if len(divContent.Links) > 1 {
			blogs.Link = divContent.Links[1].Href
		} else if len(divContent.Links) == 1 {
			blogs.Link = divContent.Links[0].Href
		} else {
			blogs.Link = ""
		}

		if len(divContent.Links) > 1 {
			blogs.Title = divContent.Links[1].Text
		} else if len(divContent.Links) == 1 {
			blogs.Title = divContent.Links[0].Text
		} else {
			blogs.Title = ""
		}

		blogs.Date = divContent.SubInfo.Date
		for _, tag := range divContent.Tags {
			blogs.Tags = append(blogs.Tags, strings.ToLower(tag.Text))
		}
		blogPosts = append(blogPosts, blogs)
	}
	return BlogResponse{
		UserID: "test",
		Posts:  blogPosts,
	}, nil
}
