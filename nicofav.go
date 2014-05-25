package main

import (
	"container/list"
	"github.com/yfujita/nico-fav-tweet/nicorank"
	"github.com/yfujita/nico-fav-tweet/tweet"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	INTERVAL            = 30
	MAX_DUPLICATE_COUNT = 100
	TWEET_LIMIT         = 5
	LOG_PATH            = "/tmp/nicofav.log"
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
			ris, err := nr.Get()
			if err != nil {

			}
			ch <- ris
			time.Sleep(INTERVAL * time.Minute)
		}
	}(ch)

	tw := tweet.NewTweet()
	tw.SetUp(CONSUMER_KEY, CONSUMER_SECRET, ATOKEN, ATOKEN_SECRET)

	for {
		ris := <-ch
		logger := NewLogger()
		logger.Logging("start main task")

		count := 0
		for i, ri := range ris {
			logger.Logging(ri.Link)
			var exists bool = false
			for e := latestVideoLists.Front(); e != nil; e = e.Next() {
				if ri.Link == e.Value {
					exists = true
					break
				}
			}

			if !exists {
				message := ri.Title + " (" + ri.Point + " points) " + ri.Link
				logger.Logging(message)
				err := tw.Message(message)
				if err != nil {
					logger.Logging("Failed to tweet message: " + message)
				}

				if MAX_DUPLICATE_COUNT < latestVideoLists.Len() {
					e := latestVideoLists.Front()
					latestVideoLists.Remove(e)
				}
				latestVideoLists.PushBack(ri.Link)
				logger.Logging("dup lists size=" + strconv.FormatInt(int64(latestVideoLists.Len()), 10))

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
		logger.Close()
	}
}

type Logger struct {
	file *os.File
}

func NewLogger() *Logger {
	if len(LOG_PATH) <= 0 {
		return nil
	}

	if _, err := os.Stat(LOG_PATH); err != nil {
		if os.IsNotExist(err) {
			fo, err := os.Create(LOG_PATH)
			if err != nil {
				return nil
			}
			fo.Close()
		}
	}

	lg := new(Logger)
	f, err := os.OpenFile(LOG_PATH, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil
	}

	lg.file = f
	return lg
}

func (lg *Logger) Close() {
	lg.file.Close()
	lg.file = nil
}

func (lg *Logger) Logging(str string) {
	log.SetOutput(lg.file)

	message := "[" + time.Now().Format(time.RFC3339) + "] " + str
	log.Println(message)
}
