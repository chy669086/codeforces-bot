package duel

import (
	"bytes"
	"codeforces-bot/src/bind"
	"codeforces-bot/src/codeforeces"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func add(id string, point int) int {
	f, _ := bind.GetUserList()
	for i := 0; i < len(f.Users); i++ {
		if f.Users[i].Id == id && f.Users[i].Time.Day() != time.Now().Day() {
			f.Users[i].Point += point
			f.Users[i].Time = time.Now()
			file, _ := os.Create("src/global/Users.json")
			defer file.Close()
			json.NewEncoder(file).Encode(f)
			return f.Users[i].Point
		}
	}
	return -1
}

func addPiont(id string) string {
	s := bind.GetUserName(id)
	if s == "" {
		return "你没有绑定账号"
	}
	result, err := codeforeces.GetStatus(s)
	if err != nil {
		return err.Error()
	}
	_, problem := GetDailyProblem()
	for _, x := range result.Result {
		if x.Problem.Name == problem.Name && x.Verdict == "OK" {
			p := add(id, problem.Rating)
			if p != -1 {
				return "你完成了今天的每日一题，获得 " + strconv.Itoa(problem.Rating) + "，现在你的总分是 " + strconv.Itoa(p)
			} else {
				return "你今天已经完成了每日一题！"
			}
		}
		if time.Unix(x.CreationTimeSeconds, 0).Day() < time.Now().Day() {
			break
		}
	}
	return "你似乎没有完成题目"
}

// 根据tag找题目
func GetProblem(rating string, tags []string) string {
	var result struct {
		Time     time.Time
		Problems []*codeforeces.Problem
	}

	toInt := func(s string) int {
		x := 0
		for _, c := range s {
			x = x*10 + int(c-'0')
		}
		return x
	}
	contain := func(tags *[]string, tag *string) bool {
		for _, x := range *tags {
			if x == *tag {
				return true
			}
		}
		return false
	}

	r := toInt(rating)

	if r < 800 || r > 3500 || r/100*100 != r {
		return "rating 应该是 800~3500 的整百数"
	}

	for i := 0; i < len(tags); i++ {
		tags[i] = strings.Replace(tags[i], "_", " ", 1)
		var s string
		if tags[i][0] == '!' {
			s = tags[i][1:]
		} else {
			s = tags[i]
		}
		if !contain(&TAGS, &s) {
			return tags[i] + " 不在 tag 里面，请检查"
		}
	}

	check()

	file, err := os.ReadFile("src/global/problem.json")
	if err != nil {
		return "获取题目列表失败，请重试"
	}
	err = json.Unmarshal(file, &result)
	if err != nil {
		return "获取题目列表失败，请重试"
	}

	var pro []*codeforeces.Problem
	for _, x := range result.Problems {
		if x.Rating == r {
			pro = append(pro, x)
		}
	}

	for _, tag := range tags {
		var p []*codeforeces.Problem
		var s string
		var f bool
		if tag[0] == '!' {
			s = tag[1:]
			f = true
		} else {
			s = tag
			f = false
		}
		for _, x := range result.Problems {
			if contain(&x.Tags, &s) {
				if !f {
					p = append(p, x)
				}
			} else if f {
				p = append(p, x)
			}
		}
		pro = p
	}
	if len(pro) == 0 {
		return "未获取题目"
	}
	var res bytes.Buffer
	problem := pro[rand.Intn(len(pro))]
	fmt.Fprintf(&res, "获取了题目是 %s，来自 %d 的 %s 题。", problem.Name, problem.ContestId, problem.Index)
	return res.String()
}

// 检查题库是否应该更新
func check() {
	var proTime struct {
		Time time.Time
	}

	file, err := os.ReadFile("src/global/problem.json")
	if err != nil {
		return
	}

	err = json.Unmarshal(file, &proTime)
	if err != nil {
		return
	}

	now := time.Now().Unix()
	lst := proTime.Time.Unix()

	if now-lst >= 86400 {
		err = FetchProblems()
		if err != nil {
			log.Fatalf("Problem list update failed.")
		}
	}
}

// 每日一题
func GetDailyProblem() (string, *codeforeces.Problem) {
	var result struct {
		Time     time.Time
		Problems []*codeforeces.Problem
	}

	check()

	file, err := os.ReadFile("src/global/problem.json")
	if err != nil {
		return "获取题目列表失败，请重试", nil
	}
	err = json.Unmarshal(file, &result)
	if err != nil {
		return "获取题目列表失败，请重试", nil
	}

	var pro []*codeforeces.Problem
	for _, x := range result.Problems {
		if x.Rating != 0 && x.Rating < 1500 && x.Rating > 10000 {
			pro = append(pro, x)
		}
	}

	dateF := func(day time.Time) int64 {
		res := day.Unix()
		res -= int64(day.Second())
		res -= int64(day.Minute() * 60)
		res -= int64(day.Hour() * 3600)
		return res
	}
	length := len(pro)
	day := dateF(time.Now())
	problem := pro[int(day%int64(length))]
	q := time.Now().Format("2006-01-02") + " 的每日一题是编号为 " + strconv.Itoa(problem.ContestId) + " 的比赛的 " + problem.Name + " 题，编号是 " + problem.Index
	return q, problem
}

// 取出当前的cf的题目列表
func FetchProblems() error {
	pro, err := codeforeces.GetProblem()
	if err != nil {
		return err
	}
	if pro.Status != "OK" {
		return fmt.Errorf("%s", pro.Status)
	}
	file, e := os.Create("src/global/problem.json")
	if e != nil {
		return e
	}

	dateF := func(day time.Time) time.Time {
		res := day.Unix()
		res -= int64(day.Second())
		res -= int64(day.Minute() * 60)
		res -= int64(day.Hour() * 3600)
		return time.Unix(res, 0)
	}

	var result = struct {
		Time     time.Time
		Problems []*codeforeces.Problem
	}{
		dateF(time.Now()),
		pro.Result.Problems,
	}
	defer file.Close()
	json.NewEncoder(file).Encode(result)
	return nil
}
