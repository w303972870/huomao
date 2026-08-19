package main

import (
	"encoding/base64"
	"github.com/w303972870/huomao/captcha"
	"io/ioutil"
	"net/http"
)

func main() {
	resp, _ := http.Get("https://www.jfbym.com/static/pic/sy/c1-4/1.png")
	body, _ := ioutil.ReadAll(resp.Body)
	captcha_class := huomao_captcha.GetCaptchaClass("yunma").Class()
	captcha_result := captcha_class.Parse(base64.StdEncoding.EncodeToString([]byte(body)))
	fmt.Println("验证码识别结果：", captcha_result)
}
