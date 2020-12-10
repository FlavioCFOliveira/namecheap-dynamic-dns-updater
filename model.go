package main

type (
	Config struct {
		IpAddressProvider string
		Profiles          []Profile
		LogToFiles        bool
		LogDirectory      string
	}
	Profile struct {
		ProfileName string
		Domain      string
		Password    string
		Hosts       []string
	}
)
