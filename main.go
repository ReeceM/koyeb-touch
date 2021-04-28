package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type KoyebApplication struct {
	token   string
	app     string
	service string
}

type Revision struct {
	Revision KoyebService
}

type KoyebService struct {
	organization_id string
	spp_id          string
	service_id      string
	id              string
	version         string
	updated_at      time.Time
	created_at      time.Time
	parent_id       string
	child_id        string
	Definition      KoyebDefinition `json:"definition"`
	statemap        map[string]interface{}
}

type KoyebDefinition struct {
	Name             string                 `json:"name"`
	Routes           []Routes               `json:"routes"`
	Ports            []Ports                `json:"ports"`
	Env              []string               `json:"env"`
	Regions          []string               `json:"regions"`
	Scaling          map[string]interface{} `json:"scaling"`
	Instance_type    string                 `json:"instance_type"`
	Deployment_group string                 `json:"deployment_group"`
	Docker           map[string]interface{} `json:"docker"`
}

type Routes struct {
	Port int    `json:"port"`
	Path string `json:"path"`
}

type Ports struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

func main() {

	var koyeb KoyebApplication

	if len(os.Args) < 2 {
		log.Fatalln("No Arguments provided, you need to run `koyeb-touch $API_TOKEN $APP_NAME $SERVICE_NAME`")
	}

	koyeb.token = os.Args[1]
	koyeb.app = os.Args[2]
	koyeb.service = os.Args[3]

	r := koyeb.latestRevision()

	koyeb.touchRevision(r.Revision)
}

func (k *KoyebApplication) latestRevision() Revision {
	// https://app.koyeb.com/v1/apps/{app_id_or_name}/services/{id_or_name}/revisions/_latest
	if k.app == "" {
		log.Fatalln("Application id required before getting the latest revision")
	}

	if k.service == "" {
		log.Fatalln("The service ID is required to fetch the definition")
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://app.koyeb.com/v1/apps/%s/services/%s/revisions/_latest", k.app, k.service), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", k.token))

	res, err := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	// fmt.Println(string(body))

	var r Revision
	jsonErr := json.Unmarshal(body, &r)

	if jsonErr != nil {
		log.Println("JSON Error")
		log.Fatalln(jsonErr)
	}

	fmt.Println(r.Revision.Definition.Name)

	if err != nil {
		log.Println("Response Error")
		log.Fatalln(err)
	}

	return r
}

func (k *KoyebApplication) touchRevision(service KoyebService) {
	url := fmt.Sprintf("https://app.koyeb.com/v1/apps/%s/services/%s", k.app, service.Definition.Name)

	fmt.Printf("touching %s\n", url)

	var jsonPost []byte
	jsonPost, marshallErr := json.Marshal(service)

	if marshallErr != nil {
		log.Fatalln(marshallErr)
	}

	client := &http.Client{}
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", k.token))
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}
