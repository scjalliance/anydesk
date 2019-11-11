package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocarina/gocsv"

	"github.com/emmaly/anydesk"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	licenseID             = kingpin.Flag("license", "AnyDesk License ID").Required().String()
	apiKey                = kingpin.Flag("apikey", "AnyDesk API Key").Required().String()
	zohoKey               = kingpin.Flag("zohokey", "Zoho Creator API Key").Required().String() // https://accounts.zoho.com/apiauthtoken/create?SCOPE=ZohoCreator/creatorapi
	zohoView              = "AnyDesk_Computers"
	aliasMatch            = regexp.MustCompile(`(?i)^scj(\d+)(-\d+)?@ad$`)
	newNamespace          = "@scj"
	longestOldAliasLength = 0
	longestNewAliasLength = 0
	now                   = time.Now()
)

// ZohoComputer is a computer in Zoho
type ZohoComputer struct {
	AssetID                 string `csv:"Asset ID"`
	WindowsName             string `csv:"Windows Name"`
	PrimaryUser             string `csv:"Primary User"`
	PrimaryUserEmailAddress string `csv:"Primary Email Address"`
}

func main() {
	kingpin.Parse()

	zohoResp, err := http.Get("https://creator.zoho.com/api/csv/it/view/" + zohoView + "?authtoken=" + (*zohoKey) + "&scope=creatorapi")
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	defer zohoResp.Body.Close()
	zohoComputers := make([]ZohoComputer, 0)
	err = gocsv.Unmarshal(zohoResp.Body, &zohoComputers)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	computers := make(map[string]ZohoComputer)
	for _, zohoComputer := range zohoComputers {
		computers[zohoComputer.AssetID] = zohoComputer
	}

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
		var primaryUser, primaryUserName, primaryUserEmail string
		if computer, ok := computers[assetID]; ok {
			primaryUserName = computer.PrimaryUser
			primaryUserEmail = computer.PrimaryUserEmailAddress
			if primaryUserEmail != "" {
				primaryUser = strings.TrimSpace(fmt.Sprintf("%s <%s>", primaryUserName, primaryUserEmail))
			}
		}
		fmt.Printf("%d = %"+strconv.Itoa(longestOldAliasLength)+"s -> %"+strconv.Itoa(longestNewAliasLength)+"s\t%s\n", client.ID, client.Alias, newAlias, primaryUser)
	}
}
