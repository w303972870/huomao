package huomao_spider

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	//"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/runtime"
	"github.com/w303972870/huomao/captcha"
	//"strings"
	"crypto/md5"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	//"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/gookit/color"
	"strconv"
	//"sync"
	"time"
	//"github.com/chromedp/chromedp/kb"
	//"io/ioutil"
	"log"
	//"mt/common"
	//"net/url"
	//"path/filepath"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/target"
	//"net/http"
	//sys "runtime"
	//"net/http/httptest"
)

/**
 * 初始化浏览器
 */
func (c *chromeSpider) brower(headless bool) (context.Context, context.CancelFunc) {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent(`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.64 Safari/537.36`),
		chromedp.Flag("headless", headless),
		chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,                   //不检查默认浏览器
		chromedp.Flag("ignore-certificate-errors", true), //忽略错误
		chromedp.Flag("disable-web-security", true),      //禁用网络安全标志
		chromedp.NoFirstRun,                              //不是首次运行
	)
	allocatorContext, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, _ := chromedp.NewContext(allocatorContext, chromedp.WithLogf(log.Printf))
	//ctx, _ = context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second) //设置后续操作读取网页包含设置的sleep不超过30秒，否则报错context deadline exceeded
	return ctx, cancel
}

/**
 * 单独页面请求
 * delay是有的页面ajax有延迟，不做延迟获取得不到完整内容
 * 返回：浏览器中的url，浏览器中的title，页面源码
 */
func (c *chromeSpider) SingHtml(url string, delay int) (string, string, string) {
	ctx, cancle := c.brower(true)
	defer cancle()
	chromedp.Run(ctx, chromedp.Navigate(url))

	if delay > 0 {
		chromedp.Run(ctx, chromedp.Sleep(time.Duration(delay)*time.Second))
	}
	var html_source string
	chromedp.Run(ctx,
		chromedp.Location(&c.currentUrl),
		chromedp.Title(&c.firstTitle),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, _ := dom.GetDocument().WithDepth(-1).Do(ctx)
			html_source, _ = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			return nil
		}),
	)
	return c.currentUrl, c.firstTitle, html_source
}

/**
 * 初始化配置
 */
func (c *chromeSpider) init() {
	c.newTargetUrl = make(chan string)
	c.ListSource = make(map[string]string)
	c.DetailSource = make(map[string]string)
	c.context, _ = c.brower(true)
}

/**
 * 执行抓取
 */
func (c *chromeSpider) Html(urls string, detail map[string]string) {

	c.config = detail
	c.actionList = ParseJsonString(c.config["action"])
	c.init()
	chromedp.Run(c.context, chromedp.Navigate(urls))

	sss, _ := chromedp.Targets(c.context)
	for _, t := range sss {
		c.currentTargetId = t.TargetID
	}

	for i := 0; i <= len(c.actionList); i++ {

		c.pageNum = -1
		c.pageNodes = make([]*cdp.Node, 1)
		c.changePageTask = chromedp.Tasks{}
		c.pageListLiNodes = make([]*cdp.Node, 1)

		if len(c.actionList) > 0 {
			if i == len(c.actionList) {
				break
			}
			chromedp.Run(c.context, c.ChromeDo(c.actionList[i]))
		}
		if c.pageNum > -1 {
			if c.nextPageAction() == false {
				continue
			}
		}

		color.Tag("mgaB").Print("开始列表页")
		if c.listAction() == true {
			c.detailAction()
		}

		if c.getDetailSource() == false {
			continue
		}

		c.listHtmlSource(c.currentTargetId)

		c.waitTargetSync.Wait()
		if c.pageNum != -1 {
			c.closeTarget(&c.context, c.currentTargetId)
		}
	}
}

/**
 * 翻页
 */
func (c *chromeSpider) nextPageAction() bool {
	chromedp.Run(c.context, c.changePageTask)
	if len(c.pageNodes) < c.pageNum-1 {
		color.Warn.Prompt("请确认PageNext的第一个选项配置是否有误，未定位到翻页，跳过翻页1")
		return false
	} else {
		chromedp.Run(c.context,
			chromedp.MouseClickNode(c.pageNodes[c.pageNum-1], chromedp.ButtonLeft),
		)
		color.Info.Tips("执行翻页成功")
	}
	return true
}

