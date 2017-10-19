package tweet

import (
	"github.com/mrjones/oauth"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

type Tweet struct {
	consumerKey    string
	consumerSecret string
	accessToken    *oauth.AccessToken
	consumer       *oauth.Consumer
}

type Friends struct {
	Ids				[]int64
}

type Followers struct {
	Ids				[]int64
}

func NewTweet() *Tweet {
	tw := new(Tweet)
	tw.accessToken = new(oauth.AccessToken)
	return tw
}

func (tw *Tweet) SetUp(consumerKey, consumerSecret, aToken, aTokenSecret string) {
	tw.consumerKey = consumerKey
	tw.consumerSecret = consumerSecret
	tw.accessToken.Token = aToken
	tw.accessToken.Secret = aTokenSecret

	tw.consumer = oauth.NewConsumer(
		tw.consumerKey,
		tw.consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		})
}

func (tw *Tweet) Message(message string) error {
	response, err := tw.consumer.Post(
		"https://api.twitter.com/1.1/statuses/update.json",
		map[string]string{
			"status": message,
		},
		tw.accessToken)

	if err != nil {
		return err
	}

	if response != nil {
		//ignore
	}
	return nil

}

func (tw *Tweet) Friends() (*Friends, error) {
	response, err := tw.consumer.Get(
		"https://api.twitter.com/1.1/friends/ids.json",
		map[string]string{
		},
		tw.accessToken)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	friends := new(Friends)
	err = json.Unmarshal(body, friends)
	if err != nil {
		panic(err)
	}
	return friends, nil
}

func (tw *Tweet) Followers() (*Followers, error) {
	response, err := tw.consumer.Get(
		"https://api.twitter.com/1.1/followers/ids.json",
		map[string]string{
		},
		tw.accessToken)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	followers := new(Followers)
	json.Unmarshal(body, followers)
	return followers, nil
}

func (tw *Tweet) Follow(id int64) error {
	response, err := tw.consumer.Post(
		"https://api.twitter.com/1.1/friendships/create.json",
		map[string]string{
			"user_id" : strconv.FormatInt(id, 10),
		},
		tw.accessToken)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	return nil
}

func (tw *Tweet) Unfollow(id int64) error {
	_, err := tw.consumer.Post(
		"https://api.twitter.com/1.1/friendships/destroy.json",
		map[string]string{
			"user_id" : strconv.FormatInt(id, 10),
		},
		tw.accessToken)
	if err != nil {
		return err
	}
	return nil
}