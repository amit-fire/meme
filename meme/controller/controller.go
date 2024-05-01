package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"

	s "meme/service"
)

type Response struct {
	Message string `json:"meme"`
}

func printTime(id string) {
	fmt.Println("id " + id + " " + time.Now().Format("2006/01/02 15:04:05:000"))
}

func GetMeme(client *redis.Client) {
	// http request handler
	http.HandleFunc("/app/memes", func(w http.ResponseWriter, r *http.Request) {

		id := r.Header.Get("id")
		printTime(id)
		if id == "" {
			http.Error(w, "missing 'id' header", http.StatusBadRequest)
			return
		}

		params := r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		meme, err := s.GetMeme(client, id, params.Get("lat"), params.Get("lon"), params.Get("query"))

		/*
			TODO: add retry mechanism with backoff policy, in case there are network issues.
			If after X aomunt of times the connection fails, place the meme in a queue that it
			could be sent later.
		*/

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			json.NewEncoder(w).Encode(Response{Message: meme})
			s.IncreaseRequestsCount(client, s.GetRequestCountId(id))
		}
	})
}

func GetTokenBalance(client *redis.Client) {
	// http request handler
	http.HandleFunc("/app/balance/", func(w http.ResponseWriter, r *http.Request) {

		pathParts := strings.Split(r.URL.Path, "/")
		id := pathParts[len(pathParts)-1]
		val, tokenErr := s.GetTokenBalance(client, id)

		msg := ""
		if tokenErr != nil {
			msg = tokenErr.Error()
		} else {
			msg = "id " + id + " has " + strconv.Itoa(val) + " tokens"
		}

		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(Response{Message: msg})
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	})

}
