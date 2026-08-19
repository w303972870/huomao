# 工具库
## 说明
    汇总了一些日常经常用到的工具方法
## 使用方法
```golang
tools := huomao_tool.GetToolClass()
_ , t := tools.NtpTime()
```
参考 [工具集示例](../example/tool "工具集示例")

## 方法列表

| 方法 | 返回值 | 说明 |  
| :---: | :--- | :--- |  
| NtpTime() | error, time.Time | 使用NTP服务器time.windows.com:123，获取NTP时间 |  
| Exist(path string) | bool, error | 判断文件或文件夹是否存在，true：存在，false：不存在，err不为空说明有错 |  
| FileHash(fname string)  | string | 获取文件的hash值，相当于md5(file) |  
| MemOfPro()  | uint64 | 获取程序当前占用的内存 |  
| MemOfParm(p interface{})  | uintptr | 获取变量占用的内存 |  
| StrInSlice(a string, list []string)  | bool | 查找一个字符串是否在一个切片中，true：存在，false：不存在 |  
| Str(parms ...string)  | string | 拼接字符串 |  
| Bye(code int) |  | 退出程序 |  
| BinName()  | string | 获取当前cgi bin名称 |  
| CurrentPath()  | string | 返回当前执行程序目录 |  
| MkDir(tdir string) |  | 创建目录 |  
| Type(parm interface{})  | string | 获取变量是什么类型 |  
| IsIp(ip string)   | bool | 判断是不是ip格式 |  
| Touch(file string) | bool, *os.File | 创建文件 |  
| IsCharAndNum(str string)  | bool | 验证只允许是字母和数字 |  
| IsEmail(email string)  | bool | 验证邮箱地址格式 |  
| IsMobile(mobileNum string) | bool | 验证手机号格式 |  

