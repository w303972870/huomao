# 加解密算法类库
## 说明
    综合了rsa，aes，md5等加解密算法
## 使用方法
```golang
encry := huomao_encry.GetEncryClass()
encry.Data = "原始字符串"
encry.Key = "密钥"
encry.AesCBC()
```
参考 [加解密示例](../example/encry "加解密算法示例")

## encry结构体
```
type encry struct {
    Data            string // 原始字符串
    EncryData       string // 加密后的字符串
    EncryBase64Data string // base64的加密后的字符串
    Key             string // 加密密钥
    blockSize       int    // 块的大小
    originData      []byte // 原始字符串处理过程的中间字符串
}

```

## encry结构体主要属性介绍

| 属性名称  | 说明 |  
| :---: | :--- |  
| Data  | 原始字符串，就是要加密的字符串 |  
| EncryData  | 将原始字符串加密后的结果，或者要解密的字符串也赋值给这个 |  
| EncryBase64Data  | 将原始字符串加密后的结果进行base64，或者要解密的字符串是base64的需要base64解密的就赋值给这个 |  
| Key  | 加密方法用到的密钥 |


## 方法列表

| 方法 | 返回值 | 说明 |  
| :---: | :--- | :--- |  
| AesCBC()  | []byte, error | aes cbc加密 |  
| AesCBCDecode() | []byte, error | aes cbc解密 |  
| AesCBCDecodeBase64()| []byte, error | 高级加密标准AES解密Base64的结果 |  
| AesECB() | []byte, error | 高级加密标准AES ECB加密 |  
| AesECBDecode() | []byte, error | 高级加密标准AES ECB解密 |  
| AesECBDecodeBase64() | []byte, error | 高级加密标准AES解密ECB Base64的结果 |  
| Md5() | string | md5加密 |  
| Md5Str(str string) | string | md5加密,与上面区别在于这里直接传入要加密的字符串，上面需要给encry.Data赋值后调用 |  
| Rand() | string | 生成一个随机数字字符串 |  
| HmacSha256() | string | HMAC SHA256加密 |  
| HmacSha512() | string | HMAC SHA512加密 |  
| HmacSha1() | string | HMAC SHA1加密 |  
| Sha256() | string | SHA256加密 |  
| Sha512() | string | SHA512加密 |  
| Base64Encode() | string | base64加密 |  
| Base64Decode() | string | base64解密 |  
| PublicBlock(rsakeys *RsaKeys)  | []byte | 获取公钥 |  
| PrivateBlock(rsakeys *RsaKeys)  | []byte | 获取私钥 |  
| Rsa(rsakeys *RsaKeys) | []byte, error | rsa加密 |  
| RsaDecode(rsakeys *RsaKeys) | []byte, error | rsa解密 |  
| RsaKeys.GenerateRSAKey() | string | 生成RSA私钥和公钥，保存到文件中,注意该方法属于RsaKeys结构体 | 

## Rsa结构体
```
type RsaKeys struct {
    Bits         int    
    KeyDir       string 
    PrivateBlock []byte 
    PrivateFile  string 
    PublicBlock  []byte 
    PublicFile   string 
}

```
## Rsa结构体主要属性介绍

| 属性名称  | 说明 |  
| :---: | :--- |  
| Bits  | 加密字长，一般赋值1024，2048 |  
| KeyDir  | 放置加密文件的路径，这个目录下文件名固定public.pem，private.pem，要使用rsa算法要么这个目录下已经存在这两个文件，要么先使用GenerateRSAKey方法生成这两个文件 |  
| PrivateBlock  | 生成的私钥 |  
| PrivateFile  | 生成的私钥文件名 |
| PublicBlock  | 生成的公钥 |
| PublicFile  | 生成的公钥文件名 |

