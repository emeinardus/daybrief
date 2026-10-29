package rss

import "encoding/xml"

// rssDocument represents the root element of an RSS document.
type rssDocument struct {
	XMLName xml.Name   `xml:"rss"`
	Channel rssChannel `xml:"channel"`
}
