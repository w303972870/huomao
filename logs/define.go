package huomao_log

var p *htLogs

type iLog interface {
	Info(message interface{})

	Error(message interface{})

	GinJson(code int, msg string, data interface{}) map[string]interface{}

	FInfoLog(file string, message string)
}

func GetLogClass() *htLogs {
	if p == nil {
		p = &htLogs{}
	}
	return p
}
