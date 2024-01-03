package duel

import "strings"

const (
	HELP = "/duel daily：每日一题（纯随机）\n/duel problem <rating> <tags>：从题库中找出 rating 难度，tag 是 tags 的题（tags 可以是多个，空格请用下划线”_“连接，多个 tag 用空格连接,rating 是必要的，tags 可以留空，加上”!“英文可以不选择该 tag）\n/duel finish：结算每日一题，解绑会导致积分清零（以后可能会持久化）\n/duel rank：前十名分数最高的\n/duel point：查询自己现在的分数"
)

func GetContent(content, id string) string {
	cmd := strings.Fields(content)
	length := len(cmd)
	if length == 0 {
		return HELP
	}
	if cmd[0] == "daily" {
		s, _ := GetDailyProblem()
		return s
	}
	if cmd[0] == "problem" {
		if len(cmd) < 2 {
			return "请输入题目 rating"
		}
		if len(cmd) >= 3 {
			return GetProblem(cmd[1], cmd[2:])
		} else {
			return GetProblem(cmd[1], make([]string, 0))
		}
	}
	if cmd[0] == "finish" {
		return addPiont(id)
	}
	if cmd[0] == "rank" {
		return getMax10()
	}
	if cmd[0] == "point" {
		return checkPoint(id)
	}
	return HELP
}
