package dom

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parseImages(doc *goquery.Document) []string {
	return getAttr(doc, "img", "src")
}

// TODO
func parseJS(doc *goquery.Document) []string { return nil }

func filterCSSLinks(links []string) []string {
	cssLinks := []string{}
	cssRegex := regexp.MustCompile(`\.css\?`)

	for _, link := range links {
		if cssRegex.MatchString(link) {
			cssLinks = append(cssLinks, link)
		}
	}

	return cssLinks
}

func parseCSS(doc *goquery.Document) []string {
	files := getAttr(doc, "link", "href")

	return filterCSSLinks(files)
}

func getAttr(doc *goquery.Document, tag string, attr string) []string {
	files := []string{}
	doc.Find(tag).Each(func(i int, s *goquery.Selection) {
		value, exists := s.Attr(attr)
		if exists {
			files = append(files, value)
		}
	})
	return files
}

func ParseAllAssets(content string) (images, js, css []string, err error) {

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))

	if err != nil {
		return nil, nil, nil, err
	}

	images = parseImages(doc)
	js = parseJS(doc)
	css = parseCSS(doc)

	return
}
