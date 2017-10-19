package main

import (
	"container/list"
	"github.com/yfujita/nico-fav-tweet/nicorank"
	"github.com/yfujita/nico-fav-tweet/tweet"
	"log"
	"os"
	"strconv"
	"time"
	"fmt"
	"bufio"
	"io"
	"flag"
	"unicode/utf8"
)

const (
	MAX_DUPLICATE_COUNT = 100
	TWEET_LIMIT         = 1
	LATEST_MOVIES_FILE = "./latest.txt"
	LOG_PATH            = "/tmp/nicofav.log"
)

func main() {
	var (
		ck string
		cs string
		atoken string
		atoken_secret string
	)

	flag.StringVar(&ck, "ck", "", "set consumer key of twitter app")
	flag.StringVar(&cs, "cs", "", "set consumer secret")
	flag.StringVar(&atoken, "at", "", "set access token of twitter bot account")
	flag.StringVar(&atoken_secret, "as", "", "access token secret")
	flag.Parse()

	latestVideoLists := getLatestVideos(LATEST_MOVIES_FILE)
	nr := nicorank.NewNicoRank()
	ris, err := nr.Get()
	if err != nil {
		panic(err.Error())
	}

	tw := tweet.NewTweet()
	tw.SetUp(ck, cs, atoken, atoken_secret)

	logger := NewLogger()
	logger.Logging("start main task")

	count := 0
	for i, ri := range ris {
		logger.Logging(ri.Link)
		var exists = false
		for e := latestVideoLists.Front(); e != nil; e = e.Next() {
			if ri.Link == e.Value {
				exists = true
				break
			}
		}

		if !exists {
			message := ri.Title + " (" + ri.Point + " points) " + ri.Link
			if utf8.RuneCountInString(message) < 130 {
				message = message + " #ニコニコ動画"
			}
			if utf8.RuneCountInString(message) < 130 {
				message = message + " #" + ri.Id
			}

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

	updateLatestVideos(LATEST_MOVIES_FILE, latestVideoLists)
	err = followBack(tw)
	if err != nil {
		panic(err)
	}
	logger.Close()
}

func getLatestVideos(path string) *list.List {
	latestVideoLists := list.New()

	fp, err := os.Open(path)
	if err != nil && os.IsExist(err) {
		panic(err)
	}

	if err == nil {
		defer fp.Close()
		reader := bufio.NewReaderSize(fp, 4096)
		for {
			line, _, err := reader.ReadLine()
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			latestVideoLists.PushBack(string(line))
		}
	}

	return latestVideoLists
}

func updateLatestVideos(path string, latestVideoLists *list.List) {
	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
	}

	fp, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	writer:= bufio.NewWriter(fp)
	for e := latestVideoLists.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
		writer.WriteString(e.Value.(string) + "\n")
	}
	writer.Flush()
}

func followBack(tw *tweet.Tweet) error {
	friends, err := tw.Friends()
	if err != nil {
		return err
	}
	followers, err := tw.Followers()
	if err != nil {
		return err
	}

	for _, follower := range followers.Ids {
		find := false
		for _, friend := range friends.Ids {
			if follower == friend {
				find = true
				break
			}
		}
		if !find {
			log.Println("follow:" + strconv.FormatInt(follower, 10))
			err = tw.Follow(follower)
			if err != nil {
				log.Println(err)
			}
		}
	}

	for _, friend := range friends.Ids {
		find := false
		for _, follower := range followers.Ids {
			if follower == friend {
				find = true
				break
			}
		}
		if !find {
			log.Println("Unfollow:" + strconv.FormatInt(friend, 64))
			err = tw.Unfollow(friend)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return nil
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
	f, err := os.OpenFile(LOG_PATH, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
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
	fmt.Println(message)
}
