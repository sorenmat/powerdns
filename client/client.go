package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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
	Content  string `json:"content"`
	Disabled bool   `json:"disabled"`
}
type CreateRecord struct {
	Rrsets []RecordSet `json:"rrsets"`
}

type pdnsclient struct {
	// baseURL is the url for the powerdns host like http://localhost:8081
	baseURL string
	apiKey  string
}

func NewClient(baseURL string, apiKey string) *pdnsclient {
	return &pdnsclient{apiKey: apiKey, baseURL: baseURL}
}

func (c *pdnsclient) AddZone(name string, nameServers []string) error {
	cl := CreateZone{
		Name:        name,
		Kind:        "Native",
		Masters:     []string{},
		Nameservers: nameServers,
	}
	b, err := json.Marshal(cl)
	if err != nil {
		log.Println(err)
		return errors.New("failure parsing zone struct to json")
	}

	url := c.baseURL + "/api/v1/servers/localhost/zones"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
		return err
	}
	req.Header.Add("X-API-Key", c.apiKey)
	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("HTTP call returned %v", resp.Status))

	}

	if resp.StatusCode != 201 {
		return errors.New(fmt.Sprintf("HTTP call returned %v", resp.StatusCode))

	}
	return nil

}

func (c *pdnsclient) AddRecord(name, dnstype, content string, ttl int, zone string) error {

	p := CreateRecord{
		Rrsets: []RecordSet{
			RecordSet{Name: name, Type: dnstype, TTL: ttl, Changetype: "REPLACE",
				Records: []Record{
					Record{
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
		return errors.New(fmt.Sprintf("HTTP call returned %v", err))

	}

	if resp.StatusCode != 204 && resp.StatusCode != 200 { // 204 No content = create, 200 = not updated but otherwise ok
		return errors.New(fmt.Sprintf("HTTP call returned %v", resp.StatusCode))

	}
	return nil

}
