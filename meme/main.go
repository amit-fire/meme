package main

import (
	"fmt"
	"log"
	"net/http"

	c "meme/controller"
	r "meme/redis"
)

func startServer() {
	fmt.Println("server is using on port 7000")
	log.Fatal(http.ListenAndServe("localhost:7000", nil))
}

func main() {

	client := r.ConnectToRedis()

	// redis client should be injected (dependency injection)
	c.GetMeme(client)
	c.GetTokenBalance(client)

	startServer()
}
