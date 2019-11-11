package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/emmaly/anydesk"
	"github.com/gentlemanautomaton/stathat"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	licenseID = kingpin.Flag("license", "AnyDesk License ID").Required().String()
	apiKey    = kingpin.Flag("apikey", "AnyDesk API Key").Required().String()
	ezKey     = kingpin.Flag("ezkey", "StatHat EZKey").Required().String()
)

func intToStr(i int) string {
	if i == -1 {
		return "∞"
	}
	return strconv.Itoa(i)
}

func main() {
	kingpin.Parse()

	stat := stathat.New().EZKey(*ezKey)

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
	stat.PostEZ("AnyDesk Sessions Online", stathat.KindValue, float64(i.Sessions.Online), nil)
	stat.PostEZ("AnyDesk Sessions Max", stathat.KindValue, float64(i.License.MaxSessions), nil)

	fmt.Printf("AnyDesk Clients Online  : %5d / %s\n", i.Clients.Online, intToStr(i.License.MaxClients))
	stat.PostEZ("AnyDesk Clients Online", stathat.KindValue, float64(i.Clients.Online), nil)
	stat.PostEZ("AnyDesk Clients Max", stathat.KindValue, float64(i.License.MaxClients), nil)

	fmt.Printf("Clients Without Alias   : %5d / %s\n", namespaceCount["@"], intToStr(-1))
	stat.PostEZ("AnyDesk Clients Without Alias", stathat.KindValue, float64(namespaceCount["@"]), nil)

	fmt.Printf("Namespace            ad : %5d / %s\n", namespaceCount["ad"], intToStr(-1))
	stat.PostEZ("AnyDesk Clients @ad", stathat.KindValue, float64(namespaceCount["ad"]), nil)

	for _, namespace := range i.License.Namespaces {
		fmt.Printf("Namespace %13s : %5d / %s\n", namespace.Name, namespaceCount[namespace.Name], intToStr(namespace.Size))
		stat.PostEZ("AnyDesk Clients @"+namespace.Name, stathat.KindValue, float64(namespaceCount[namespace.Name]), nil)
	}
}
