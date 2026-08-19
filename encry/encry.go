package huomao_encry

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
)

/*
  加密算法示例
    var RsaKeys driver_encry.RsaKeys
    var Encry driver_encry.Encry

    RsaKeys.KeyDir = fmt.Sprint( args.IniConfig.RunDir , "/keys/" )
    RsaKeys.GenerateRSAKey()

    fmt.Println(fmt.Sprint( RsaKeys.KeyDir , "/" , RsaKeys.PrivateFile ))
    fmt.Println(fmt.Sprint( RsaKeys.KeyDir , "/" , RsaKeys.PublicFile ))
    fmt.Println(string(RsaKeys.PublicBlock))
    fmt.Println(string(RsaKeys.PrivateBlock))

    Encry.Data = "id.hf-cloud.cn"
    Encry.Rsa( &RsaKeys )
    fmt.Println(Encry.EncryBase64Data)

    Encry.RsaDecode( &RsaKeys )
    fmt.Println(Encry.Data)
*/

/*高级加密标准AES加密*/
func (e *hencry) AesCBC() ([]byte, error) {
	var k []byte

	if len(e.Key) <= 0 {
		k = []byte(DESENCRYPTKEY)
	} else {
		k = []byte(e.Key)
	}
	//创建加密算法的实例
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}
	data := []byte(e.Data)
	if len(data) <= 0 {
		return nil, errors.New("Data:需要加密的字符串为空")
	}

	//获取块的大小
	e.blockSize = block.BlockSize()

	//对数据进行填充,让数据的长度满足加密需求
	//取余计算长度,判断加密的文本是不是blockSize的倍数,如果不是的话把多余的长度计算出来,用于补齐长度
	padding := e.blockSize - len(data)%e.blockSize
	//补齐
	//Repeat: 把切片[]byte{byte(padding)}复制padding个然后合并成新的字节切片返回
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	e.originData = append(data, padText...)

	//采用aes加密方式中的CBC加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:e.blockSize])
	encryData := make([]byte, len(e.originData))

	//执行加密
	blockMode.CryptBlocks(encryData, e.originData)
	e.EncryBase64Data = base64.StdEncoding.EncodeToString(encryData)
	e.EncryData = string(encryData)

	//返回
	return encryData, nil
}

/*高级加密标准AES解密*/
func (e *hencry) AesCBCDecode() ([]byte, error) {
	var k []byte

	if len(e.Key) <= 0 {
		k = []byte(DESENCRYPTKEY)
	} else {
		k = []byte(e.Key)
	}
	//创建加密算法的实例
	block, err := aes.NewCipher(k)

	if err != nil {
		return nil, err
	}
	if len(e.EncryData) <= 0 {
		return nil, errors.New("EncryData:解密字符串为空")
	}
	//获取块的大小
	e.blockSize = block.BlockSize()

	//创建加密实例
	blockMode := cipher.NewCBCDecrypter(block, k[:e.blockSize])
	originData := make([]byte, len([]byte(e.EncryData)))

	//该函数可也用来加密也可也用来解密
	blockMode.CryptBlocks(originData, []byte(e.EncryData))

	//取出填充的字符串
	//获取数据长度
	length := len(originData)
	if length <= 0 {
		return nil, errors.New("加密字符串长度不符合要求")
	}
	//获取填充字符串的长度
	unPadding := int(originData[length-1])
	if length < unPadding {
		return nil, errors.New("填充字符串长度不符合要求")
	}
	e.Data = string(originData[:(length - unPadding)])
	//截取切片,删除填充的字节,并且返回明文
	return originData[:(length - unPadding)], nil
}

/*高级加密标准AES解密Base64的结果*/
func (e *hencry) AesCBCDecodeBase64() ([]byte, error) {
	if len(e.EncryBase64Data) <= 0 {
		return nil, errors.New("EncryBase64Data:Base64解密字符串为空")
	}
	//解码base64字符串
	pwdByte, err := base64.StdEncoding.DecodeString(e.EncryBase64Data)
	if err != nil {
		return nil, err
	}
	e.EncryData = string(pwdByte)

	//执行aes解密
	return e.AesCBCDecode()
}

