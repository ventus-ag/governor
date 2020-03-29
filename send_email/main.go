package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/gorilla/mux"
)

const (
	layoutUS = "January 2, 2006, 15:04:05"
)

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
	Data            mail   `json:"data"`
}

type loginAuth struct {
	username, password string
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/dapr/subscribe", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]string{"email"})
	})

	router.HandleFunc("/email", messageHandler).Methods("POST")

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

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	pwd := os.Getenv("VENTUS_EMAIL_PASSWORD")
	from := os.Getenv("NO_REPLY_EMAIL")
	auth := LoginAuth(from, pwd)
	to := data.Data.Email
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + data.Data.Subject + "\n" + mime + "\n" +
		data.Data.Message

	err = smtp.SendMail("smtp.office365.com:587", auth,
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Fatalf("smtp error: %s", err)
	}

	log.Println("Payload: \n" + string(body))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// LoginAuth is used for smtp login auth
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}
