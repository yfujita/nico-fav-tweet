package main

import (
	"github.com/mrjones/oauth"
)

type Tweet struct {
	consumerKey    string
	consumerSecret string
	accessToken    *oauth.AccessToken
	consumer       *oauth.Consumer
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

func (tw *Tweet) Message(message string) {
	response, err := tw.consumer.Post(
		"https://api.twitter.com/1.1/statuses/update.json",
		map[string]string{
			"status": message,
		},
		tw.accessToken)

	if err != nil {
		//ignore
	}

	if response != nil {
		//ignore
	}
}
