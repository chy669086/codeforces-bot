package main

import (
	"codeforces-bot/src/bind"
	"codeforces-bot/src/duel"
	"codeforces-bot/src/report"
	"context"
	"log"
	"math/rand"
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

const (
	HeartBeat = "/heartbeat"
	Duel      = "/duel"
	Bind      = "/bind"
	Report    = "/report"
)

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

var fafeng = []string{
	"你骄傲吗？你有什么骄傲的？女人，我已经开始上第三节课了，真正的栋梁之材，是每周五天课，五天的早八，我以为我在杀了十年的鱼，心已经和我的杀鱼刀一样冷了，但在点开课表后，我冰冷的泪划过我的脸庞，原来这世间还有可以伤害",
	"猛地推门而入！推门而入！推门而入！推门而入！推门而入！推门而入！推门而入！推门而入！推门而入！推门而入！推门而入！推门而入！推门而入！推门而入！推门而入！推门而入！推门而入！推门而入！推了99次之后发现推的是",
	"红色是毁灭 蓝色是冷漠 绿色是伪装 白色是虚无 粉色是虚伪 紫色是神秘 橙色是愤怒 黑色是归宿 黄色发给我",
	"地球没照我 样转？硬撑罢了！地球我没 照样转？硬撑罢了！地球没我照样转？硬罢撑 了！地球没我样照 转？硬撑罢了！地没球 我照样转？硬撑罢了！地球没我照样转？硬撑罢了！",
	"撒了一地\n▼▲►▼▲►▼▼▲▼▲►▼▲►▼▲►▼▲►▼▲►▼▼▲►▼►▼▲►▼▼▲►▼▲►▼▲►▼▼►▼▲►▼▲►▼▼▲▼▲►▼▲►▼►▼▲►▼▲►▼▲▲►▼▲►▼►▼▲▼▲►▼▲▲►▼►▼\n呜呜呜赔我的妙脆角，赔我的妙脆角。",
	"我起床了 这个点起床的人 是未来之星 是国家栋梁 是都市小说里的商业大鳄 是的自律者 是相亲节目里的心动嘉宾 是自然界的 是世间所有丑与恶的唾弃者 是世间所以美与好的创造者",
	"如果不回信息会使你愉快的话你就不要回了 我只是一个渺小的存在 不会让你注意到我 即使我发再多信息也没用 得不到始终就是得不到 我累了",
	"没错你是个明白人我明白你明白的意思我也是明白人明白人就应该明白我明白你明白的意思只要大家都明白明白人应明白我明白你明白的意思这样网络环境就是充满明白人明白其他明白人所明白的事",
	"确实 前期貂蝉确实打不过张良，张良一二技能控制加伤害清线很快 但我建议貂蝉第一件出吸血书因为这样可以随时补充状态 而是真的非常好吃 特别是泡了汤汁的油条和方便面 加上麻酱真的是人间美味",
}

// atMessageEventHandler 处理 @机器人 的消息
func atMessageEventHandler(event *dto.WSPayload, data *dto.WSATMessageData) error {
	res := message.ParseCommand(data.Content) //去掉@结构和清除前后空格
	log.Println("cmd = " + res.Cmd + " content = " + res.Content)
	cmd := res.Cmd
	content := res.Content

	switch cmd {
	case HeartBeat:
		api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID, Content: "我的心脏还在跳动呢！"})
	case Duel:
		rp := duel.GetContent(content, data.Author.ID)
		api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID, Content: rp})
	case Bind:
		id := data.Author.ID
		msg := bind.Exe(id, res.Content)
		api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID, Content: msg})
	case Report:
		api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID, Content: report.GetReport(content, data.Author.ID)})
	case "贴贴":
		api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID, Content: "<@" + data.Author.ID + "> 贴贴"})
	default:
		api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID, Content: "<@" + data.Author.ID + ">" + fafeng[rand.Intn(len(fafeng))]})
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
