package duel

import (
	"codeforces-bot/src/codeforeces"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

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
		if err == nil {
			log.Fatalf("Problem list update successfully.")
		} else {
			log.Fatalf("Problem list update failed.")
		}
	}
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

	var result = struct {
		Time     time.Time
		Problems []*codeforeces.Problem
	}{
		time.Now(),
		pro.Result.Problems,
	}
	defer file.Close()
	json.NewEncoder(file).Encode(result)
	return nil
}
