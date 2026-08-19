package main

import (
	"fmt"
	"github.com/w303972870/huomao/encry"
)

func main() {
	encry := huomao_encry.GetEncryClass()
	encry.Data = "原始字符串"
	encry.Key = "DFA84B10B7ACDD25"
	aesencode, _ := encry.AesCBC()
	fmt.Println("AesCBC加密结果获取方式一：", aesencode)
	fmt.Println("AesCBC加密结果获取方式二：", encry.EncryData)
	fmt.Println("AesCBC加密结果的base64后的结果：", encry.EncryBase64Data)

	decode_example := huomao_encry.GetEncryClass()
	decode_example.EncryData = string(aesencode)
	decode_example.Key = "DFA84B10B7ACDD25"
	aesdecode, _ := decode_example.AesCBCDecode()
	fmt.Println("AesCBC解密结果获取方式一：", string(aesdecode))
	fmt.Println("AesCBC解密结果获取方式二：", decode_example.Data)

	decode_example.EncryBase64Data = encry.EncryBase64Data
	aesbase64decode, _ := decode_example.AesCBCDecodeBase64()
	fmt.Println("AesCBC Base64方式解密结果1：", string(aesbase64decode))
	fmt.Println("AesCBC Base64方式解密结果2：", decode_example.Data)

	md5_example := huomao_encry.GetEncryClass()
	md5_example.Data = "原始字符串"
	md5str := md5_example.Md5()
	fmt.Println("md5结果：", md5str)
	fmt.Println("md5结果：", md5_example.Md5Str("原始字符串"))

	var RsaKeys huomao_encry.RsaKeys
	rsa_example := huomao_encry.GetEncryClass()
	RsaKeys.KeyDir = "/Users/huomao/keys/"
	RsaKeys.GenerateRSAKey()
	fmt.Println("生成私钥文件：", fmt.Sprint(RsaKeys.KeyDir, "/", RsaKeys.PrivateFile))
	fmt.Println("生成公钥文件：", fmt.Sprint(RsaKeys.KeyDir, "/", RsaKeys.PublicFile))
	fmt.Println("公钥内容：", string(RsaKeys.PublicBlock))
	fmt.Println("私钥内容：", string(RsaKeys.PrivateBlock))

	rsa_example.Data = "火卯"
	rsa_example.Rsa(&RsaKeys)
	fmt.Println("展示rsa加密后的结果：", rsa_example.EncryData)
	fmt.Println("展示rsa加密后的base64结果：", rsa_example.EncryBase64Data)

	rsa_example.RsaDecode(&RsaKeys)
	fmt.Println("展示rsa解密后的结果：", rsa_example.Data)

}
