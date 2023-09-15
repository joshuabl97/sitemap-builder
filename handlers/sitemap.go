package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
)

type RequestData struct {
	BaseURL string `json:"url" validate:"required"`
	Depth   int    `json:"depth" validate:"required"`
}

type SiteMap struct {
	BaseURL       string
	Endpoint      string
	CurrentDepth  int
	TotalDepth    int
	Children      []*SiteMap
	Parents       []*SiteMap
	ExternalLinks []string
}

var TraveresedEndpoints map[string]SiteMap

func (h *HandlerHelper) SitemapGen(w http.ResponseWriter, r *http.Request) {
	// decode the request JSON into the RequestData struct
	var requestData RequestData
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		h.l.Error().Msg("Invalid JSON request")
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	// validate the request body
	validate := validator.New()
	if err := validate.Struct(requestData); err != nil {
		h.l.Error().Msg("Invalid JSON request")
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	// first call is using the base url
	// the base_url is reused to construct subsequent calls
	// i.e https://google.com + /path can be used for GET /getLinks {"url": "https://google.com/path"}
	base_url := requestData.BaseURL
	reqData := Link{URL: base_url}
	reqJSON, err := json.Marshal(reqData)
	if err != nil {
		h.l.Error().Msg("Error encoding JSON")
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	// create a new request to GET /getlinks
	getLinksRequest, err := http.NewRequest("GET", "http://localhost:"+h.port+"/getLinks", bytes.NewBuffer(reqJSON))
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// set the request header
	getLinksRequest.Header.Set("Content-Type", "application/json")

	// perform the initial request
	resp, err := http.DefaultClient.Do(getLinksRequest)
	if err != nil {
		http.Error(w, "Error making request to /getLinks", http.StatusInternalServerError)
		return
	}

	// decode the response JSON
	var links []Link
	if err := json.NewDecoder(resp.Body).Decode(&links); err != nil {
		http.Error(w, "Error decoding JSON response", http.StatusInternalServerError)
		return
	}

	// instantiate the root of the sitemap
	sitemap := &SiteMap{
		BaseURL:      base_url,
		Endpoint:     "/",
		CurrentDepth: 1,
		TotalDepth:   requestData.Depth,
	}

	// add the basepath to traveresed endpoints
	TraveresedEndpoints["/"] = *sitemap

	// queue for endpoints that have not yet been searched
	var queue []string

	// find and attach external links to sitemap and build queue
	for _, link := range links {
		if !strings.HasPrefix(link.URL, "/") {
			sitemap.ExternalLinks = append(sitemap.ExternalLinks, link.URL)
			continue
		}
		queue = append(queue, link.URL)
	}

	fmt.Printf("%v\n", sitemap)
	fmt.Printf("%v\n", queue)
}
