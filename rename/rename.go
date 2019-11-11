package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/emmaly/anydesk"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	licenseID             = kingpin.Flag("license", "AnyDesk License ID").Required().String()
	apiKey                = kingpin.Flag("apikey", "AnyDesk API Key").Required().String()
	aliasMatch            = regexp.MustCompile(`(?i)^scj(\d+)(-\d+)?@ad$`)
	newNamespace          = "@scj"
	longestOldAliasLength = 0
	longestNewAliasLength = 0
	now                   = time.Now()
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
			if oldAliasLength := len(client.Alias); longestOldAliasLength < oldAliasLength {
				longestOldAliasLength = oldAliasLength
			}
			assetID := match[1]
			if client.Online {
				if newAliasLength := len(assetID + newNamespace); longestNewAliasLength < newAliasLength {
					longestNewAliasLength = newAliasLength
				}
				clientMap[assetID] = client
			}
		}
	}

	for assetID, client := range clientMap {
		newAlias := assetID + newNamespace
		// time.Sleep(time.Millisecond * 250) // be a little polite about it...
		// err := ad.ClientAlias(client.ID, newAlias)
		// if err != nil {
		// 	fmt.Printf("Error setting Client ID %d to %s: %s\n", client.ID, newAlias, err.Error())
		// 	continue
		// }
		fmt.Printf("%d = %"+strconv.Itoa(longestOldAliasLength)+"s -> %"+strconv.Itoa(longestNewAliasLength)+"s\n", client.ID, client.Alias, newAlias)
	}
}
