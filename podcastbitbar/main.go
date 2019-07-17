/*
Copyright 2019 github.com/ucirello and cirello.io/podcastbitbar

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	targetURL := flag.String("url", "", "")
	flag.Parse()
	resp, err := http.Get(*targetURL)
	fatalOnErr(err)
	data, err := ioutil.ReadAll(resp.Body)
	fatalOnErr(err)
	resp.Body.Close()
	var feed feed
	fatalOnErr(xml.Unmarshal(data, &feed))
	if len(feed.Channel) != 1 {
		log.Fatal("not enough channels")
	}
	podcast := feed.Channel[0]
	item := podcast.Items[0]
	fmt.Println("🎧")
	fmt.Println("---")
	fmt.Println(podcast.Title, "-", item.Title, ` | href="`+item.Enclosure.URL+`"`)
}

func fatalOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type feed struct {
	XMLName xml.Name  `xml:"rss"`
	Channel []podcast `xml:"channel"`
}

type podcast struct {
	XMLName        xml.Name `xml:"channel"`
	Title          string   `xml:"title"`
	Link           string   `xml:"link"`
	Description    string   `xml:"description"`
	Category       string   `xml:"category,omitempty"`
	Cloud          string   `xml:"cloud,omitempty"`
	Copyright      string   `xml:"copyright,omitempty"`
	Docs           string   `xml:"docs,omitempty"`
	Generator      string   `xml:"generator,omitempty"`
	Language       string   `xml:"language,omitempty"`
	LastBuildDate  string   `xml:"lastBuildDate,omitempty"`
	ManagingEditor string   `xml:"managingEditor,omitempty"`
	PubDate        string   `xml:"pubDate,omitempty"`
	Rating         string   `xml:"rating,omitempty"`
	SkipHours      string   `xml:"skipHours,omitempty"`
	SkipDays       string   `xml:"skipDays,omitempty"`
	TTL            int      `xml:"ttl,omitempty"`
	WebMaster      string   `xml:"webMaster,omitempty"`

	Items []item `xml:"item"`
}

type item struct {
	XMLName          xml.Name  `xml:"item"`
	GUID             string    `xml:"guid"`
	Title            string    `xml:"title"`
	Link             string    `xml:"link"`
	Description      string    `xml:"description"`
	AuthorFormatted  string    `xml:"author,omitempty"`
	PubDateFormatted string    `xml:"pubDate,omitempty"`
	Enclosure        enclosure `xml:"enclosure"`
}

type enclosure struct {
	XMLName         xml.Name `xml:"enclosure"`
	URL             string   `xml:"url,attr"`
	Length          int64    `xml:"-"`
	LengthFormatted string   `xml:"length,attr"`
	TypeFormatted   string   `xml:"type,attr"`
}
