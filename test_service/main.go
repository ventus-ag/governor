package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const (
	layoutUS = "January 2, 2006, 15:04:05"
	topic    = "email"
)

type mail struct {
	Addr    string `json:"addr"`
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
}

func main() {
	for true {
		log.Println("START")

		d := mail{
			Addr:    "smtp.office365.com:587",
			From:    "dmitriy.yarovoy@ventus.ag",
			To:      "masterhorn89@gmail.com",
			Message: "We will be happy to send you notifications when you need it.",
		}

		publish(d)

		log.Println("go sleep")
		time.Sleep(time.Minute * 5)
	}

}

func publish(d mail) {
	body, err := json.Marshal(d)
	if err != nil {
		log.Fatalln(err)
	}

	URL := "http://localhost:3500/v1.0/publish/" + topic
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 500 {
		log.Fatalln("500 on publishing new message into the topic")
	} else {
		log.Println("200 on publishing new message into the topic")
	}
}

//test
