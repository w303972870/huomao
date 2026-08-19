package huomao_captcha

type iParseClass interface {
	Class() iParse
}

type iParse interface {
	Parse(img string) string
}

func GetCaptchaClass(which string) iParseClass {
	if which == "yunma" {
		return &yunma{}
	}
	return nil
}
