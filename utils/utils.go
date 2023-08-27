package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// get valid urls from args and return as slice
func ParseURI(args []string) []string {

	validURLs := []string{}

	for _, arg := range args {
		u, err := url.Parse(arg)
		if err == nil && u.Scheme != "" && u.Host != "" {
			validURLs = append(validURLs, u.String())
		}
	}

	return validURLs
}

// http://baidu.com/a.html or https://autify.com/
// pages/hash(url)/index.html
//
//	------------------xxhshshhs.css
func GetDirName(urlString string) string {
	return "pages/" + HashContent(urlString)
}

func GetFileName() string {
	return "index.html"
}

func HashContent(content string) string {
	hasher := sha256.New()
	hasher.Write([]byte(content))
	hashBytes := hasher.Sum(nil)
	// todo
	return hex.EncodeToString(hashBytes)[:8] // Using first 8 characters of hash as name
}

func SaveFile(dirName, fileName, content string) error {
	filePath := filepath.Join(dirName, fileName)
	err := ioutil.WriteFile(filePath, []byte(content), os.ModePerm)
	return err
}

// https://autify.com/wp-content/uploads/2022/02/27_mediba.svg -> dirName/d3d3df53.svg
func UrlRewrite(preUrl, afterUrl, content string) string {
	return strings.Replace(content, preUrl, afterUrl, -1)
}
