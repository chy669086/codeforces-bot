package report

import (
	"bytes"
	"fmt"
)

var tranInt map[string]string = map[string]string{
	"January":   "1",
	"February":  "2",
	"March":     "3",
	"April":     "4",
	"May":       "5",
	"June":      "6",
	"July":      "7",
	"August":    "8",
	"September": "9",
	"October":   "10",
	"November":  "11",
	"December":  "12",
}

var tran map[string]string = map[string]string{
	"January":   " 1 月",
	"February":  " 2 月",
	"March":     " 3 月",
	"April":     " 4 月",
	"May":       " 5 月",
	"June":      " 6 月",
	"July":      " 7 月",
	"August":    " 8 月",
	"September": " 9 月",
	"October":   " 10 月",
	"November":  " 11 月",
	"December":  " 12 月",
}

type report struct {
	date    string
	max     int
	maxPro  string
	cnt     int
	passCnt int
}

func (r *report) String() string {
	if r.cnt == 0 {
		return "你" + r.date + "好像没有提交过题目"
	}
	var tmp bytes.Buffer
	fmt.Fprintf(&tmp, "你%s提交了 %d 次，正确率为 %.2f，期中最难的题目是 %s，难度为 %d。", r.date, r.cnt, float64(r.passCnt)/float64(r.cnt), r.maxPro, r.max)
	return tmp.String()
}
