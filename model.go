package main

type (
	Config struct {
		IpAddressProvider string
		Profiles          []Profile
	}
	Profile struct {
		ProfileName string
		Domain      string
		Password    string
		Hosts       []string
	}
)
