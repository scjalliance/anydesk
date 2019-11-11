package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/emmaly/anydesk"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	licenseID  = kingpin.Flag("license", "AnyDesk License ID").Required().String()
	apiKey     = kingpin.Flag("apikey", "AnyDesk API Key").Required().String()
	aliasMatch = regexp.MustCompile(`(?i)^scj(\d+)(-\d+)?@ad$`)
)

func main() {
	kingpin.Parse()

	ad, err := anydesk.New(*apiKey, *licenseID, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	_, err = ad.AuthTest()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	clients, err := ad.Clients(&anydesk.ClientsOptions{
		IncludeOffline: true,
	})
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	namespaceCount := make(map[string]int)
	for _, client := range clients.Clients {
		namespace := "-none-"
		nameParts := strings.Split(client.Alias, "@")
		if len(nameParts) > 1 {
			namespace = nameParts[1]
		}
		namespaceCount[namespace]++
	}

	i, err := ad.SysInfo()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("AnyDesk Sessions Online : %d/%d\n", i.Sessions.Online, i.License.MaxSessions)
	fmt.Printf("AnyDesk Clients Online  : %d/%d\n", i.Clients.Online, i.License.MaxClients)
	fmt.Printf("Clients Without Alias   : %d/-1\n", namespaceCount["-none-"])
	fmt.Printf("Namespace            ad : %d/-1\n", namespaceCount["ad"])
	for _, namespace := range i.License.Namespaces {
		fmt.Printf("Namespace %13s : %d/%d\n", namespace.Name, namespaceCount[namespace.Name], namespace.Size)
	}
}
