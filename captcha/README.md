# 验证码识别
## 说明
    方法暂时实现了云码接口方式，后续再实现其他OCR和自建模型识别
## 使用方法
```golang
# 通过传入参数yunma来调用云码api的方式，目前仅支持这一个参数
captcha_class := huomao_captcha.GetCaptchaClass("yunma").Class()
captcha_result := captcha_class.Parse(base64.StdEncoding.EncodeToString([]byte(body)))
```
参考 [验证码识别示例](../example/captcha "验证码识别示例")

## 支持方式
- [x] 云码
- [ ] 百度
- [ ] 自建模型