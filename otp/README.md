# OTP
## 说明
    使用Golang实现的otp算法
## 使用方法
```golang
otp := huomao_otp.GetOtpClass()
secret := otp.GetSecret()
```
参考 [OTP示例](../example/otp "OTP示例")

## 方法介绍

| 方法 | 返回值 | 说明 |  
| :---: | :--- | :--- |  
| GetSecret()  | string | 获取一个密钥，这个密钥需要记住，且不能外泄，要严格保密，后续计算比较都需要用到 |  
| GetCode(secret string)  | string, error | 根据密钥获取一个实时的动态码 |  
| GetQrcode(name string, ukey string)   | string | 传入一个名称和密钥，名称是用于展示谁的二维码的名称，一般是用户登录账号，生成一个otp专属的链接，拿到这个链接就可以生产二维码了 |  
| GetQrcodeUrl(name string, ukey string)   | string| 参数跟上面一样，区别就是返回的是一个二维码图片的http url |  
| VerifyCode(secret, code string)   | bool, error| 验证动态码是否有效 |  