/*
 *   高级加密标准AES ECB加密
 *   ECB模式: mysql中AES_DECRYPT函数的实现方式
 *   主要关注三点:
 *       1.调用aes.NewCipher([]byte)是加密关键字key的生成方式, 即下面的generateKey方法
 *       2.分组分块加密的加密方式
 *       3.mysql中一般需要HEX函数来转化数据格式
 *   加密: HEX(AES_ENCRYPT('关键信息', '***—key'))
 *   解密: AES_DECRYPT(UNHEX('关键信息'), '***-key’)
 *   所以调用AESEncrypt或者AESDecrypt方法之后, 使用hex.EncodeToString()转化
 */
func (e *hencry) AesECB() ([]byte, error) {
	var key []byte

	if len(e.Key) <= 0 {
		key = []byte(DESENCRYPTKEY)
	} else {
		key = []byte(e.Key)
	}

	src := []byte(e.Data)
	if len(src) <= 0 {
		return nil, errors.New("Data:需要加密的字符串为空")
	}

	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(src) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, src)
	pad := byte(len(plain) - len(src))
	for i := len(src); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted := make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(src); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}
	e.EncryData = string(encrypted)
	e.EncryBase64Data = base64.StdEncoding.EncodeToString(encrypted)

	return encrypted, nil
}

/*高级加密标准AES ECB解密*/
func (e *hencry) AesECBDecode() ([]byte, error) {
	var key []byte

	if len(e.Key) <= 0 {
		key = []byte(DESENCRYPTKEY)
	} else {
		key = []byte(e.Key)
	}

	encrypted := []byte(e.EncryData)
	if len(encrypted) <= 0 {
		return nil, errors.New("EncryData:需要解密的字符串为空")
	}

	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted := make([]byte, len(encrypted))
	//
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	e.Data = string(decrypted[:trim])
	return decrypted[:trim], nil
}

/*高级加密标准AES解密ECB Base64的结果*/
func (e *hencry) AesECBDecodeBase64() ([]byte, error) {
	if len(e.EncryBase64Data) <= 0 {
		return nil, errors.New("EncryBase64Data:Base64解密字符串为空")
	}
	//解码base64字符串
	pwdByte, err := base64.StdEncoding.DecodeString(e.EncryBase64Data)
	if err != nil {
		return nil, err
	}
	e.EncryData = string(pwdByte)

	//执行aes解密
	return e.AesECBDecode()
}

/*md5加密*/
func (e *hencry) Md5() string {
	if e.Data == "" {
		return ""
	}
	h := md5.New()
	h.Write([]byte(e.Data)) // 需要加密的字符串为
	return hex.EncodeToString(h.Sum(nil))
}

/*md5加密*/
func (e *hencry) Md5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str)) // 需要加密的字符串为
	return hex.EncodeToString(h.Sum(nil))
}

/*生成一个随机数*/
func (e *hencry) Rand() string {
	result, _ := rand.Int(rand.Reader, big.NewInt(100000))
	return result.String()
}

/*HMAC SHA256加密*/
func (e *hencry) HmacSha256() string {
	m := hmac.New(sha256.New, []byte(e.Key))
	m.Write([]byte(e.Data))
	return hex.EncodeToString(m.Sum(nil))
}

/*HMAC SHA512加密*/
func (e *hencry) HmacSha512() string {
	m := hmac.New(sha512.New, []byte(e.Key))
	m.Write([]byte(e.Data))
	return hex.EncodeToString(m.Sum(nil))
}

/*HMAC SHA1加密*/
func (e *hencry) HmacSha1() string {
	m := hmac.New(sha1.New, []byte(e.Key))
	m.Write([]byte(e.Data))
	return hex.EncodeToString(m.Sum(nil))
}

/*SHA256加密*/
func (e *hencry) Sha256() string {
	h := sha256.New()
	h.Write([]byte(e.Data))
	return hex.EncodeToString(h.Sum(nil))
}

/*SHA512加密*/
func (e *hencry) Sha512() string {
	h := sha512.New()
	h.Write([]byte(e.Data))
	return hex.EncodeToString(h.Sum(nil))
}

