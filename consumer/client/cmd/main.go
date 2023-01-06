package main

import (
	"github.com/doruk581/cdc-wholesale-workshop-go/consumer/client"
	"log"
	"net/url"
	"time"
)

var token = time.Now().Format("2006-01-02T15:04")

func main() {
	u, _ := url.Parse("http://localhost:8080")
	client := &client.Client{
		BaseURL: u,
	}

	users, err := client.WithToken(token).GetProduct(10)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(users)
}
