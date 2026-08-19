package huomao_otp

type otp struct {
	//二维码的尺寸
	Width  int
	Height int
}

var p *otp

type iOtp interface {

	// 获取秘钥
	GetSecret() string

	// 获取动态码
	GetCode(secret string) (string, error)

	// 获取动态码二维码内容
	GetQrcode(name string, ukey string) string

	// 获取动态码二维码图片地址,这里是第三方二维码api
	GetQrcodeUrl(name string, ukey string) string

	// 验证动态码
	VerifyCode(secret, code string) (bool, error)
}

func GetOtpClass() *otp {
	if p == nil {
		p = &otp{}
	}
	return p
}
