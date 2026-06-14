package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

type URLShortener []struct {
	Path string `json:"path"`
	Url  string `json:"url"`
}

type FileType string

const (
	Json FileType = "json"
	Yaml FileType = "yaml"
)

func main() {
	filePath := flag.String("f", "./url.yml", "Defines the file path. Default is 'url.yml' - Supported format: JSON, YAML")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	splits := strings.Split(*filePath, ".")
	extension := splits[len(splits)-1]
	var fileType FileType

	switch extension {
	case "yml", "yaml":
		fileType = Yaml
	case "json":
		fileType = Json
	default:
		fmt.Printf("Unsupported file type: %s, terminating!", extension)
		return
	}

	content, err := os.ReadFile(*filePath)
	if err != nil {
		fmt.Printf("Failed to read the file: %s", err)
		return
	}

	var urlShorteners URLShortener

	switch fileType {
	case Yaml:
		if err := yaml.Unmarshal(content, &urlShorteners); err != nil {
			fmt.Printf("Failed to unmarshal YAML: %s", err)
			return
		}
	case Json:
		if err := json.Unmarshal(content, &urlShorteners); err != nil {
			fmt.Printf("Failed to unmarshal JSON: %s", err)
			return
		}
	}

	for _, urlShortener := range urlShorteners {
		mux.Handle("/"+urlShortener.Path, http.RedirectHandler(urlShortener.Url, http.StatusFound))
	}

	fmt.Printf("\nRegistered %d paths", len(urlShorteners))
	fmt.Println("Serving server on :8080")
	http.ListenAndServe(":8080", mux)
}
