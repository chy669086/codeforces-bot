package bind

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type UserList struct {
	Users []*User
}

type User struct {
	Id      string
	Account string
	Point   int
	Time    time.Time
}

func Exe(id, content string) string {
	if content == "" {
		return HELP
	}
	cmd := strings.Fields(content)
	if len(cmd) > 2 {
		return HELP
	}
	if cmd[0] == "del" {
		if err := DelUser(id); err != nil {
			return err.Error()
		}
		return "解除绑定成功"
	} else if cmd[0] == "account" {
		if len(cmd) == 1 {
			return "请输入账号"
		}
		if err := BindUser(id, cmd[1]); err != nil {
			return err.Error()
		}
		return "绑定用户 " + cmd[1] + " 成功"
	} else if cmd[0] == "check" {
		s := GetUserName(id)
		if s != "" {
			return "你绑定的账号是 " + s
		} else {
			return "你好像没有绑定账号"
		}
	}
	return HELP
}

func GetUserName(id string) string {
	f, _ := GetUserList()
	for _, user := range f.Users {
		if user.Id == id {
			return user.Account
		}
	}
	return ""
}

func DelUser(id string) error {
	f, _ := GetUserList()
	var tmp []*User
	var i int = -1
	for j, user := range f.Users {
		if user.Id == id {
			i = j
		} else {
			tmp = append(tmp, user)
		}
	}
	if i == -1 {
		return fmt.Errorf("你似乎没有绑定账号！")
	}

	file, _ := os.Create("src/global/Users.json")
	defer file.Close()

	json.NewEncoder(file).Encode(UserList{
		Users: tmp,
	})
	return nil
}

func BindUser(id, account string) error {
	if checkUser(id) {
		return fmt.Errorf("你已经绑定过账号了！")
	}

	AddUser(id, account)
	return nil
}

func checkUser(id string) bool {
	f, _ := GetUserList()
	for _, user := range f.Users {
		if user.Id == id {
			return true
		}
	}
	return false
}

func AddUser(id, account string) error {
	f, err := GetUserList()
	if err != nil {
		return err
	}
	var tmp = User{
		id,
		account,
		0,
		time.Unix(0, 0),
	}
	f.Users = append(f.Users, &tmp)
	fmt.Println(tmp)

	file, err := os.OpenFile("src/global/Users.json", os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	defer file.Close()
	json.NewEncoder(file).Encode(f)
	return nil
}

func GetUserList() (*UserList, error) {
	file, err := os.ReadFile("src/global/Users.json")
	if err != nil {
		return nil, err
	}
	var result UserList
	err = json.Unmarshal(file, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
