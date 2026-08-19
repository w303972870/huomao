package huomao_encry

var p *hencry

/*AES加密算法的默认密钥*/
const DESENCRYPTKEY = "-JAbO_*4@LJrtfG4EB7B^!3ky$/^Mn1P"

type iEncry interface {
	/*高级加密标准AES加密*/
	AesCBC() ([]byte, error)

	/*高级加密标准AES解密*/
	AesCBCDecode() ([]byte, error)

	/*高级加密标准AES解密Base64的结果*/
	AesCBCDecodeBase64() ([]byte, error)

	/*高级加密标准AES ECB加密*/
	AesECB() ([]byte, error)

	/*高级加密标准AES ECB解密*/
	AesECBDecode() ([]byte, error)

	/*高级加密标准AES解密ECB Base64的结果*/
	AesECBDecodeBase64() ([]byte, error)

	/*md5加密*/
	Md5() string

	/*md5加密*/
	Md5Str(str string) string

	/*生成一个随机数*/
	Rand() string

	/*HMAC SHA256加密*/
	HmacSha256() string

	/*HMAC SHA512加密*/
	HmacSha512() string

	/*HMAC SHA1加密*/
	HmacSha1() string

	/*SHA256加密*/
	Sha256() string

	/*SHA512加密*/
	Sha512() string

	/*base64加密*/
	Base64Encode() string

	/*base64解密*/
	Base64Decode() string

	// 获取公钥
	PublicBlock(rsakeys *RsaKeys) []byte

	// 获取私钥
	PrivateBlock(rsakeys *RsaKeys) []byte

	// 加密
	Rsa(rsakeys *RsaKeys) ([]byte, error)

	// 解密
	RsaDecode(rsakeys *RsaKeys) ([]byte, error)

	/* 生成RSA私钥和公钥，保存到文件中 */
	GenerateRSAKey()
}

/*
 *   对称加密, 加解密都使用的是同一个密钥, 其中的代表就是AES
 *   非对加解密, 加解密使用不同的密钥, 其中的代表就是RSA
 *   签名算法, 如MD5、SHA1、HMAC等, 主要用于验证，防止信息被修改, 如：文件校验、数字签名、鉴权协议
 */
type hencry struct {
	Data            string // 原始字符串
	EncryData       string // 加密后的字符串
	EncryBase64Data string // base64的加密后的字符串
	Key             string // 加密密钥
	blockSize       int    // 块的大小
	originData      []byte // 原始字符串处理过程的中间字符串
}

type RsaKeys struct {
	Bits         int    // 加密字长，一般赋值1024，2048
	KeyDir       string // 放置加密文件的路径
	PrivateBlock []byte // 生成的私钥
	PrivateFile  string // 生成的私钥文件名
	PublicBlock  []byte // 生成的公钥
	PublicFile   string // 生成的公钥文件名
}

func GetEncryClass() *hencry {
	return &hencry{}
}
