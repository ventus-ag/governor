package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
)

var (
	layoutUS       = "January 2, 2006, 15:04:05"
	stateStoreName = "statestore"
	stateURL       = "http://localhost:3500/v1.0/state/" + stateStoreName
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

type state struct {
	Key   string     `json:"key"`
	Value stateValue `json:"value"`
}

type stateValue struct {
	Mail bool   `json:"mail"`
	Date string `json:"date"`
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/id", getData).Methods("GET")
	router.HandleFunc("/order", getOrder).Methods("GET")
	router.HandleFunc("/neworder", newOrder).Methods("GET")

	srv := &http.Server{
		Handler: router,
		Addr:    ":80",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Starting server at port 80")
	log.Fatal(srv.ListenAndServe())

}

func getData(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Fatalln("Error parsing parameters: ", err)
	}
	pID := params["projectid"][0]
	log.Println(pID)

	req, err := http.NewRequest("GET", stateURL+"/"+pID, nil)
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(bodyBytes)
	if err != nil {
		log.Fatalln(err)
	}
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("Got new get order request.")
	url := stateURL + "/order"
	url = "http://localhost:3500/v1.0/state/statestore/order"
	log.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 204 {
		log.Println("Status code 204")
	} else if resp.StatusCode == 200 {
		log.Println("Status code 200")
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	log.Println(bodyBytes)
	log.Println(string(bodyBytes))
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(bodyBytes)
	// err = json.NewEncoder(w).Encode(bodyBytes)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Finished with get order request.")
}

func newOrder(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Fatalln("Error parsing parameters: ", err)
	}
	order := params["order"][0]
	log.Println("Got new order: " + order)

	data := `[{"key": "order","value": "` + order + `"}]`
	req, err := http.NewRequest("POST", stateURL, bytes.NewBuffer([]byte(data)))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 201 {
		log.Println("Status code 201")
	} else if resp.StatusCode == 400 {
		log.Println("Status code 400")
	}
	log.Println("Successfully persisted state.")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
