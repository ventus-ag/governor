package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	layoutUS = "January 2, 2006, 15:04:05"
	topic    = "email"
)

type data struct {
	Max       int    `json:"max"`
	Current   int    `json:"current"`
	Treshold  int    `json:"treshold"`
	QuotaName string `json:"quota_name"`
	Email     string `json:"email"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Date      string `json:"date"`
}

type mail struct {
	Message string `json:"message"`
	Subject string `json:"subject"`
	Email   string `json:"email"`
}

type event struct {
	ID              string `json:"id"`
	Source          string `json:"source"`
	Type            string `json:"type"`
	Specversion     string `json:"specversion"`
	Datacontenttype string `json:"datacontenttype"`
	Data            data   `json:"data"`
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/dapr/subscribe", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]string{"cpu"})
	})

	router.HandleFunc("/cpu", messageHandler).Methods("POST")

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}

func messageHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
		return
	}
	var data event

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("New message - creating msg")

	templateFileName := "message.html"

	subject := "Ventus Cloud Support - Governor"

	// msg := "Dear User. We would like to inform you, that on project: " + data.Data.Name +
	// 	" you have passed treshhold in 60% of " + data.Data.QuotaName + " and using " + strconv.Itoa(data.Data.Current) +
	// 	" of " + strconv.Itoa(data.Data.Max) + " cores. If you want, you can scale up your resources here."

	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		log.Fatalln(err)
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data.Data); err != nil {
		log.Fatalln(err)
	}

	msg := buf.String()
	// log.Println("New msg:", msg)

	d := mail{
		Email:   data.Data.Email,
		Subject: subject,
		Message: msg,
	}

	publish(d)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Println("Published new e-mail")
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
