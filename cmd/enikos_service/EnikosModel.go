package main

import "encoding/xml"

// Rss was generated 2024-08-25 20:05:21 by https://xml-to-go.github.io/ in Ukraine.
type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Content string   `xml:"content,attr"`
	Wfw     string   `xml:"wfw,attr"`
	Dc      string   `xml:"dc,attr"`
	Atom    string   `xml:"atom,attr"`
	Sy      string   `xml:"sy,attr"`
	Slash   string   `xml:"slash,attr"`
	Media   string   `xml:"media,attr"`
	Channel struct {
		Text  string `xml:",chardata"`
		Title string `xml:"title"`
		Link  struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Description     string `xml:"description"`
		LastBuildDate   string `xml:"lastBuildDate"`
		Language        string `xml:"language"`
		UpdatePeriod    string `xml:"updatePeriod"`
		UpdateFrequency string `xml:"updateFrequency"`
		Generator       string `xml:"generator"`
		Item            []struct {
			Text     string   `xml:",chardata"`
			Title    string   `xml:"title"`
			Link     string   `xml:"link"`
			Comments string   `xml:"comments"`
			Creator  string   `xml:"creator"`
			PubDate  string   `xml:"pubDate"`
			Category []string `xml:"category"`
			Guid     struct {
				Text        string `xml:",chardata"`
				IsPermaLink string `xml:"isPermaLink,attr"`
			} `xml:"guid"`
			Description string `xml:"description"`
			CommentRss  string `xml:"commentRss"`
			Image       string `xml:"image"`
			Content     struct {
				Text   string `xml:",chardata"`
				URL    string `xml:"url,attr"`
				Medium string `xml:"medium,attr"`
				Width  string `xml:"width,attr"`
				Height string `xml:"height,attr"`
				Player struct {
					Text string `xml:",chardata"`
					URL  string `xml:"url,attr"`
				} `xml:"player"`
				Title struct {
					Text string `xml:",chardata"`
					Type string `xml:"type,attr"`
				} `xml:"title"`
				Description struct {
					Text string `xml:",chardata"`
					Type string `xml:"type,attr"`
				} `xml:"description"`
				Thumbnail struct {
					Text string `xml:",chardata"`
					URL  string `xml:"url,attr"`
				} `xml:"thumbnail"`
				Rating struct {
					Text   string `xml:",chardata"`
					Scheme string `xml:"scheme,attr"`
				} `xml:"rating"`
			} `xml:"content"`
			Enclosure struct {
				Text   string `xml:",chardata"`
				URL    string `xml:"url,attr"`
				Length string `xml:"length,attr"`
				Type   string `xml:"type,attr"`
			} `xml:"enclosure"`
		} `xml:"item"`
	} `xml:"channel"`
}
