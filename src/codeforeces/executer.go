package codeforeces

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// 获取题目列表
func GetProblem(tag ...string) (*ProblemList, error) {
	q := url.QueryEscape(strings.Join(tag, "&"))
	resp, err := http.Get(ProblemUrl + "?tag=" + q)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result ProblemList
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// 查询一个人的比赛记录
func GetRantingChange(name string) (*RantingResultList, error) {
	qUrl := RatingUrl + "?handle=" + name
	resp, err := http.Get(qUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result RantingResultList
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// 查找一个人某场比赛的赛时提交记录
func GetPassedProblemInContest(name string, contestId int) (*ResultList, error) {
	q := ContestStatusUrl + "?handle=" + name + "&contestId=" + strconv.Itoa(contestId)
	resp, err := http.Get(q)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result ResultList
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// 查找一个人的所有提交记录
func GetStatus(name string) (*ResultList, error) {
	q := StatusUrl + "?handle=" + name
	resp, err := http.Get(q)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result ResultList
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
