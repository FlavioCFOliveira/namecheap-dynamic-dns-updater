package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var cfg Config
var CurrentIpAddress string

func init() {

	// get configuration file path
	configFlag := flag.String("configFlag", "config.json", "a string")
	flag.Parse()

	// Reads the file content
	jsonFile, err := os.Open(*configFlag)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	jsonContent, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the config file content
	json.Unmarshal(jsonContent, &cfg)

	// Get The current public IpAddress
	ip, err := GetPublicIpAddress()
	if err != nil {
		log.Fatalf("Cannot obtain ip address. %s", err)
	}
	CurrentIpAddress = ip
}

func main() {

	if len(cfg.Profiles) > 0 {
		for _, p := range cfg.Profiles {

			log.Printf("Updating profile '%s':", p.ProfileName)
			if len(p.Hosts) > 0 {

				log.Printf(" - %v hosts found for the domain '%s':", len(p.Hosts), p.Domain)
				for _, h := range p.Hosts {
					e := UpdateHost(p.Domain, h, p.Password)
					if e == nil {
						log.Printf("  - Host '%s' updated to %s;", h, CurrentIpAddress)
					} else {
						log.Printf("  - Error updating host '%s': %s;", h, e)
					}
				}

			} else {
				log.Printf(" - no hosts to update")
			}

		}
	}
}

func UpdateHost(d string, h string, p string) error {

	uri := fmt.Sprintf("https://dynamicdns.park-your-domain.com/update?host=%s&domain=%s&password=%s&ip=%s", h, d, p, CurrentIpAddress)

	resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		/*body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.*/
		log.Print(uri)
		return errors.New(fmt.Sprintf("%v - %s", resp.StatusCode, resp.Status))
	}

	return nil
}

func GetPublicIpAddress() (string, error) {

	resp, err := http.Get(cfg.IpAddressProvider)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body), nil
}
