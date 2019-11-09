package main

import (
	"fmt"
	"os"
	"regexp"

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

	var clientMap = make(map[string]anydesk.Client)

	for _, client := range clients.Clients {
		match := aliasMatch.FindStringSubmatch(client.Alias)
		if len(match) > 1 {
			assetID := match[1]
			if existingClient, ok := clientMap[assetID]; !ok || existingClient.OnlineTime < client.OnlineTime {
				clientMap[assetID] = client
			}
		}
	}

	for assetID, client := range clientMap {
		alias := assetID + "@scj"
		err := ad.ClientAlias(client.ID, alias)
		if err != nil {
			fmt.Printf("Error setting Client ID %d to %s: %s\n", client.ID, alias, err.Error())
			continue
		}
		fmt.Printf("%d = %s\n", client.ID, alias)
	}
}
