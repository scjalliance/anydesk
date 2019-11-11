package http

import (
	"net/http"
	"os"

	"github.com/scjalliance/anydesk/stats"
)

// HandleHTTP handles HTTP for running in Google Cloud Functions
func HandleHTTP(w http.ResponseWriter, r *http.Request) {
	licenseID, ok := os.LookupEnv("ANYDESK_LICENSE_ID")
	if !ok {
		http.Error(w, "Missing LICENSE in server-side configuration", http.StatusInternalServerError)
		return
	}

	apiKey, ok := os.LookupEnv("ANYDESK_APIKEY")
	if !ok {
		http.Error(w, "Missing APIKEY in server-side configuration", http.StatusInternalServerError)
		return
	}

	var ezKey string
	if r.Method == "POST" {
		ezKey, _ = os.LookupEnv("STATHAT_EZKEY")
	}

	err := stats.Run(licenseID, apiKey, ezKey, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