/**
 * 查找详情页
 */
func (c *chromeSpider) listAction() bool {
	color.Tag("ylw0").Print("...查找列表")
	chromedp.Run(c.context, chromedp.Sleep(5*time.Second),
		chromedp.WaitReady(c.config["parent"]),
		chromedp.Nodes(c.config["parent"], &c.pageListLiNodes),
	)
	if len(c.pageListLiNodes) > 0 {
		c.newTargetUrl = make(chan string, len(c.pageListLiNodes))
		color.Tag("ylw0").Print("...已找到...搜集列表", len(c.pageListLiNodes), "个")
	} else {
		color.Tag("red").Print("搜集列表为空")
		return false
	}

	return true
}

/**
 * 激活详情页
 */
func (c *chromeSpider) detailAction() {

	for _, n := range c.pageListLiNodes {
		if n == nil {
			continue
		}

		ch := c.addNewTabListener()

		var nnode []cdp.NodeID
		nnode = append(nnode, n.NodeID)
		chromedp.Run(c.context,
			chromedp.Click(nnode, chromedp.ByNodeID),
		)
		/*
			if len(ch) == 0 {
				modifier := input.ModifierCtrl
				if sys.GOOS == "darwin" {
					modifier = input.ModifierCommand
				}
				chromedp.Run(c.context,
					chromedp.MouseClickNode(n, chromedp.ButtonModifiers(modifier)),
				)
			}
		*/
		<-ch
		color.Tag("mgaB").Print("·")
	}
	color.Tag("mgaB").Println("OK")
}

/**
 * 获取详情页信息
 */
func (c *chromeSpider) getDetailSource() bool {
	listTargets := c.listTargets()
	if len(listTargets) < len(c.pageListLiNodes) || len(c.newTargetUrl) == 0 {
		color.Warn.Prompt("配置项parent可能有问题，定位到的原始点击打不开详情页面2")
		return false
	}

	color.Tag("blue").Print("开始获取详情页源码", len(listTargets), "个")

	for _, t := range listTargets {
		if c.currentTargetId == t.TargetID {
			continue
		}
		color.Tag("blue").Print("·")
		c.waitTargetSync.Add(1)
		newurl := <-c.newTargetUrl
		go c.detailHtmlSource(t.TargetID, newurl)

		newctx, _ := chromedp.NewContext(c.context, chromedp.WithTargetID(t.TargetID))
		chromedp.Run(newctx, chromedp.ActionFunc(func(ctxt context.Context) error {
			c.closeTarget(&newctx, t.TargetID)
			return nil
		}))

	}
	//c.waitTargetSync.Wait()
	color.Tag("blue").Print("OK")
	return true
}

/**
 * 获取列表页信息,返回存储源码的key
 */
func (c *chromeSpider) listHtmlSource(t target.ID) (htmlSourceKey string) {
	color.Tag("blue").Print("...进入列表页")
	newCtx, _ := c.activate(t)
	color.Tag("blue").Println("OK")
	//chromedp.RunResponse(*newCtx)
	chromedp.Run(*newCtx,
		chromedp.Location(&c.currentUrl),
		chromedp.Title(&c.firstTitle),
		chromedp.Sleep(1*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, _ := dom.GetDocument().WithDepth(-1).Do(ctx)
			html_list_source, _ := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			h := md5.New()
			h.Write([]byte(html_list_source)) // 需要加密的字符串为
			htmlSourceKey = hex.EncodeToString(h.Sum(nil))
			c.ListSource[htmlSourceKey] = html_list_source
			return nil
		}),
	)
	//ioutil.WriteFile(fmt.Sprint(t, ".html"), []byte(html_list_source), 0777)
	return
}

/**
 * 获取详情页信息,详情页使用url作为键值，返回键值
 */
