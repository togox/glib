
package trps

import (
	"github.com/sclevine/agouti"
	"io/ioutil"
	"log"
	"github.com/cooleo/slugify"
	"time"
	"os"
)

type Options struct {
	Url string
	Sitename string
	SeleniumServerURL string
	Browser string
}

func delaySecond(n time.Duration) {
	time.Sleep(n * time.Second)
}

func Fetch(options Options) (*os.File, string){

	fileName := slugify.Slugify(options.Url)  + ".html"
	capabilities := agouti.NewCapabilities().Browser(options.Browser).With("javascriptEnabled")
	page, err := agouti.NewPage(options.SeleniumServerURL, agouti.Desired(capabilities))
	if err != nil {
		log.Println("start triposo crawl failed to open page", options.Url, " with error :",err)
	}
	if page != nil {
		err := page.Navigate(options.Url)
		if err != nil {
			log.Println("Failed to navigate:", err)
		} else {
			log.Println("Start triposo crawling with Url :", options.Url)
			sum := 1
			loadMore := page.Find("#more-link")
			text, _ := loadMore.Text()
			log.Println("button text:%s", text)
			delaySecond(3)
			for text != "" {
				sum += 1
				if(sum >= 50) {
					break;
				}
				log.Println("do loading more with Url ", options.Url, ", total page:", sum)
				loadMore.MouseToElement()
				delaySecond(2)
				loadMore := page.Find("#more-link")
				text, _ := loadMore.Text()
				log.Println("button text:%s", text)
				if text != "" {
					loadMore.Click()
					delaySecond(2)
				} else {
					log.Println("do loading more with Url ", options.Url, ", button text null")
					break
				}
			}
			log.Println("start getting page html and save to S3")
			data, _ := page.HTML()
			html := []byte(data)
			err := ioutil.WriteFile(fileName, html, 0644)

			if err != nil {
				log.Println("Error :%s", err)
			} else {
				file, err := os.Open(fileName) // For read access.
				if err != nil {
					log.Fatal(err)
				} else {
					page.Destroy()
					return file, fileName
				}
			}
		}
	} else {
		log.Println("Page null")
		page.Destroy()
		return nil,""
	}

	return nil,""
}
