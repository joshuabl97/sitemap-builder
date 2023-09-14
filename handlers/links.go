package handlers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/html"
)

type Link struct {
	URL string `json:"url"`
}

func (h *HandlerHelper) GetPageLinks(w http.ResponseWriter, r *http.Request) {
	// decode the request JSON into the RequestData struct
	var requestData Link
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		h.l.Error().Msg("Invalid JSON request")
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	// fetch the HTML content from the URL
	resp, err := http.Get(requestData.URL)
	if err != nil {
		h.l.Error().Msg("Error fetching URL")
		http.Error(w, "Error fetching URL", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// create a slice to store the links
	links := []Link{}

	// parse the HTML content
	tokenizer := html.NewTokenizer(resp.Body)

	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			// end of the document
			// marshal the links slice to JSON
			jsonData, err := json.Marshal(links)
			if err != nil {
				h.l.Error().Msg("Error encoding JSON")
				http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
				return
			}
			// set the content type header to JSON
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
			return
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				// found an <a> tag
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						// append the link to the slice
						links = append(links, Link{URL: attr.Val})
					}
				}
			}
		}
	}
}
