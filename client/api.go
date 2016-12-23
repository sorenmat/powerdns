package client

type PowerClient interface {
	GetZone(name string) (*GetZone, error)
	AddZone(name string, nameServers []string) (error, int)
	AddSOARecord(name, primaryDNS, admin string, refreshSeconds, failedRefresh, authoritativeTimeout, negativeTTL int, zone string) (error, int)
	AddSRVRecord(service, proto, name string, ttl int, priority int, weight int, port, target, zone string) (error, int)
	AddRecord(name, dnstype, content string, ttl int, zone string) (error, int)
}
