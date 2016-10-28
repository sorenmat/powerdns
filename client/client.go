package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"strconv"
)

type CreateZone struct {
	Name        string   `json:"name"`
	Kind        string   `json:"kind"`
	Masters     []string `json:"masters"`
	Nameservers []string `json:"nameservers"`
}

type RecordSet struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	TTL        int      `json:"ttl"`
	Changetype string   `json:"changetype"`
	Records    []Record `json:"records"`
}
type Record struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	Disabled bool   `json:"disabled"`
}
type CreateRecord struct {
	Rrsets []RecordSet `json:"rrsets"`
}

type GetZone struct {
	Account        string        `json:"account"`
	Dnssec         bool          `json:"dnssec"`
	ID             string        `json:"id"`
	Kind           string        `json:"kind"`
	LastCheck      int           `json:"last_check"`
	Masters        []interface{} `json:"masters"`
	Name           string        `json:"name"`
	NotifiedSerial int           `json:"notified_serial"`
	Rrsets         []struct {
		Comments []interface{} `json:"comments"`
		Name     string        `json:"name"`
		Records  []struct {
			Content  string `json:"content"`
			Disabled bool   `json:"disabled"`
		} `json:"records"`
		TTL      int    `json:"ttl"`
		Type     string `json:"type"`
	} `json:"rrsets"`
	Serial         int    `json:"serial"`
	SoaEdit        string `json:"soa_edit"`
	SoaEditAPI     string `json:"soa_edit_api"`
	URL            string `json:"url"`
}

type PowerClient struct {
	// baseURL is the url for the powerdns host like http://localhost:8081
	baseURL string
	apiKey  string
}

func NewClient(baseURL string, apiKey string) *PowerClient {
	return &PowerClient{apiKey: apiKey, baseURL: baseURL}
}

func (c *PowerClient) GetZone(name string) (*GetZone, error) {
	url := c.baseURL + "/api/v1/servers/localhost/zones" + "/" + name
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("X-API-Key", c.apiKey)
	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error calling PowerDNS resource %v", err))
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Error calling PowerDNS %v got statuscode '%v' with error %v", url, resp.StatusCode, err))
	}
	jd := json.NewDecoder(resp.Body)
	zonedata := &GetZone{}
	err = jd.Decode(zonedata)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error decoding JSON from PowerDNS resource %v", err))
	}
	return zonedata, nil
}

func (c *PowerClient) AddZone(name string, nameServers []string) error {
	cl := CreateZone{
		Name:        name,
		Kind:        "Native",
		Masters:     []string{},
		Nameservers: nameServers,
	}
	b, err := json.Marshal(cl)
	if err != nil {
		return errors.New("failure parsing zone struct to json")
	}
	url := c.baseURL + "/api/v1/servers/localhost/zones"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Add("X-API-Key", c.apiKey)
	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		status := ""
		if resp != nil {
			status = resp.Status
		}
		return errors.New(fmt.Sprintf("HTTP call to %v '%v' returned %v", req.Method, url, status))

	}
	if resp.StatusCode != 201 {
		rb, _ := ioutil.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("HTTP call returned %v with %v\n\t%v", resp.StatusCode, resp.Status, string(rb)))

	}
	return nil

}

func (c *PowerClient) AddSOARecord(name, primaryDNS, admin string, refreshSeconds, failedRefresh, authoritativeTimeout, negativeTTL int, zone string) error {
	/*
	The SOA record includes the following details:

	* The primary name server for the domain, which is ns1.example.com or the first name server in the vanity name server list for vanity name servers.
	* The responsible party for the domain, which is admin.example.com = admin@example.com.
	* A timestamp that changes when you update your domain name.
	* The number of seconds before the zone should be refreshed.
	* The number of seconds before a failed refresh should be retried.
	* The limit in seconds before a zone is considered no longer authoritative.
	* The negative result TTL.
*/
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	content := fmt.Sprintf("%v %v %v %v %v %v %v", primaryDNS, admin, timestamp, refreshSeconds, failedRefresh, authoritativeTimeout, negativeTTL)
	return c.AddRecord(name, "SOA", content, 30, zone)
}
func (c *PowerClient) AddRecord(name, dnstype, content string, ttl int, zone string) error {

	p := CreateRecord{
		Rrsets: []RecordSet{
			{
				Name: name,
				Type: dnstype,
				TTL: ttl,
				Changetype: "REPLACE",
				Records: []Record{
					{
						Name: name,
						Content:  content,
						Disabled: false,
					},
				}},
		},
	}

	b, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
		return errors.New("failure parsing record struct to json")
	}

	url := c.baseURL + "/api/v1/servers/localhost/zones/" + zone

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(b))
	if err != nil {
		log.Fatal("Error creating request", err)
		return err
	}
	req.Header.Add("X-API-Key", c.apiKey)
	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("HTTP call returned %v with content %v", err, string(body)))

	}

	if resp.StatusCode != 204 && resp.StatusCode != 200 {
		// 204 No content = create, 200 = not updated but otherwise ok
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("HTTP call returned %v\nPowerDNS response: %v", resp.StatusCode, string(body)))

	}
	return nil

}
