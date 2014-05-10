package main

import (
	"container/list"
	"fmt"
	"github.com/yfujita/nico-fav-tweet/nicorank"
	"github.com/yfujita/nico-fav-tweet/tweet"
	"time"
)

const (
	INTERVAL            = 30
	MAX_DUPLICATE_COUNT = 100
	TWEET_LIMIT         = 5
	CONSUMER_KEY        = "set consumer key of twitter app"
	CONSUMER_SECRET     = "set consumer secret"
	ATOKEN              = "set access token of twitter bot account"
	ATOKEN_SECRET       = "set access token secret"
)

func main() {
	latestVideoLists := list.New()

	ch := make(chan []*nicorank.RankInfo, 10)
	go func(ch chan []*nicorank.RankInfo) {
		nr := nicorank.NewNicoRank()
		for {
			fmt.Println("start nico task")
			ris := nr.Get()
			ch <- ris
			time.Sleep(INTERVAL * time.Minute)
		}
	}(ch)

	tw := tweet.NewTweet()
	tw.SetUp(CONSUMER_KEY, CONSUMER_SECRET, ATOKEN, ATOKEN_SECRET)

	for {
		ris := <-ch
		fmt.Println(time.Now())
		fmt.Println("start main task")

		count := 0
		for i, ri := range ris {
			fmt.Println(ri.Link)
			var exists bool = false
			for e := latestVideoLists.Front(); e != nil; e = e.Next() {
				if ri.Link == e.Value {
					exists = true
					break
				}
			}

			if !exists {
				message := "@meumeu69 " + ri.Title + " (" + ri.Point + " points) " + ri.Link
				fmt.Println(message)
				tw.Message(message)
				if MAX_DUPLICATE_COUNT < latestVideoLists.Len() {
					e := latestVideoLists.Front()
					latestVideoLists.Remove(e)
				}
				latestVideoLists.PushBack(ri.Link)
				fmt.Println("dup lists size=%d", latestVideoLists.Len())

				count++
				if i > TWEET_LIMIT {
					count++
				}
				if count >= TWEET_LIMIT {
					break
				}
				time.Sleep(1 * time.Second)
			}
		}
	}
}
