package main

import "encoding/xml"

type RssNewsBeast struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Media   string   `xml:"media,attr"`
	Version string   `xml:"version,attr"`
	Content string   `xml:"content,attr"`
	Wfw     string   `xml:"wfw,attr"`
	Dc      string   `xml:"dc,attr"`
	Atom    string   `xml:"atom,attr"`
	Sy      string   `xml:"sy,attr"`
	Slash   string   `xml:"slash,attr"`
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
			Text      string `xml:",chardata"`
			Title     string `xml:"title"`
			Link      string `xml:"link"`
			Thumbnail struct {
				Text  string `xml:",chardata"`
				URL   string `xml:"url,attr"`
				Width string `xml:"width,attr"`
			} `xml:"thumbnail"`
			Enclosure struct {
				Text   string `xml:",chardata"`
				URL    string `xml:"url,attr"`
				Type   string `xml:"type,attr"`
				Length string `xml:"length,attr"`
			} `xml:"enclosure"`
			PubDate     string   `xml:"pubDate"`
			Category    []string `xml:"category"`
			Description string   `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}
