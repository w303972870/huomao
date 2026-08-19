package huomao_spider

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"sync"
)

type chromeSpider struct {
	context      context.Context
	Debug        bool
	ListSource   map[string]string `列表页源码`
	DetailSource map[string]string `详情页源码`

	config          map[string]string         `url抓取配置详情`
	actionList      [][][]map[string][]string `临时变量，动作配置列表`
	currentUrl      string                    `临时变量，当前标签页的url`
	firstTitle      string                    `临时变量，当前标签页的<title>`
	pageNodes       []*cdp.Node               `临时变量,储存翻页按钮列表`
	pageListLiNodes []*cdp.Node               `临时变量,储存抓取列表`
	captchaImg      []byte                    `临时变量,存储验证码`
	currentTargetId target.ID                 `临时变量,记录当前激活的列表页`
	pageNum         int                       `临时变量,需要第几页的列表页`
	changePageTask  chromedp.Tasks            `临时变量,翻页`
	waitTargetSync  sync.WaitGroup            `任务列表`
	newTargetUrl    chan string
}

var need_next_page int
var change_page_task chromedp.Tasks
var screen_shot, captcha_img []byte

func GetSpiderClass() *chromeSpider {
	return &chromeSpider{}
}
