package client

import (
	"fmt"
	"time"
	"strconv"
	"strings"
)

type MockClient struct {
	PowerClient
	zonedata   []CreateZone
	recorddata []CreateRecord
}

func (m *MockClient) GetZone(name string) (*GetZone, error) {
	for _, r := range m.zonedata {
		if r.Name == name {
			return &GetZone{
				Name: r.Name,
			}, nil

		}
	}
	return nil, nil
}
func (m *MockClient) AddZone(name string, nameServers []string) error {
	cl := CreateZone{
		Name:        name,
		Kind:        "Native",
		Masters:     []string{},
		Nameservers: nameServers,
	}
	m.zonedata = append(m.zonedata, cl)
	return nil
}
func (m *MockClient) AddSOARecord(name, primaryDNS, admin string, refreshSeconds, failedRefresh, authoritativeTimeout, negativeTTL int, zone string) error {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	content := fmt.Sprintf("%v %v %v %v %v %v %v", primaryDNS, admin, timestamp, refreshSeconds, failedRefresh, authoritativeTimeout, negativeTTL)
	return m.AddRecord(name, "SOA", content, 30, zone)

}
func (m *MockClient) AddSRVRecord(service, proto, name string, ttl int, priority int, weight int, port, target, zone string) error {
	if !strings.HasSuffix(name, ".") {
		name = name + "."
	}
	content := fmt.Sprintf("_%v._%v.%v %v IN %v %v", service, proto, name, ttl, weight, port, target)
	return m.AddRecord(name, "SRV", content, ttl, zone)
}

func (m *MockClient) AddRecord(name, dnstype, content string, ttl int, zone string) error {
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
	m.recorddata = append(m.recorddata, p)
	return nil
}
