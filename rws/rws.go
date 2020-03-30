package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/limits"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
)

var (
	layoutUS       = "January 2, 2006, 15:04:05"
	topic          = "cpu"
	stateStoreName = "statestore"
	stateURL       = "http://localhost:3500/v1.0/state/" + stateStoreName
)

type data struct {
	MaxCores         int    `json:"maxCores"`
	CurrentCores     int    `json:"currentCores"`
	Treshold         int    `json:"treshold"`
	Email            string `json:"email"`
	ID               string `json:"id"`
	Name             string `json:"name"`
	Date             string `json:"date"`
	MaxRAM           int    `json:"maxRam"`
	CurrentRAM       int    `json:"currentRam"`
	MaxInstances     int    `json:"maxInstances"`
	CurrentInstances int    `json:"currentInstances"`
}

type email struct {
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

type getUserResp struct {
	PortalID    string `json:"portal_id"`
	PortalName  string `json:"portal_name"`
	PortalEmail string `json:"portal_email"`
}

var openstackOpts = gophercloud.AuthOptions{
	IdentityEndpoint: os.Getenv("OPENSTACK_IDENTITY_ENDPOINT"),
	DomainName:       "Default",
	TenantName:       os.Getenv("OPENSTACK_TENANT_NAME"),
	Username:         os.Getenv("OPENSTACK_USERNAME"),
	Password:         os.Getenv("OPENSTACK_PASSWORD"),
}

var projectURL = "http://localhost:3500/v1.0/invoke/gvr-get-client/method/get?name="

func main() {

	for true {
		projects := getAllProjects()
		for _, project := range projects {
			// FOR TROUBLESHOOTING
			if project.ID != "292d78952e584d25b0c71deb2eb06d55" {
				continue
			}
			client := getEmail(project.Name)
			log.Println(client)

			if client.PortalEmail == "" {
				log.Println("There is no e-mail for this Project, skipping")
				continue
			}

			limits := getLimits(project.ID)
			d := data{
				MaxCores:     limits.MaxTotalCores,
				CurrentCores: limits.TotalCoresUsed,
				Treshold:     60,
				Name:         client.PortalName,
				// Email:            client.PortalEmail,
				Email:            "masterhorn89@gmail.com",
				ID:               project.ID,
				Date:             time.Now().Format(layoutUS),
				MaxRAM:           limits.MaxTotalRAMSize,
				CurrentRAM:       limits.TotalRAMUsed,
				MaxInstances:     limits.MaxTotalInstances,
				CurrentInstances: limits.TotalInstancesUsed,
			}

			// FOR TROUBLESHOOTING
			// if d.ID == "292d78952e584d25b0c71deb2eb06d55" {
			// 	d.CurrentCores = 30
			// 	d.CurrentRAM = 40000
			// 	d.CurrentInstances = 8
			// }

			// CPU
			if verifyTreshold(d.MaxCores, d.CurrentCores, d.Treshold) {
				log.Println("CPU treshold reached for: " + d.ID)
				if projectGetState(d.ID, "cpu") == false {
					log.Println("Sending email for: " + d.ID)
					projectSaveState(d.ID, true, "cpu")
					e := email{
						Max:       d.MaxCores,
						Current:   d.CurrentCores,
						Treshold:  60,
						QuotaName: "Cores",
						Name:      d.Name,
						Email:     d.Email,
						ID:        d.ID,
						Date:      time.Now().Format(layoutUS),
					}
					publish(e)
				} else {
					log.Println("Email were already sent for: " + d.ID)
				}
			} else {
				log.Println("CPU treshold not reached for: " + d.ID)
				if projectGetState(d.ID, "cpu") == true {
					log.Println("Reseting indicator of sent email for: " + d.ID)
					projectSaveState(d.ID, false, "cpu")
				}
			}
			// RAM
			if verifyTreshold(d.MaxRAM, d.CurrentRAM, d.Treshold) {
				log.Println("RAM treshold reached for: " + d.ID)
				if projectGetState(d.ID, "RAM") == false {
					log.Println("Sending email for: " + d.ID)
					projectSaveState(d.ID, true, "RAM")
					e := email{
						Max:       d.MaxRAM,
						Current:   d.CurrentRAM,
						Treshold:  60,
						QuotaName: "RAM",
						Name:      d.Name,
						Email:     d.Email,
						ID:        d.ID,
						Date:      time.Now().Format(layoutUS),
					}
					publish(e)
				} else {
					log.Println("Email were already sent for: " + d.ID)
				}
			} else {
				log.Println("RAM treshold not reached for: " + d.ID)
				if projectGetState(d.ID, "RAM") == true {
					log.Println("Reseting indicator of sent email for: " + d.ID)
					projectSaveState(d.ID, false, "RAM")
				}
			}
			// INSTANCES
			if verifyTreshold(d.MaxInstances, d.CurrentInstances, d.Treshold) {
				log.Println("Instances treshold reached for: " + d.ID)
				if projectGetState(d.ID, "Instances") == false {
					log.Println("Sending email for: " + d.ID)
					projectSaveState(d.ID, true, "Instances")
					e := email{
						Max:       d.MaxInstances,
						Current:   d.CurrentInstances,
						Treshold:  60,
						QuotaName: "Instances",
						Name:      d.Name,
						Email:     d.Email,
						ID:        d.ID,
						Date:      time.Now().Format(layoutUS),
					}
					publish(e)
				} else {
					log.Println("Email were already sent for: " + d.ID)
				}
			} else {
				log.Println("Instances treshold not reached for: " + d.ID)
				if projectGetState(d.ID, "Instances") == true {
					log.Println("Reseting indicator of sent email for: " + d.ID)
					projectSaveState(d.ID, false, "Instances")
				}
			}
		}
		log.Println("END")
		time.Sleep(time.Minute * 30)
	}
}

func publish(e email) {
	body, err := json.Marshal(e)
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
}

func verifyTreshold(max int, current int, treshold int) bool {

	x := (max * treshold) / 100

	if current >= x {
		return true
	}
	return false
}

func getLimits(id string) limits.Absolute {
	provider, err := openstack.AuthenticatedClient(openstackOpts)
	if err != nil {
		log.Fatalf("Error with authentication provider from gophercloud/openstack: %s\n", err)
	}

	epOpts := gophercloud.EndpointOpts{
		Region: os.Getenv("OPENSTACK_REGION_NAME"),
	}

	clientCompute, err := openstack.NewComputeV2(provider, epOpts)
	if err != nil {
		log.Fatalf("Error with getting compute: %s\n", err)
	}

	getOpts := limits.GetOpts{
		TenantID: id,
	}

	limits, err := limits.Get(clientCompute, getOpts).Extract()
	if err != nil {
		log.Fatalf("Error with getting quotas: %s\n", err)
	}

	return limits.Absolute
}

func getAllProjects() []projects.Project {
	provider, err := openstack.AuthenticatedClient(openstackOpts)
	if err != nil {
		log.Fatalf("Error with authentication provider from gophercloud/openstack: %s\n", err)
	}

	epOpts := gophercloud.EndpointOpts{
		Region: os.Getenv("OPENSTACK_REGION_NAME"),
	}

	clientProject, err := openstack.NewIdentityV3(provider, epOpts)
	if err != nil {
		log.Fatalf("Error with getting compute: %s\n", err)
	}

	listOpts := projects.ListOpts{
		Enabled: gophercloud.Enabled,
	}

	allpj, err := projects.List(clientProject, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allProjects, err := projects.ExtractProjects(allpj)
	if err != nil {
		panic(err)
	}
	return allProjects
}

func projectGetState(id string, quota string) bool {
	req, err := http.NewRequest("GET", stateURL+"/"+quota+"_"+id, nil)
	if err != nil {
		log.Fatalln("projectGetState - new request: ", err)
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("projectGetState - do request: ", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if string(bodyBytes) == "" {
		return false
	}
	var sV stateValue
	err = json.Unmarshal(bodyBytes, &sV)
	if err != nil {
		log.Fatalln("projectGetState - unmarshal: ", err)
	}
	// log.Println(sV)
	return sV.Mail
}

func projectSaveState(id string, mail bool, quota string) {

	sV := stateValue{
		Mail: mail,
		Date: time.Now().Format(layoutUS),
	}

	s := state{
		Key:   quota + "_" + id,
		Value: sV,
	}

	states := []state{s}

	jsonS, err := json.Marshal(states)
	if err != nil {
		log.Fatalln("projectSaveState - marshal: ", err)
	}

	req, err := http.NewRequest("POST", stateURL, bytes.NewBuffer(jsonS))
	if err != nil {
		log.Fatalln("projectSaveState - new request: ", err)
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("projectSaveState - do request: ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		log.Println("Successfully persisted state.")
	}
}

func getEmail(name string) getUserResp {
	URL := projectURL + name
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatalln("getEmail: Error creating request", err)
	}
	// log.Println(URL)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("getEmail: Error sending http request", err)
	}
	defer resp.Body.Close()

	var g getUserResp

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return g
	}

	// log.Println(string(bodyBytes))
	err = json.Unmarshal(bodyBytes, &g)
	if err != nil {
		log.Fatalln("getEmail: Error with unmarshaling response", err)
	}
	return g
}
