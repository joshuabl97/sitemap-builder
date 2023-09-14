package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func main() {
	url := "https://www.iana.org/help/example-domains"

	getPageLinks(url)
}

func getPageLinks(url string) {
	// Fetch the HTML content from the URL
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}
	defer resp.Body.Close()

	// Parse the HTML content
	tokenizer := html.NewTokenizer(resp.Body)

	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			// End of the document
			return
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				// Found an <a> tag
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						fmt.Println("Link:", attr.Val)
					}
				}
			}
		}
	}
}
