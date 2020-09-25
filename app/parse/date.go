package parse

import (
	"Meink/app/util"
	"fmt"
	"time"
)

const (
	//基本的时间格式化
	Date_Fromat               = "2006-01-02 15:04:05"
	Date_Format_With_Timezone = "2006-01-02 15:04:05 -0700"
)

func Date(dateStr string) time.Time {
	date, err := time.Parse(fmt.Sprintf(Date_Format_With_Timezone), dateStr)
	if err != nil {
		date, err = time.ParseInLocation(fmt.Sprintf(Date_Fromat), dateStr, time.Now().Location())
		if err != nil {
			util.MError("文章时间解析错误:" + err.Error())
		}
	}
	return date
}
