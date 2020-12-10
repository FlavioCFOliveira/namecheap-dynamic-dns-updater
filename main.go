package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	cfg              Config
	CurrentIpAddress string
	logFile          *os.File
	SuccessLogger    *log.Logger
	ErrorLogger      *log.Logger
)

func init() {

	// get configuration file path
	configFlag := flag.String("config", "config.json", "a string")
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

	// prepare for file logging
	if cfg.LogToFiles {
		if len(cfg.LogDirectory) == 0 {
			cfg.LogDirectory, _ = os.Getwd()
		}
		filename := time.Now().Format("20060102.log")

		logFile, err := os.OpenFile(fmt.Sprintf("%s/%s", cfg.LogDirectory, filename), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		wrt := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(wrt)
		log.SetPrefix("[INFO] ")
	}

	SuccessLogger = log.New(log.Writer(), "[SUCCESS] ", log.LstdFlags)
	ErrorLogger = log.New(log.Writer(), "[ERROR] ", log.LstdFlags)
}

func main() {
	defer logFile.Close()

	if len(cfg.Profiles) > 0 {
		for _, p := range cfg.Profiles {

			log.Printf("Updating profile '%s':", p.ProfileName)
			if len(p.Hosts) > 0 {

				log.Printf(" - %v hosts found for the domain '%s':", len(p.Hosts), p.Domain)
				for _, h := range p.Hosts {
					e := UpdateHost(p.Domain, h, p.Password)
					if e == nil {
						SuccessLogger.Printf("  - Host '%s' updated to %s;", h, CurrentIpAddress)
					} else {
						ErrorLogger.Printf("  - Error updating host '%s': %s;", h, e)
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
		log.Print(uri)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
		}
		log.Print(body)

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
