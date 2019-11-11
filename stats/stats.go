package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/emmaly/anydesk"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	licenseID  = kingpin.Flag("license", "AnyDesk License ID").Required().String()
	apiKey     = kingpin.Flag("apikey", "AnyDesk API Key").Required().String()
	aliasMatch = regexp.MustCompile(`(?i)^scj(\d+)(-\d+)?@ad$`)
)

func intToStr(i int) string {
	if i == -1 {
		return "∞"
	}
	return strconv.Itoa(i)
}

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
		namespace := "@"
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

	fmt.Printf("AnyDesk Sessions Online : %5d / %s\n", i.Sessions.Online, intToStr(i.License.MaxSessions))
	fmt.Printf("AnyDesk Clients Online  : %5d / %s\n", i.Clients.Online, intToStr(i.License.MaxClients))
	fmt.Printf("Clients Without Alias   : %5d / %s\n", namespaceCount["@"], intToStr(-1))
	fmt.Printf("Namespace            ad : %5d / %s\n", namespaceCount["ad"], intToStr(-1))
	for _, namespace := range i.License.Namespaces {
		fmt.Printf("Namespace %13s : %5d / %s\n", namespace.Name, namespaceCount[namespace.Name], intToStr(namespace.Size))
	}
}
