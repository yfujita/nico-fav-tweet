package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	MAX_LENGTH = 10
)

type NicoRank struct {
	Url string
}

type RankInfo struct {
	Title string
	Link  string
	Point string
}

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Item    []Item   `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
}

func NewNicoRank() *NicoRank {
	nr := new(NicoRank)
	nr.init()
	return nr
}

func (nr *NicoRank) Get() []*RankInfo {
	resp, err := http.Get(nr.Url)
	if err != nil {
		fmt.Printf("error\n")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return parse(body)
}

func (nr *NicoRank) init() {
	nr.Url = "http://www.nicovideo.jp/ranking/fav/hourly/all?rss=2.0"
}

func parse(r []byte) []*RankInfo {
	ris := make([]*RankInfo, MAX_LENGTH, MAX_LENGTH)
	rss := Rss{}
	xml.Unmarshal(r, &rss)
	pointTag := "nico-info-number\">"
	for i, item := range rss.Channel.Item {
		if i >= MAX_LENGTH {
			break
		}
		ri := new(RankInfo)

		rankIndex := strings.Index(item.Title, "位：")
		ri.Title = item.Title[rankIndex+len("位："):]

		ri.Link = item.Link

		pointIndex := strings.Index(item.Description, pointTag)
		substr := item.Description[pointIndex+len(pointTag):]
		pointLastIndex := strings.Index(substr, "<")
		ri.Point = substr[:pointLastIndex]

		ris[i] = ri
	}
	return ris
}
