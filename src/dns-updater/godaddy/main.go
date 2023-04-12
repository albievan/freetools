package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	apiKey    string
	apiSecret string
)

type DnsRecord struct {
	Domain   string `json:"domain,omitempty"`
	Type     string `json:"type,omitempty"`
	Name     string `json:"name,omitempty"`
	Data     string `json:"data"`
	TTL      int    `json:"ttl,omitempty"`
	Port     int    `json:"port,omitempty"`
	Service  string `json:"service,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Priority int    `json:"priority,omitempty"`
}

func (r *DnsRecord) ToString() (string, error) {
	recordBytes, err := json.Marshal(r)

	if err != nil {
		log.Println("Marshal error")
		log.Printf("[{\"error\":\"%s\"}]", err.Error())
		return fmt.Sprintf("[{\"error\":\"%s\"}]", err.Error()), err
	}
	return string(recordBytes), nil
}

func main() {
	startTime := time.Now()

	var updateInterval time.Duration
	var recordsFile string

	flag.DurationVar(&updateInterval, "interval", time.Minute*5, "Update interval (e.g., -interval 5m)")
	flag.StringVar(&recordsFile, "file", "domains.json", "File with all the DNS records defined(e.g., -file dnsrecords.json)")
	flag.StringVar(&apiKey, "apikey", "", "Godaddy api key (e.g. -apikey xyz)")
	flag.StringVar(&apiSecret, "apisecret", "", "Godaddy api secret (e.g -apisecret abc)")
	flag.Parse()

	if apiKey == "" {
		apiKey = os.Getenv("DDNS_KEY")
	}
	if apiSecret == "" {
		apiSecret = os.Getenv("DDNS_SEC")
	}

	if apiKey == "" || apiSecret == "" {
		log.Panicln("No api key or secret found")
	}

	for {
		ip, err := getPublicIP()
		if err != nil {
			log.Fatalf("Error getting public IP: %v", err)
		} else {
			log.Printf("Dynamic IP: %s", ip)
		}

		records, err := readRecordsFromFile(recordsFile)
		if err != nil {
			log.Fatalf("Error reading records from file: %v", err)
		}

		log.Printf("records: %v", records)
		for _, r := range records {
			if r.Data == "" {
				r.Data = ip
			}
			log.Printf("Domain: %s, %s", r.Domain, r.Data)

			err := updateDNSRecord(r)
			if err != nil {
				log.Printf("Error updating DNS record: %v", err)
			} else {
				log.Printf("DNS record updated successfully")
			}
		}

		ts := time.Since(startTime)
		duration := fmt.Stringer.String(ts)
		durationSec := float64(ts / time.Second)

		log.Printf("[%s] Waiting %v, execution time: %s(%vs), next: %s, \n", time.Now(), updateInterval, duration, durationSec, time.Now().Add(updateInterval))
		time.Sleep(updateInterval)
	}
}

func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func updateDNSRecord(record DnsRecord) error {
	client := &http.Client{}

	s, _ := record.ToString()
	log.Printf("raw dns record: %s", s)

	_domain := record.Domain
	_type := record.Type
	_name := record.Name

	record.Domain = ""
	record.Type = ""
	record.Name = ""

	recordBytes, err := json.Marshal([]DnsRecord{record})
	if err != nil {
		log.Println("Marshal error")
		return err
	}

	log.Printf("marshaled dns record: %v", string(recordBytes))
	log.Printf("https://api.godaddy.com/v1/domains/%s/records/%s/%s", _domain, _type, _name)
	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/%s/%s", _domain, _type, _name),
		strings.NewReader(string(recordBytes)),
	)
	if err != nil {
		log.Println("NewRequest Error")
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", apiKey, apiSecret))

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Execute error (client.do)")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("update failed: %d %s", resp.StatusCode, string(body))
	}

	return nil
}

func readRecordsFromFile(filename string) ([]DnsRecord, error) {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	//log.Printf("file dns record: %s", string(fileBytes))
	var records []DnsRecord
	err = json.Unmarshal(fileBytes, &records)
	if err != nil {
		return nil, err
	}

	return records, nil
}
