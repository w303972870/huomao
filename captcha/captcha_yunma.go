package huomao_captcha

import (
	"bytes"
	//"encoding/base64"
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"net/http"
	//"time"
)

type yunma struct{}

func (y *yunma) Class() iParse {
	return &yunmaParse{
		customUrl: "https://www.jfbym.com/api/YmServer/customApi",
		token:     "",
	}
}

type yunmaParse struct {
	customUrl string
	token     string
}

type result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data data   `json:"data"`
}
type data struct {
	Code      int    `json:"code"`
	CaptchaId string `json:"captchaId"`
	RecordId  string `json:"recordId"`
	Data      string `json:"data"`
}

func (y *yunmaParse) Parse(img string) string {
	config := map[string]interface{}{}
	config["image"] = img
	config["type"] = "10110"
	config["token"] = y.token
	configData, _ := json.Marshal(config)
	body := bytes.NewBuffer([]byte(configData))
	resp, _ := http.Post(y.customUrl, "application/json;charset=utf-8", body)
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)

	var res result
	_ = json.Unmarshal([]byte(data), &res)
	if 10000 == res.Code {
		return res.Data.Data
	}
	return ""
}
