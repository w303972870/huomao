package huomao_log

import (
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
	"log"
	"os"
	"runtime"
)

type htLogs struct{}

func (mt *htLogs) Info(message interface{}) {
	color.Info.Println(message)
}

func (mt *htLogs) Error(message interface{}) {
	color.Error.Println(message)
	os.Exit(1)
}

/*将接口的返回统一格式化*/
func (mt *htLogs) GinJson(code int, msg string, data interface{}) map[string]interface{} {
	return gin.H{"code": code, "msg": msg, "data": data}
}

func (mt *htLogs) FInfoLog(file string, message string) {
	logFile, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if nil != err {
		color.Error.Println(err.Error())
	}
	loger := log.New(logFile, "", log.Ldate|log.Ltime)
	loger.SetFlags(log.Ldate | log.Ltime)
	loger.Println("[", runtime.NumGoroutine(), "]", message)
	if err := logFile.Close(); err != nil {
		color.Error.Println(err)
	}
}
