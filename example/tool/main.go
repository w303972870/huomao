package main

import (
	"fmt"
	"github.com/w303972870/huomao/tool"
)

func main() {

	tools := huomao_tool.GetToolClass()

	err, t := tools.NtpTime()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("当前NTP时间：", t.Format("2006-01-02 15:04:05"))
	}

	mem := tools.MemOfPro()
	fmt.Println("内存数：", mem)

	var s1 = []string{"ok", "ok1", "ok2", "ok3", "ok4", "ok5"}
	if isset := tools.StrInSlice("ok", s1); isset == true {
		fmt.Println("isset")
	} else {
		fmt.Println("not in")
	}

}
