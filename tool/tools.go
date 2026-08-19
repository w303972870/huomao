package huomao_tool

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"time"
	"unsafe"
)

/*获取NTP时间*/
func (t *function) NtpTime() (error, time.Time) {
	conn, err := net.Dial("udp", "time.windows.com:123")
	if err != nil {
		return err, time.Now()
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(15 * time.Second)); err != nil {
		return err, time.Now()
	}

	req := &packet{Settings: 0x1B}

	if err := binary.Write(conn, binary.BigEndian, req); err != nil {
		return err, time.Now()
	}
	rsp := &packet{}
	if err := binary.Read(conn, binary.BigEndian, rsp); err != nil {
		return err, time.Now()
	}
	secs := float64(rsp.TxTimeSec) - ntpEpochOffset
	nanos := (int64(rsp.TxTimeFrac) * 1e9) >> 32

	return nil, time.Unix(int64(secs), nanos)
}

/*判断文件或文件夹是否存在，true：存在，false：不存在，err不为空说明有错*/
func (t *function) Exist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/*文件hash*/
func (t *function) FileHash(fname string) string {
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
		return ""
	}
	defer f.Close()

	br := bufio.NewReader(f)

	h := sha1.New()
	_, err = io.Copy(h, br)

	if err != nil {
		panic(err)
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

/*获取程序占用的内存*/
func (t *function) MemOfPro() uint64 {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return mem.Alloc
}

/*获取变量占用的内存*/
func (t *function) MemOfParm(p interface{}) uintptr {
	return unsafe.Sizeof(p)
}

/*查找一个字符串是否在一个切片中，true：存在，false：不存在*/
func (t *function) StrInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

/*拼接字符串*/
func (t *function) Str(parms ...string) string {
	var buffer bytes.Buffer

	for _, parm := range parms {
		buffer.WriteString(parm)
	}
	return buffer.String()
}

/*退出程序*/
func (t *function) Bye(code int) {
	os.Exit(code)
}

/*获取当前cgi bin名称*/
func (t *function) BinName() string {
	_, file := filepath.Split(os.Args[0])
	return file
}

/*返回当前执行目录*/
func (t *function) CurrentPath() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

/*创建目录*/
func (t *function) MkDir(tdir string) {
	if is, _ := t.Exist(tdir); !is {
		os.MkdirAll(tdir, os.ModePerm)
	}
}

/*获取变量类型*/
func (t *function) Type(parm interface{}) string {
	switch t := parm.(type) {
	case bool:
		return "bool"
	case int:
		return "int"
	default:
		return fmt.Sprintf("%T", t)
	}
}

/*判断ip格式*/
func (t *function) IsIp(ip string) bool {
	if addr := net.ParseIP(ip); addr == nil {
		return false
	} else {
		return true
	}
}

/*创建文件*/
func (t *function) Touch(file string) (bool, *os.File) {
	f, err := os.Create(file)
	if err != nil {
		return false, nil
	}
	return true, f
}

/*验证只允许是字母和数字*/
func (t *function) IsCharAndNum(str string) bool {
	pattern := `^[A-Za-z0-9]+$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(str)
}

/*验证邮箱地址格式*/
func (t *function) IsEmail(email string) bool {
	//pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	//pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	pattern := "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

/*验证手机号格式*/
func (t *function) IsMobile(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}
