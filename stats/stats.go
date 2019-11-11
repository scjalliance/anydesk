package stats

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/emmaly/anydesk"
	"github.com/gentlemanautomaton/stathat"
)

func intToStr(i int) string {
	if i == -1 {
		return "∞"
	}
	return strconv.Itoa(i)
}

// Run the stats
func Run(licenseID string, apiKey string, ezKey string, out io.Writer) error {
	stat := stathat.New()
	if ezKey != "" {
		stat = stat.EZKey(ezKey)
	} else {
		stat = stat.Noop()
	}

	ad, err := anydesk.New(apiKey, licenseID, nil)
	if err != nil {
		return err
	}

	_, err = ad.AuthTest()
	if err != nil {
		return err
	}

	now := time.Now()
	clients, err := ad.Clients(&anydesk.ClientsOptions{
		IncludeOffline: true,
	})
	if err != nil {
		return err
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
		fmt.Fprintf(out, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(out, "AnyDesk Sessions Online : %5d / %s\n", i.Sessions.Online, intToStr(i.License.MaxSessions))
	stat.PostEZ("AnyDesk Sessions Online", stathat.KindValue, float64(i.Sessions.Online), &now)
	stat.PostEZ("AnyDesk Sessions Max", stathat.KindValue, float64(i.License.MaxSessions), &now)

	fmt.Fprintf(out, "AnyDesk Clients Online  : %5d / %s\n", i.Clients.Online, intToStr(i.License.MaxClients))
	stat.PostEZ("AnyDesk Clients Online", stathat.KindValue, float64(i.Clients.Online), &now)
	stat.PostEZ("AnyDesk Clients Max", stathat.KindValue, float64(i.License.MaxClients), &now)

	fmt.Fprintf(out, "Clients Without Alias   : %5d / %s\n", namespaceCount["@"], intToStr(-1))
	stat.PostEZ("AnyDesk Clients Without Alias", stathat.KindValue, float64(namespaceCount["@"]), &now)

	fmt.Fprintf(out, "Namespace            ad : %5d / %s\n", namespaceCount["ad"], intToStr(-1))
	stat.PostEZ("AnyDesk Clients @ad", stathat.KindValue, float64(namespaceCount["ad"]), &now)

	for _, namespace := range i.License.Namespaces {
		fmt.Fprintf(out, "Namespace %13s : %5d / %s\n", namespace.Name, namespaceCount[namespace.Name], intToStr(namespace.Size))
		stat.PostEZ("AnyDesk Clients @"+namespace.Name, stathat.KindValue, float64(namespaceCount[namespace.Name]), &now)
	}

	return nil
}
