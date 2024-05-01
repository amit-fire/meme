package meme_service

import (
	"errors"

	"github.com/go-redis/redis"
)

func GetRequestCountId(id string) string {
	return id + "_count"
}

func GetMeme(client *redis.Client, id string, lat string, lon string, query string) (string, error) {
	proceed, err := canProceed(client, id)

	if err != nil {
		return "", err
	}

	if proceed {
		// TODO: validate lat, lon, query
		meme := getMeme(id, lat, lon, query)
		return meme, nil
	}

	return "", errors.New("could not generate meme for id " + id)
}

func canProceed(client *redis.Client, id string) (bool, error) {
	tokens, err := getTokens(client, id)
	if err != nil {
		return false, err
	}

	/*
		the user has at least one token, which means there's a point to continue checking
		if the number of times the api has been called is below the number of tokens the user has
	*/
	if tokens > 0 {
		/*
			in a real-life scenario, the tokens and request count would be in separate db instances.
			for simplicity, a single instance is used.
			to distinguish between the keys, the suffix "_count" is added to the id for
			tracking the number of times a call has reached the api.
		*/
		count, err := getRequestsCount(client, GetRequestCountId(id))

		if err == nil && count < tokens {
			return true, nil
		}
	}

	return false, errors.New(id + " does not have tokens")
}

func getTokens(client *redis.Client, id string) (int, error) {
	val, err := client.Get(id).Int()
	if err != nil {
		return -1, errors.New(id + " does not exist")
	}

	return val, nil
}

func getRequestsCount(client *redis.Client, id string) (int, error) {
	item := client.Get(id)

	// first request. count key doesn't exist
	if item.Val() == "" {
		return 0, nil
	} else {
		return item.Int()
	}
}

func IncreaseRequestsCount(client *redis.Client, id string) {
	client.Incr(id)
}

func getMeme(id string, lat string, lon string, query string) string {
	return id + "_" + lat + "_" + lon + "_" + query
}

func GetTokenBalance(client *redis.Client, id string) (int, error) {
	tokens, err := client.Get(id).Int()
	if err != nil {
		return -1, errors.New(id + " does not exist")
	}

	count, err := getRequestsCount(client, GetRequestCountId(id))

	if err != nil {
		return -1, errors.New("error fetching data for id " + id)
	}

	return tokens - count, nil
}

func GetTokenBalance2(client *redis.Client, id string) (int, error) {
	val, err := client.Get(id).Int()
	if err != nil {
		return -1, errors.New(id + " does not exist")
	}

	return val, nil
}
