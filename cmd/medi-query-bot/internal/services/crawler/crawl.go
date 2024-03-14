package crawler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"med-chat-bot/internal/meta"
	"med-chat-bot/pkg/cfg"
	"net/http"
	"time"
)

type ICrawlerService interface {
	Crawl(c *gin.Context) (*meta.BasicResponse, error)
}

type crawlerService struct {
}

type CrawlerServiceArgs struct {
}

func NewCrawlerService(args CrawlerServiceArgs) ICrawlerService {
	return &crawlerService{}
}

func (_this *crawlerService) Crawl(c *gin.Context) (*meta.BasicResponse, error) {
	service, err := selenium.NewChromeDriverService(viper.GetString(cfg.ConfigSeleniumPath), 4444)
	if err != nil {
		panic(err)
	}
	defer service.Stop()
	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := chrome.Capabilities{
		Prefs: map[string]interface{}{"profile.managed_default_content_settings.images": 2},
		Path:  "",
		Args:  []string{
			//"--headless",
			//"--window-size=1200x600",
		},
	}
	caps.AddChrome(chromeCaps)

	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		panic(err)
	}
	defer driver.Quit()
	err = driver.SetImplicitWaitTimeout(10 * time.Second)
	if err != nil {
		return nil, err
	}
	err = driver.Get("https://vnexpress.net/")
	if err != nil {
		return nil, err
	}
	time.Sleep(3 * time.Second)
	title, err := driver.Title()
	if err != nil {
		panic(err)
	}

	// Print the title
	fmt.Println("Title:", title)

	response := &meta.BasicResponse{
		Meta: meta.Meta{
			Code:    http.StatusOK,
			Message: "Success",
		},
		Data: nil,
	}
	return response, nil
}