func (c *chromeSpider) detailHtmlSource(t target.ID, url string) string {
	var htmlsource string
	/*newCtx, _ := chromedp.NewContext(c.context, chromedp.WithTargetID(t))
	chromedp.Run(newCtx,
		//chromedp.Sleep(2*time.Second),
		chromedp.Location(&c.currentUrl),
		chromedp.Title(&c.firstTitle),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, _ := dom.GetDocument().WithDepth(-1).Do(ctx)
			htmlsource, _ = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			c.DetailSource[c.currentUrl] = htmlsource
			if c.currentUrl == "about:blank" {
				color.Info.Tips("about:blank：可能是网络问题导致页面未加载完成")
			}
			c.closeTarget(&ctx, t)
			c.waitTargetSync.Done()
			return nil
		}),
	)*/

	delay, _ := strconv.Atoi(c.config["delay"])
	if delay > 0 {
		_, _, htmlsource = c.SingHtml(url, delay)
	} else {
		_, _, htmlsource = c.SingHtml(url, 0)
	}
	c.DetailSource[url] = htmlsource
	c.waitTargetSync.Done()

	/*
		if delay > 0 {
			chromedp.Run(c.context, chromedp.Sleep(time.Duration(delay)*time.Second))
		}
		chromedp.Run(*newCtx,
			//chromedp.Sleep(2*time.Second),
			chromedp.Location(&c.currentUrl),
			chromedp.Title(&c.firstTitle),
			chromedp.ActionFunc(func(ctx context.Context) error {
				i := 0
				for {
					if c.currentUrl != "about:blank" {
						node, _ := dom.GetDocument().WithDepth(-1).Do(ctx)
						html_list_source, _ := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
						c.DetailSource[c.currentUrl] = html_list_source
						break
					} else {
						chromedp.Run(ctx, chromedp.Location(&c.currentUrl), chromedp.Title(&c.firstTitle), chromedp.Sleep(1*time.Second))
					}
					i++
					if i > 5 {
						break
					}
				}
				if c.currentUrl == "about:blank" {
					color.Info.Tips("about:blank：可能是网络问题导致页面未加载完成")
				}
				c.closeTarget(newCtx, t)
				c.waitTargetSync.Done()
				return nil
			}),
		)
	*/
	//ioutil.WriteFile(fmt.Sprint(t, ".html"), []byte(htmlsource), 0777)
	return c.currentUrl
}

/**
 * 关闭标签
 */
func (c *chromeSpider) closeTarget(ctxt *context.Context, t target.ID) error {
	chromedp.Cancel(*ctxt)
	err := chromedp.Run(*ctxt, chromedp.ActionFunc(func(ctxt context.Context) error {
		target.CloseTarget(t).Do(ctxt)
		return nil
	}))
	if err != nil {
		//fmt.Println(err)
	}
	return nil
}

/**
 * 激活标签，返回激活标签的context
 */
func (c *chromeSpider) activate(t target.ID) (*context.Context, context.CancelFunc) {
	err := chromedp.Run(c.context, chromedp.ActionFunc(func(ctxt context.Context) error {
		err := target.ActivateTarget(t).Do(ctxt)
		if err != nil {
			return err
		}
		return nil
	}))

	if err != nil {
		fmt.Println(err)
	}
	newCtx, cancel := chromedp.NewContext(c.context, chromedp.WithTargetID(t))
	return &newCtx, cancel
}

/**
 * 返回标签列表
 */
func (c *chromeSpider) listTargets() []*target.Info {
	targets, err := chromedp.Targets(c.context)
	if err != nil {
		fmt.Println(err)
	}
	var out []*target.Info
	for _, v := range targets {
		//fmt.Println("target: ", v.TargetID.String(), " type: ", v.Type, " url: ", v.URL)
		if v.Type != "page" {
			color.Info.Tips("target of type not page")
			continue
		}
		out = append(out, v)
	}
	return out
}

/**
 * 开启网络请求监测
 */