/*base64加密*/
func (e *hencry) Base64Encode() string {
	return string(base64.StdEncoding.EncodeToString([]byte(e.Data)))
}

/*base64解密*/
func (e *hencry) Base64Decode() string {
	a, err := base64.StdEncoding.DecodeString(e.EncryBase64Data)
	if err != nil {
		return ""
	}
	return string(a)
}

// 读取公/私钥
func (e *hencry) readKeys(path string) ([]byte, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, stat.Size())
	file.Read(buf)
	defer file.Close()
	return buf, nil
}

// 获取公钥
func (e *hencry) PublicBlock(rsakeys *RsaKeys) []byte {
	if rsakeys.PublicBlock == nil {
		rsakeys.PublicBlock, _ = e.readKeys(fmt.Sprint(rsakeys.KeyDir, "/public.pem"))
	}
	return rsakeys.PublicBlock
}

// 获取私钥
func (e *hencry) PrivateBlock(rsakeys *RsaKeys) []byte {
	if rsakeys.PrivateBlock == nil {
		rsakeys.PrivateBlock, _ = e.readKeys(fmt.Sprint(rsakeys.KeyDir, "/private.pem"))
	}
	return rsakeys.PrivateBlock
}

// 加密
func (e *hencry) Rsa(rsakeys *RsaKeys) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(e.PublicBlock(rsakeys))
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	result, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(e.Data))
	if err != nil {
		return nil, err
	}
	e.EncryData = string(result)
	e.EncryBase64Data = base64.StdEncoding.EncodeToString([]byte(e.EncryData))
	//加密
	return result, nil
}

// 解密
func (e *hencry) RsaDecode(rsakeys *RsaKeys) ([]byte, error) {
	//解密
	block, _ := pem.Decode(e.PrivateBlock(rsakeys))

	if block == nil {
		return nil, errors.New("private key error!")
	}

	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	result, err := rsa.DecryptPKCS1v15(rand.Reader, priv, []byte(e.EncryData))
	if err != nil {
		return nil, err
	}
	e.Data = string(result)

	// 解密
	return result, nil
}

/*
 * 生成RSA私钥和公钥，保存到文件中
 * @bits 证书的大小
 * @dir 证书保存路径
 */
func (e *RsaKeys) GenerateRSAKey() {
	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	//Reader是一个全局、共享的密码用强随机数生成器
	if e.Bits == 0 {
		e.Bits = 2048
	}
	if e.PrivateFile == "" {
		e.PrivateFile = "private.pem"
	}
	if e.PublicFile == "" {
		e.PublicFile = "public.pem"
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, e.Bits)
	if err != nil {
		panic(err)
	}
	//保存私钥
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	//使用pem格式对x509输出的内容进行编码
	//创建文件保存私钥
	privateFile, err := os.Create(fmt.Sprint(e.KeyDir, "/", e.PrivateFile))
	if err != nil {
		panic(err)
	}
	defer privateFile.Close()
	//构建一个pem.Block结构体对象
	privateBlock := pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}
	//将数据保存到文件
	pem.Encode(privateFile, &privateBlock)
	e.PrivateBlock = pem.EncodeToMemory(&privateBlock)

	//保存公钥
	//获取公钥的数据
	publicKey := privateKey.PublicKey
	//X509对公钥编码
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	//pem格式编码
	//创建用于保存公钥的文件
	publicFile, err := os.Create(fmt.Sprint(e.KeyDir, "/", e.PublicFile))
	if err != nil {
		panic(err)
	}
	defer publicFile.Close()
	//创建一个pem.Block结构体对象
	publicBlock := pem.Block{Type: "RSA Public Key", Bytes: X509PublicKey}
	//保存到文件
	pem.Encode(publicFile, &publicBlock)
	e.PublicBlock = pem.EncodeToMemory(&publicBlock)

	//fmt.Println("RSA Public Key:" , fmt.Sprint( e.KeyDir , "/" , e.PublicFile ) )
	//fmt.Println("RSA Private Key:" , fmt.Sprint( e.KeyDir , "/" , e.PrivateFile ) )
}

func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}
