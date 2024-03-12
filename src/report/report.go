package report

import (
	"bytes"
	"codeforces-bot/src/bind"
	"codeforces-bot/src/codeforeces"
	"fmt"
	"strings"
	"time"
)

func GetReport(content, id string) string {
	cmd := strings.Fields(content)
	if len(cmd) == 0 {
		return HELP
	}
	if cmd[0] == "daily" {
		return GetDailyReport(id)
	}
	if cmd[0] == "month" {
		return GetMonthsReport(id)
	}
	if cmd[0] == "self" {
		return GetContestReport(bind.GetUserName(id))
	} else if cmd[0] == "account" {
		if len(cmd) < 2 {
			return HELP
		}
		return GetContestReport(cmd[1])
	}
	toInt := func(s string) int {
		x := 0
		for _, c := range s {
			x = x*10 + int(c-'0')
		}
		return x
	}
	if cmd[0] == "get" {
		if len(cmd) == 3 {
			if toInt(cmd[1]) < 2020 || toInt(cmd[1]) > time.Now().Year() {
				return "请输入正确的年份"
			}
			if toInt(cmd[2]) < 1 || toInt(cmd[2]) > 12 {
				return "月份应该是 1~12 的整数"
			}
			return GetAssignReport(id, cmd[1], cmd[2])
		}
	}
	return HELP
}

func GetContestReport(name string) string {

	f, err := codeforeces.GetRantingChange(name)
	if err != nil {
		return err.Error()
	}

	var maxRating = 0
	var maxContest string

	for _, result := range f.Result {
		if result.NewRating > maxRating {
			maxRating = result.NewRating
			maxContest = result.ContestName
		}
	}

	var s bytes.Buffer

	fmt.Fprintf(&s, "%s 参加了 %d 场比赛，在 %s 达到最高 rating %d。", name, len(f.Result), maxContest, maxRating)
	return s.String()
}

func GetAssignReport(id, year, month string) string {
	f, _ := bind.GetUserList()
	account := ""

	for _, user := range f.Users {
		if user.Id == id {
			account = user.Account
			break
		}
	}
	if account == "" {
		return "你没有绑定账号！"
	}
	submits, _ := codeforeces.GetStatus(account)
	var rep report
	rep.date = "在 " + year + " 年 " + month + " 月"
	toInt := func(s string) int {
		x := 0
		for _, c := range s {
			x = x*10 + int(c-'0')
		}
		return x
	}
	for _, submit := range submits.Result {
		if time.Unix(submit.CreationTimeSeconds-5*3600, 0).Year() < toInt(year) ||
			(time.Unix(submit.CreationTimeSeconds-5*3600, 0).Year() == toInt(year) &&
				toInt(tranInt[time.Unix(submit.CreationTimeSeconds-5*3600, 0).Month().String()]) < toInt(month)) {
			break
		}
		if time.Unix(submit.CreationTimeSeconds-5*3600, 0).Year() == toInt(year) &&
			toInt(tranInt[time.Unix(submit.CreationTimeSeconds-5*3600, 0).Month().String()]) == toInt(month) {
			rep.cnt++
			if submit.Verdict == "OK" {
				rep.passCnt++
				if submit.Problem.Rating > rep.max {
					rep.max = submit.Problem.Rating
					rep.maxPro = submit.Problem.Name
				}
			}
		}
	}
	return rep.String()
}

func GetMonthsReport(id string) string {
	f, _ := bind.GetUserList()
	account := ""

	for _, user := range f.Users {
		if user.Id == id {
			account = user.Account
			break
		}
	}
	if account == "" {
		return "你没有绑定账号！"
	}
	submits, _ := codeforeces.GetStatus(account)
	var rep report
	rep.date = "在" + tran[time.Now().Month().String()]
	for _, submit := range submits.Result {
		if time.Unix(submit.CreationTimeSeconds-5*3600, 0).Month() != time.Now().Month() {
			break
		}
		rep.cnt++
		if submit.Verdict == "OK" {
			rep.passCnt++
			if submit.Problem.Rating > rep.max {
				rep.max = submit.Problem.Rating
				rep.maxPro = submit.Problem.Name
			}
		}
	}
	return rep.String()
}

func GetDailyReport(id string) string {
	f, _ := bind.GetUserList()
	account := ""

	for _, user := range f.Users {
		if user.Id == id {
			account = user.Account
			break
		}
	}
	if account == "" {
		return "你没有绑定账号！"
	}
	submits, _ := codeforeces.GetStatus(account)
	var rep report
	rep.date = "今天"
	for _, submit := range submits.Result {
		if time.Unix(submit.CreationTimeSeconds, 0).Day() != time.Now().Day() {
			break
		}
		rep.cnt++
		if submit.Verdict == "OK" {
			rep.passCnt++
			if submit.Problem.Rating > rep.max {
				rep.max = submit.Problem.Rating
				rep.maxPro = submit.Problem.Name
			}
		}
	}
	return rep.String()
}
