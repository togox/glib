package sobr

import (
	"github.com/sclevine/agouti"
	"io/ioutil"
	"log"
	"github.com/cooleo/slugify"
	"time"
	"os"
	"fmt"
)

type Options struct {
	Url string
	Sitename string
	SeleniumServerURL string
	Browser string
	Username string
	Password string
}


func delaySecond(n time.Duration) {
	time.Sleep(n * time.Second)
}

func clickLoadMore(page *agouti.Page) {
	loadMoreButton := page.Find(".show-more-button");
	if (loadMoreButton != nil) {
		 loadMoreButton.Click();
	}
}
func loopLoadMore(page *agouti.Page){
	clickLoadMore(page)
	delaySecond(2)
	loadMore := page.Find(".show-more-button")


	if (loadMore != nil) {
		  text, _ := loadMore.Text()
			fmt.Println("loadMore Text:", text)
			if (text != ""){
			   loopLoadMore(page);
			}
	}
}

func Fetch(options Options) (*os.File, string,string){
	fileName := slugify.Slugify(options.Url)  + ".html"
	fileNameJson := slugify.Slugify(options.Url) + ".json"

	capabilities := agouti.NewCapabilities().Browser(options.Browser).With("javascriptEnabled")
	page, err := agouti.NewPage(options.SeleniumServerURL, agouti.Desired(capabilities))
	if err != nil {
		log.Println("start triposo crawl failed to open page", options.Url, " with error :",err)
	}
	if page != nil {
		err := page.Navigate(options.Url)
		if err != nil {
			fmt.Println("Failed to navigate:", err)
		} else {
			fmt.Println("stawrt")
		}
		page.Find("body > div.page > header > div.in > div > a.btn.fb.sub").Click();
		page.Find("body > div > section > div > div > div > div.item-list > div.item-big > a").Click();
		page.Find("#email").SendKeys(options.Username)
		page.Find("#pass").SendKeys(options.Password)
		page.Find("#loginbutton").Click()
		loopLoadMore(page)
		data, _ := page.HTML()
		html := []byte(data)
		writeErr := ioutil.WriteFile(fileName, html, 0644)
		if writeErr != nil {
			log.Println("Error :%s", writeErr)
		} else {
			file, err := os.Open(fileName) // For read access.
			if err != nil {
				log.Fatal(err)
			} else {
				page.Destroy()
				return file, fileName,fileNameJson
			}
		}
	} else {
		log.Println("Page null")
		return nil,"",""
	}
	return nil,"",""
}
