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

// This struct holds the needed details about the application and API
// The token is the API token
// While app and service are the string name or id of the target Koyeb App
type KoyebApplication struct {
	token   string
	app     string
	service string
}

// This is the main type returned from the Revisions Api
// the Revision is a service
type Revision struct {
	Revision KoyebService
}

// This is the service, all lowercase values aren't needed in the resulted JSON
// The other details are still there if we need them
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

// The Service Definition is the important part
// it will hold the details that are returned to the actual endpoint
// All values are required by the endpoint
// The json mapping makes it lowercase or the api will fail as they don't match
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

	// Create the new Koyeb instance details
	var koyeb KoyebApplication

	if len(os.Args) < 2 {
		log.Fatalln("No Arguments provided, you need to run `koyeb-touch $API_TOKEN $APP_NAME $SERVICE_NAME`")
	}

	// Load the values from the command line
	// @TODO move to cobra or use flags with an env file option/string to parse
	koyeb.token = os.Args[1]
	koyeb.app = os.Args[2]
	koyeb.service = os.Args[3]

	// fetch the latest revision from the KoyebApplication type. this will fetch for a single service
	r := koyeb.latestRevision()

	// 'Touch' the service to trigger a redploy and build of the service off the latest docker image
	koyeb.touchRevision(r.Revision)
}

// --------------- latest Revision ---------
// This function fetches the latest revision
// does not need to accept any params, it knows what it needs to know from the type
// it will call the api for the latest revision of the selected service
// most of the time the default is `main`.
// --------------------------------
// the reason for doing this is to get the definition of the last deploy
// we only want to reuse that information as it worked likely, but also allows usage
// of any new values set by the user in the dashboard.
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

// ---------------- touchRevision ------------
// this calls the PUT HTTP method on the service that has been retrieved
// N.B. at the moment, the service name is now retrived from the definition, not the value passed in the args
// this will use the definition that has been marshalled from the type and send it back to Koyeb
// it will print the result
// @todo maybe I can make the output be nice, like xyz definition is now active using this id....
//         it will look nice than the dump of JSON in the end
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

	if res.StatusCode == 200 {
		fmt.Println("Completed Successfully")
	}
	// body, _ := ioutil.ReadAll(res.Body)

	fmt.Println("Result is not determinate, please check Koyeb dashboard")
}
