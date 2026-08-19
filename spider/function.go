package huomao_spider

import (
	"encoding/json"
)

func ParseJsonString(str string) [][][]map[string][]string {
	var parse_data [][][]map[string][]string
	err := json.Unmarshal([]byte(str), &parse_data)
	if err != nil || is_not_null(parse_data) == true {
		return nil
	}
	return parse_data
}

func is_not_null(action_list [][][]map[string][]string) bool {
	is_null := true

	for _, page := range action_list {
		for _, list := range page {
			if len(list) > 0 {
				for _, item := range list {
					if len(item) > 0 {
						_, ok := item["Click"]
						_, ok1 := item["SetValue"]
						_, ok2 := item["Captcha"]
						_, ok3 := item["Wait"]
						_, ok4 := item["PageNext"]
						if !ok && !ok1 && !ok2 && !ok3 && !ok4 {
							continue
						}
						if ok && len(item["Click"]) > 0 {
							is_null = false
						}
						if ok1 && len(item["SetValue"]) > 1 {
							is_null = false
						}
						if ok2 && len(item["Captcha"]) > 2 {
							is_null = false
						}
						if ok3 && len(item["Wait"]) > 0 {
							is_null = false
						}
						if ok4 && len(item["PageNext"]) > 1 {
							is_null = false
						}

					}
				}
			}
		}
	}
	return is_null
}
