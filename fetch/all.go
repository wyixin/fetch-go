package fetch

import (
	"encoding/json"
	"errors"
	"fetch-go/dom"
	"fetch-go/utils"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Fetch struct {
	WG     *sync.WaitGroup
	Input  *FetchInput
	Output *FetchOutput
}

// input base url
// output all contents, including:
// 1. the html response of this base url
// 2. the content of the static files included in the above html response
type FetchInput struct {
	BaseURL string    `json:"url"`
	Time    time.Time `json:"ts"`

	// cookie, proxy ???
}

type StaticFile struct {
	BaseURL  string `json:"url"`
	HashName string `json:"hash_name"`
}

func (f *StaticFile) DownLoadTo(dirName string) (err error) {
	if f.BaseURL == "" {
		return errors.New("baseurl empty!")
	}

	filepath := dirName + "/" + f.HashName
	fmt.Println(filepath)
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(f.BaseURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (f *StaticFile) GetHashName() error {
	if f.BaseURL == "" {
		return errors.New("baseurl empty!")
	}

	u, err := url.Parse(f.BaseURL)
	if err == nil && u.Path != "" {
		f.HashName = utils.HashContent(f.BaseURL) + filepath.Ext(u.Path)
	}
	return nil
}

type FetchOutput struct {
	//	Header: output http header?
	Time       time.Time     `json:"ts"`
	BaseURL    string        `json:"url"`
	CSSFiles   []*StaticFile `json:"css"`
	ImageFiles []*StaticFile `json:"imgs"`
	JSFiles    []*StaticFile `json:"js"`
	Body       string        `json:"-"`
}

func fetchURL(url string) (string, error) {

	fmt.Println(`now processing url:` + url)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (f *Fetch) fulFill(files []*StaticFile, urls []string, dirName string) {

	for k, v := range urls {
		f.WG.Add(1)

		static := &StaticFile{
			BaseURL: v,
		}
		static.GetHashName()

		go func() {
			defer f.WG.Done()
			err := static.DownLoadTo(dirName)
			fmt.Println(err)
		}()

		files[k] = static
	}
}

func (f *Fetch) Wait() {
	f.WG.Wait()
}

func (f *Fetch) FetchALL() error {

	in := f.Input
	dirName := utils.GetDirName(in.BaseURL)
	content, err := fetchURL(in.BaseURL)

	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return err
	}

	images, js, css, err := dom.ParseAllAssets(content)
	if err != nil {
		fmt.Println("Error on parse document string:", err)
		return err
	}

	cssFiles := make([]*StaticFile, len(css))
	jsFiles := make([]*StaticFile, len(js))
	imgFiles := make([]*StaticFile, len(images))

	err = os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return err
	}

	f.fulFill(cssFiles, css, dirName)
	f.fulFill(imgFiles, images, dirName)

	f.Output = &FetchOutput{
		BaseURL:    in.BaseURL,
		CSSFiles:   cssFiles,
		ImageFiles: imgFiles,
		JSFiles:    jsFiles,
		Body:       content,
		Time:       time.Now(),
	}

	return nil
}

func (f *Fetch) SavePage() error {

	out := f.Output
	dirName := utils.GetDirName(out.BaseURL)
	fileName := utils.GetFileName()

	// save static file first
	// then rewrite the links to local file system
	rewritedOutPut := out.Body
	cssFiles := out.CSSFiles
	for _, c := range cssFiles {
		rewritedOutPut = utils.UrlRewrite(c.BaseURL, c.HashName, rewritedOutPut)
	}

	imgFiles := out.ImageFiles
	for _, c := range imgFiles {
		fmt.Println(c.HashName, c.BaseURL)
		rewritedOutPut = utils.UrlRewrite(c.BaseURL, c.HashName, rewritedOutPut)
	}

	// and save
	utils.SaveFile(dirName, fileName, rewritedOutPut)

	// save metadata.json
	toJson, err := json.Marshal(f)
	if err != nil {
		fmt.Println("to json err")
		return err
	}

	utils.SaveFile(dirName, "metadata.json", string(toJson))
	return nil
}

func (f *Fetch) MPrint() error {
	dirName := utils.GetDirName(f.Input.BaseURL)
	filePath := dirName + "/" + "metadata.json"

	data, err := os.ReadFile(filePath)
	if err != nil {
		//		fmt.Println("Error decoding JSON:", err)
		fmt.Println("URL IS NOT REQUESTED YET")
		return err
	}

	err = json.Unmarshal(data, f)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return err
	}

	fmt.Println("Request URL is:", f.Input.BaseURL)
	fmt.Println("Last request time is:", f.Output.Time)
	fmt.Println("CSS File Downloaded:", len(f.Output.CSSFiles))
	fmt.Println("Images Downloaded:", len(f.Output.ImageFiles))
	return nil
}
