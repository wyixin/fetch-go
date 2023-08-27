package main

import (
	"fetch-go/fetch"
	"fetch-go/utils"
	"fmt"
	"sync"

	flag "github.com/spf13/pflag"
)

var printMetadata bool
var validURLs []string

func init() {
	flag.BoolVar(&printMetadata, "metadata", false, "Display Metadata or not")
}

func main() {
	flag.Parse()
	args := flag.Args()
	validURLs = utils.ParseURI(args)

	fmt.Println("metadata?", printMetadata)
	fmt.Println(validURLs)

	// extract links from content
	// fetch and save date from above links
	// rewrite links to suitable for local file system
	// save the rewited content to index.html file
	url := validURLs[0]
	input := &fetch.FetchInput{
		BaseURL: url,
	}

	instance := fetch.Fetch{
		WG: &sync.WaitGroup{},
	}

	output, _ := instance.FetchALL(input)

	// fmt.Println("Fetched content:", output)
	instance.SavePage(output)
	instance.Wait()
}
