# 日志消息工具
## 说明
    方便统一消息展示用的
## 使用方法
```golang
logs := huomao_log.GetLogClass()
logs.Info("提示消息")
```

## 方法介绍

| 方法 | 返回值 | 说明 |  
| :---: | :--- | :--- |  
| Info(message interface{})  |  | 输出信息类消息提示 |  
| Error(message interface{})  |  | 输出错误类消息，并结束程序 |  
| GinJson(code int, msg string, data interface{})   | map[string]interface{} | 为了统一gin框架的返回结果，统一返回格式：gin.H{"code": code, "msg": msg, "data": data} |  
| FInfoLog(file string, message string)   | | 将日志写入文件 |  