func (c *chromeSpider) addNetWorkListener() {
	chromedp.ListenTarget(c.context, func(v interface{}) {
		switch ev := v.(type) {
		case *network.EventRequestWillBeSent:
			fmt.Println(ev.Request.Method, ev.Request.URL, ev.Request.PostData, ev.Request.Headers)
		case *page.EventJavascriptDialogOpening:
			fmt.Println("closing alert:", ev.Message)
			t := page.HandleJavaScriptDialog(true) //false自动关闭alert对话框
			go func() {
				if err := chromedp.Run(c.context, t); err != nil {
					// 关闭失败或出错了
				}
			}()
		case *runtime.EventConsoleAPICalled:
			fmt.Printf("* console.%s call:\n", ev.Type)
			for _, arg := range ev.Args {
				fmt.Printf("%s - %s\n", arg.Type, arg.Value)
			}
		case *runtime.EventExceptionThrown:
			// Since ts.URL uses a random port, replace it.
			s := ev.ExceptionDetails.Error()
			fmt.Printf("* %s\n", s)
		}
	})
}

/**
 * 注册新tab标签的监听服务
 */
func (c *chromeSpider) addNewTabListener() <-chan target.ID {
	return chromedp.WaitNewTarget(c.context, func(info *target.Info) bool {
		//fmt.Println(info.URL)
		if info.URL != "" {
			c.newTargetUrl <- info.URL
		}
		return info.URL != ""
	})
}

/**
 * 截全屏
 */
func (c *chromeSpider) Screenshot() []byte {
	var screen_shot []byte
	chromedp.Run(c.context,
		chromedp.CaptureScreenshot(&screen_shot),
	)
	return screen_shot
}

/*
设置以下动作是按顺序的，sleep是为了有的ajax加载不完
*/
func (c *chromeSpider) ChromeDo(do_list [][]map[string][]string) chromedp.Tasks {
	chrome_task := chromedp.Tasks{}

	for _, list := range do_list {
		if len(list) < 1 {
			continue
		}
		for _, do := range list {
			if len(do) < 1 {
				continue
			}
			for key, action := range do {
				if key == "SetValue" && len(action) > 1 {
					chrome_task = append(chrome_task, chromedp.WaitReady(action[0]), chromedp.SetValue(action[0], action[1]))
				} else if key == "Click" && len(action) > 0 {
					chrome_task = append(chrome_task,
						chromedp.WaitReady(action[0], chromedp.ByQuery),
						chromedp.Click(action[0], chromedp.ByQuery),
					)
				} else if key == "Captcha" && len(action) > 2 {
					chrome_task = append(chrome_task, chromedp.WaitVisible(action[0]),
						chromedp.WaitReady(action[0]), chromedp.WaitReady(action[1]), chromedp.WaitReady(action[2]))
					chrome_task = append(chrome_task, chromedp.Screenshot(action[0], &captcha_img,
						chromedp.NodeVisible, chromedp.ByID), c.action_captcha(action))
				} else if key == "Wait" && len(action) > 0 {
					chrome_task = append(chrome_task, chromedp.WaitVisible(action[0], chromedp.ByQuery))
				} else if key == "PageNext" && len(action) > 1 {
					c.pageNum, _ = strconv.Atoi(action[1])
					c.changePageTask = append(c.changePageTask, chromedp.WaitVisible(action[0]), chromedp.Nodes(action[0], &c.pageNodes))
				}
			}
		}
	}
	return chrome_task
}

/*操作验证码*/
func (c *chromeSpider) action_captcha(action []string) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(context.Context) error {
		captcha_class := huomao_captcha.GetCaptchaClass("yunma").Class()
		result := captcha_class.Parse(base64.StdEncoding.EncodeToString(captcha_img))
		task := chromedp.Tasks{}
		task = append(task, chromedp.SendKeys(action[1], result, chromedp.ByID), chromedp.Click(action[2])) //chromedp.Sleep(5*time.Second),

		chromedp.Run(c.context, task)
		//err := ioutil.WriteFile("captcha_img.png", captcha_img, 0777)
		return nil
	})
}
