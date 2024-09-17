package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	// Parse command-line arguments
	url := flag.String("url", "", "URL to download HTML from")
	dirName := flag.String("name", "", "Name of the directory to save files")
	keepJS := flag.Bool("keep-js", false, "Keep all <script> tags in the saved HTML")
	flag.Parse()

	if *url == "" || *dirName == "" {
		fmt.Println("Usage: go run main.go -url <URL> -name <Directory Name> [--no-js]")
		return
	}

	// Create the directory
	err := os.Mkdir(*dirName, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	// Download the HTML content
	resp, err := http.Get(*url)
	if err != nil {
		fmt.Printf("Error downloading HTML: %v\n", err)
		return
	}
	defer resp.Body.Close()

	htmlFilePath := filepath.Join(*dirName, "index.html")
	htmlFile, err := os.Create(htmlFilePath)
	if err != nil {
		fmt.Printf("Error creating HTML file: %v\n", err)
		return
	}
	defer htmlFile.Close()

	htmlContent, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading HTML content: %v\n", err)
		return
	}

	// Parse the HTML content
	doc, err := html.Parse(strings.NewReader(string(htmlContent)))
	if err != nil {
		fmt.Printf("Error parsing HTML: %v\n", err)
		return
	}

	// Find and download CSS files
	var cssFiles []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "link" {
			for _, attr := range n.Attr {
				if attr.Key == "rel" && attr.Val == "stylesheet" {
					for _, attr := range n.Attr {
						if attr.Key == "href" {
							cssFiles = append(cssFiles, attr.Val)
							break
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	for _, cssFile := range cssFiles {
		cssURL := cssFile
		if !strings.HasPrefix(cssURL, "http") {
			cssURL = *url + cssFile
		}

		resp, err := http.Get(cssURL)
		if err != nil {
			fmt.Printf("Error downloading CSS file %s: %v\n", cssFile, err)
			continue
		}
		defer resp.Body.Close()

		cssFileName := filepath.Base(cssFile)
		cssFilePath := filepath.Join(*dirName, cssFileName)
		cssFile, err := os.Create(cssFilePath)
		if err != nil {
			fmt.Printf("Error creating CSS file %s: %v\n", cssFileName, err)
			continue
		}
		defer cssFile.Close()

		_, err = io.Copy(cssFile, resp.Body)
		if err != nil {
			fmt.Printf("Error saving CSS file %s: %v\n", cssFileName, err)
			continue
		}

		// Update the HTML to reference the locally downloaded CSS file
		updateHTMLReferences(doc, cssFilePath, cssFileName)
	}

	// Optionally keep all <script> tags
    if !*keepJS {
        removeScriptTags(doc)
    }

	// Save the updated HTML
	err = html.Render(htmlFile, doc)
	if err != nil {
		fmt.Printf("Error saving updated HTML: %v\n", err)
		return
	}

	fmt.Println("Download complete.")
}

func updateHTMLReferences(n *html.Node, oldHref, newHref string) {
	if n.Type == html.ElementNode && n.Data == "link" {
		for i, attr := range n.Attr {
			if attr.Key == "href" && attr.Val == oldHref {
				n.Attr[i].Val = newHref
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		updateHTMLReferences(c, oldHref, newHref)
	}
}

func removeScriptTags(n *html.Node) {
    for c := n.FirstChild; c != nil; {
        next := c.NextSibling
        if c.Type == html.ElementNode && c.Data == "script" {
            n.RemoveChild(c)
        } else {
            removeScriptTags(c)
        }
        c = next
    }
}
