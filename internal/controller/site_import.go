package controller

import (
	"daily/internal/service"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// ImportFromURL handles importing content from a provided URL
func ImportFromURL(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get URL from query parameters
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "URL parameter is required", http.StatusBadRequest)
		return
	}

	// Parse the URL and extract site information
	siteInfo, err := service.ParseSiteInfoFromURL(url)
	if err != nil {
		http.Error(w, "Failed to parse URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Save the site info to database
	err = service.SaveSiteInfo(siteInfo)
	if err != nil {
		http.Error(w, "Failed to save site info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the site info as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(siteInfo)
}
