package main

import (
	"fmt"
	"github.com/w303972870/huomao/otp"
)

func main() {
	otp := huomao_otp.GetOtpClass()
	secret := otp.GetSecret()
	fmt.Println("生成密钥：", secret)

	code, _ := otp.GetCode(secret)
	fmt.Println("生成一个动态码：", code)

	ok, _ := otp.VerifyCode(secret, code)
	if ok == true {
		fmt.Println("校验动态码", code, "正确")
	} else {
		fmt.Println("校验动态码", code, "错误")
	}

	code, _ = otp.GetCode(secret)
	fmt.Println("重新生成一个动态码：", code)

	ok, _ = otp.VerifyCode(secret, code)
	if ok == true {
		fmt.Println("校验动态码", code, "正确")
	} else {
		fmt.Println("校验动态码", code, "错误")
	}

	qrcode := otp.GetQrcode("wangdianchen", secret)
	fmt.Println("生成otp：", qrcode)

	qrcodeurl := otp.GetQrcodeUrl("wangdianchen", secret)
	fmt.Println("生成二维码url：", qrcodeurl)

}
