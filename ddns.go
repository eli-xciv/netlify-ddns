package netlify-ddns

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var netlifyUrl string = "https://api.netlify.com/api/v1/dns_zones/%s/dns_records"
var ipUrl string = "https://api.ipify.org/"
var bearerToken string = "Bearer %s"
var netlifyApiToken string = ""
var netlifyDnsUrl string = ""
var netlifySanitizedDnsUrl string = ""
var previousId string = ""
var previousIp string = ""

type Get_Dns_Response []struct {
	Hostname    string `json:"hostname"`
	Id          string `json:"id"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Ttl         int    `json:"ttl"`
	Priority    int    `json:"priority"`
	Dns_zone_id string `json:"dns_zone_id"`
	Site_id     string `json:"site_id"`
	Flag        int    `json:"flag"`
	Tag         string `json:"tag"`
	Managed     bool   `json:"managed"`
}

type Add_Dns_Request struct {
	Hostname string `json:"hostname"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	Ttl      int    `json:"ttl"`
}

func getCurrentNetlifyDns() (string, string) {
	url := fmt.Sprintf(netlifyUrl, netlifySanitizedDnsUrl)
	bearer := fmt.Sprintf(bearerToken, netlifyApiToken)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var netlifyResponse Get_Dns_Response

	if err := json.Unmarshal(body, &netlifyResponse); err != nil {
		log.Fatalln(err)
	}
	return netlifyResponse[0].Id, netlifyResponse[0].Value
}

func updateNetlifyDns() {
	log.Println("Checking DDNS status")
	oldId, oldIp := getCurrentNetlifyDns()

	req, err := http.NewRequest("GET", ipUrl, nil)
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if string(body) == oldIp {
		log.Printf("IP's are the same, doing nothing")
	} else {
		log.Println("DDNS IP's have changed")
		log.Println(oldIp, " has changed to ", string(body), " ; Updating netlify with current IP")

		var addRequest Add_Dns_Request

		addRequest.Hostname = netlifyDnsUrl
		addRequest.Value = string(body)
		addRequest.Ttl = 3600
		addRequest.Type = "A"
		// add new record to netlify

		// delete old record
		log.Println("removing old record of id:", oldId)
		//req, err = http.NewRequest("Delete",  ,nil)
	}

}

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

func main() {
	fmt.Println("Netlify DNS Client")
	fmt.Println("Written By: Eli-XCIV (https://github.com/eli-xciv)")

	netlifyDnsUrl = os.Getenv("NETLIFY_URL")
	netlifySanitizedDnsUrl = replaceAtIndex(netlifyDnsUrl, '_', strings.LastIndex(netlifyDnsUrl, "."))
	netlifyApiToken = os.Getenv("NETLIFY_API_TOKEN")

	log.Println("Starting Service")
	for true {
		updateNetlifyDns()
		time.Sleep(5 * time.Minute)
	}

}
