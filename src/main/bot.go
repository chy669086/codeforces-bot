package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
	yaml "gopkg.in/yaml.v2"
)

// Config 定义了配置文件的结构
type Config struct {
	AppID uint64 `yaml:"appid"` //机器人的appid
	Token string `yaml:"token"` //机器人的token
}

var config Config
var api openapi.OpenAPI
var ctx context.Context

// 第一步： 获取机器人的配置信息，即机器人的appid和token
func init() {
	content, err := os.ReadFile("src/global/config.yaml")
	if err != nil {
		log.Println("读取配置文件出错， err = ", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		log.Println("解析配置文件出错， err = ", err)
		os.Exit(1)
	}
	log.Println(config)
}

// atMessageEventHandler 处理 @机器人 的消息
func atMessageEventHandler(event *dto.WSPayload, data *dto.WSATMessageData) error {
	res := message.ParseCommand(data.Content) //去掉@结构和清除前后空格
	log.Println("cmd = " + res.Cmd + " content = " + res.Content)
	cmd := res.Cmd
	// content := res.Content
	switch cmd {
	case "/heartbeat":
		api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID, Content: "我的心脏还在跳动呢！"})
	}
	return nil
}

func main() {
	//第二步：生成token，用于校验机器人的身份信息
	token := token.BotToken(config.AppID, config.Token)
	//第三步：获取操作机器人的API对象
	api = botgo.NewOpenAPI(token).WithTimeout(3 * time.Second)
	//获取context
	ctx = context.Background()
	//第四步：获取websocket
	ws, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Fatalln("websocket错误， err = ", err)
		os.Exit(1)
	}

	var atMessage event.ATMessageEventHandler = atMessageEventHandler

	intent := websocket.RegisterHandlers(atMessage)     // 注册socket消息处理
	botgo.NewSessionManager().Start(ws, token, &intent) // 启动socket监听
}
