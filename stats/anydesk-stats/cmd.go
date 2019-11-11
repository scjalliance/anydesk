package main

import (
	"fmt"
	"os"

	"github.com/scjalliance/anydesk/stats"
	"gopkg.in/alecthomas/kingpin.v2"
)

// main runs from the command line
func main() {
	var (
		licenseID = kingpin.Flag("license", "AnyDesk License ID").Envar("ANYDESK_LICENSE_ID").String()
		apiKey    = kingpin.Flag("apikey", "AnyDesk API Key").Envar("ANYDESK_APIKEY").String()
		ezKey     = kingpin.Flag("ezkey", "StatHat EZKey").Envar("STATHAT_EZKEY").String()
	)

	kingpin.Parse()

	if *licenseID == "" || *apiKey == "" {
		fmt.Println("Error: missing license or apikey")
		kingpin.Usage()
		os.Exit(1)
	}

	err := stats.Run(*licenseID, *apiKey, *ezKey, os.Stdout)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
