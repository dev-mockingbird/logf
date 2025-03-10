package logf

import "encoding/json"

func JSON(data any) string {
	bs, _ := json.Marshal(data)
	return string(bs)
}
