package main

import (
	"fmt"
	"github.com/w303972870/huomao/spider"
)

func main() {
	config := map[string]string{
		"parent": "a[class=img-wrapper_29V76]",
		"delay":  "2",
	}
	spider := huomao_spider.GetSpiderClass()
	_, _, index_single_html := spider.SingHtml("https://www.baidu.com", 2)
	spider.Html("https://top.baidu.com/board?tab=realtime", config)
	fmt.Println(index_single_html)
	fmt.Println(spider.ListSource)
	fmt.Println(spider.DetailSource)
}
