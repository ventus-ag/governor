package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var token = os.Getenv("PORTAL_TOKEN")
var basePortalURL = os.Getenv("PORTAL_URL")
var ren = render.New()

type getUserResp struct {
	PortalID    string `json:"portal_id"`
	PortalName  string `json:"portal_name"`
	PortalEmail string `json:"portal_email"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/get", get).Queries("name", "{name}").Methods("GET")

	log.Printf("server started")
	log.Fatal(http.ListenAndServe(":3000", r))
}

func get(w http.ResponseWriter, r *http.Request) {
	projectName := r.FormValue("name")
	// fields := strings.Fields(testProjectName)
	fields := strings.Fields(projectName)
	id := fields[len(fields)-1]
	var name string
	for i := 0; i < len(fields)-2; i++ {
		name = name + fields[i] + " "
	}
	name = strings.TrimSuffix(name, " ")
	log.Println(id)
	log.Println(name)
	resp := getClient(id)

	w.Header().Set("Content-Type", "application/json")
	ren.JSON(w, http.StatusOK, resp)
}

func getClient(id string) getUserResp {
	URL := basePortalURL + "clients/" + id
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	// log.Println(string(bodyBytes))

	// work with Portal client
	var c portalClient
	err = json.Unmarshal(bodyBytes, &c)
	if err != nil {
		log.Fatalln(err)
	}
	portalResp := getUserResp{
		PortalID:    strconv.Itoa(c.ID),
		PortalName:  c.Name,
		PortalEmail: c.Email,
	}
	return portalResp
}
