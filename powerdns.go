//go:generate stringer -type=
package powerdns

import (
	"io"
	"net"
	"net/http"
	"net/url"
)

// PowerDNSClient provides a wrapper around common API requests and some helpers around net/http.Client
type PowerDNSClient struct {
	Server string
	ApiKey string
}

type URL url.URL

type AuthoritativeZoneKind int
type RecursiveZoneKind int

const (
	Native AuthoritativeZoneKind = iota
	Master
	Slave
)

const (
	RecNative RecursiveZoneKind = iota
	RecForwarded
)

type Server struct {
	Type       string `json:"type"`
	Id         string `json:"id"`
	URL        string `json:"url"`
	DaemonType string `json:"daemon_type"`
	Version    string `json:"version"`
	ConfigURL  string `json:"config_url"`
	ZonesURL   string `json:"zones_url"`
}

type Config struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Zone struct {
	Id               string   `json:"id"`
	Name             string   `json:"name"`
	Type             string   `json:"type"`
	URL              string   `json:"url"`
	Kind             string   `json:"kind"`
	Serial           int      `json:"serial"`
	NotifiedSerial   int      `json:"notified_serial"`
	Masters          []net.IP `json:"masters"`
	DNSSEC           bool     `json:"dnssec"`
	NSEC3Param       bool     `json:"nsec3param"`  // UNSUPPORTED UPSTREAM
	NSEC3Narrow      bool     `json:"nsec3narrow"` // UNSUPPORTED UPSTREAM
	Presigned        bool     `json:"presigned"`
	SOAEdit          string   `json:"soa_edit"`
	SOAEditAPI       string   `json:"soa_edit_api"`
	Account          string   `json:"account"`
	Nameservers      []string `json:"nameservers"`
	Servers          []string `json:"servers"`
	RecursionDesired bool     `json:"recursion_desired"`
	RRSets           []RRSet  `json:"rrsets"`
}

type RRSet struct {
	Name     string
	Type     string // FIXME: enumerate them all and make an enum
	TTL      int
	Records  []Record
	Comments []Comment
}

type Record struct {
	Content  string
	Disabled bool
}

type Comment struct {
	Content    string
	Account    string
	ModifiedAt int
}

func (p PowerDNSClient) prepareRequest(method, url string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, p.Server+url, body)
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"X-API-Key": []string{p.ApiKey},
		"Accept":    []string{"text/json"},
	}
	return req, err
}
