# 爬虫
## 说明
    基于chromedp实现，其中支持两种抓取方式，下面再做说明
## 使用方法
```golang
spider := huomao_spider.GetSpiderClass()
_, _, index_single_html := spider.SingHtml("https://www.baidu.com", 2)
spider.Html("https://top.baidu.com/board?tab=realtime", config)
```
参考 [爬虫示例](../example/spider "爬虫示例")

## 两种抓取方式

| 方式 | 方法 | 说明 |  
| :---: | :--- | :--- |  
| 列表方式 | Html(url string,**config map[string]string**) | 这种方式传入的url是一个种子页面，就是一个列表页面，抓取列表项和每一个列表项的源码，通过ListSource获取种子页源码，结果是map[string]string结构，key是一个hash值，因为有的时候两个页面的url会相同，比如ajax加载数据的列表页，value是页面的源码，通过DetailSource获取每一个列表项目的源码，结构是map[string]string，key是每一个列表项的url，value是该url的源码 |  
| 单页面方式 | SingHtml(url string) | 通过传入url获取单个页面的源码，有三个返回值：浏览器url，浏览器标题，页面源码 |  

## 列表模式的config参数介绍

### 格式
> map[string]string

### 说明
> 原本是从配置文件中读取配置，后来转成参数格式，但是参数依然可以来源于配置文件读取后传递

### 配置文件格式
#### 这个格式就是参数config的格式
```
[http://www.ccgp-jiangsu.gov.cn/jiangsu/cggg_search.html]
parent = tbody > tr[href] #能找到要抓取列表的部分，也就是通过这个进行页面查找返回就是想要的数据列表，点击找到的列表能到详情页，必须是可点击的，同时以这个结果为父级提取文本
delay = 2
#有的时候需要进行一些动作操作才会出现我们要抓取的列表，这才有了这个action配置
#最外层使用一对"""包裹，[]里面是一个一个页面，每一个页面用{}包含，两个[]就是一次要采集两个页面
#{}里面目前支持四个操作SetValue和Click，SetValue是设置文本框或者选择框等元素的值，Click是点击页面上的元素，配置SetValue和Click有顺序的
#SetValue里面有两个值，第一个是要赋值文本框或者选择框的id或者class，第二个是要赋的值
#Click，有一个值，是要点击的按钮或者标签的id或者class
#下面例子的意思是，每次采集两个页面，第一页面先设置id是region的值为110000，再点一下id是searchBtn的按钮，读取页面内容
#第二个页面是先设置id是region的值为110000，再点一下id是searchBtn的按钮，最后再点一下id是img_next的按钮，读取页面内容
#Captcha，暂只支持一个验证码,共三个值，第一个是验证码图片的，第二个是输入验证码的文本框，第三个是输入验证码要点击的确定按钮
#PageNext，翻页用的，第一个选项是列表项，通过这个列表项能发挥分页条上的元素列表，第二个是点击第几个按钮进行翻页，从1开始
action = """
[
        [
                [{"Wait":["#validateCodeImg"]}],
                [{"Click":["#date > li:nth-child(2) > a"]}],
                [{"Click":["#cgtype > li:nth-child(1) > a"]}],
                [{"SetValue":["#region","110000"]}],
                [{"Captcha":["#validateCodeImg","#validateCode","button[class=q_submit]"]}]
        ],
        [	
                [{"PageNext":[".ui-paging-container > ul > li","10"]}]
        ]
]
"""
```
#### 参数详解
| 参数 | 是否必须 | 值 | 说明 |  
| :---: | :--- | :--- | :--- |  
| parent | 是|-| 这个只有调用列表抓起方法HTML时有效，是能找到要抓取列表的部分，也就是通过这个进行页面查找返回就是想要的数据列表，点击找到的列表能到详情页，必须是可点击的，同时以这个结果链接提取文本 |  
| delay | 否 |-| 有的时候点击列表弹开的新页面为了防止抓起，ajax延后执行，要是抓取早了就不是真正的结果页面，所以有了这个配置项，用于抓取列表页详情源码的时候等待多少秒 |  
| action | 否 || 有的时候需要进行一些动作操作才会出现我们要抓取的列表，这才有了这个action配置，最外层使用一对"""包裹，[]里面是一个一个页面，每一个页面用{}包含，两个[]就是一次要采集两个页面，{}里面目前支持5个操作Wait/Captcha/PageNext/SetValue/Click,这里的配置有顺序执行的，SetValue是设置文本框或者选择框等元素的值，Click是点击页面上的元素，Wait是等待某个元素加载完成，Captcha是有验证码，PageNex是翻页操作 |  
|||Click|Click，有一个值，是要点击的按钮或者标签的id或者class或者直接使用selector|
|||SetValue|SetValue里面有两个值，第一个是要赋值文本框或者选择框的id或者class或者直接使用selector，第二个是要赋的值|
|||Captcha|暂只支持一个验证码,共三个值，第一个是验证码图片的，第二个是输入验证码的文本框，第三个是输入验证码要点击的确定按钮|
|||PageNext|翻页用的，第一个选项是列表项，通过这个列表项能发现分页条上的元素列表，第二个是点击第几个列表项元素进行翻页，从1开始|
|||Wait|有一个值，只有等到这个元素加载完毕了才开始进行采集|

#### 解读上面配置文件的意思
> 等id是validateCodeImg的元素显示出来，然后点击selector是“#date > li:nth-child(2) > a”的元素，然后点击selector是“#cgtype > li:nth-child(1) > a”的元素，然后将id是region的文本框输入一个值“110000”，然后根据selector的id是validateCodeImg找到验证码图片，进行识别，将识别结果放入id是validateCode的文本框中，再点击button[class=q_submit]，最后开始抓取

> 等上面一步抓取完了，点击翻页按钮，通过selector是.ui-paging-container > ul > li找到翻页按钮列表，点击第10个，最后抓取页面

> 等上面两步执行完，返回抓取结果

