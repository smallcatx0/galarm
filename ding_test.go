package galarm_test

import (
	"galarm"
	"log"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var C *viper.Viper

func ConfInit() {
	C = viper.New()
	C.SetConfigFile("./unittest_env.yaml")
	C.SetConfigType("toml")
	err := C.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	if C.GetString("v") != "1.0" {
		log.Fatal("读取配置文件失败")
	}
}

func init() {
	ConfInit()
	galarm.InitDingAla(C.GetString("ding_hook"), C.GetString("ding_secret"))
}

func TestTextMsg(t *testing.T) {
	assert := assert.New(t)

	ding := galarm.DingAlarmNew(C.GetString("ding_hook"), C.GetString("ding_secret"))
	err := ding.Text("测试普通消息", "多行文本内容", "自定义消息体").AtPhones("18681636749").Send()
	ding.Text("消息粘滞").Send()
	assert.NoError(err)
}

func TestMDMsg(t *testing.T) {
	assert := assert.New(t)

	err := galarm.DingAla.Markdown("markdown 消息标题",
		"### 三级标题",
		"> 引用",
		"内容",
		"![screenshot](https://dummyimage.com/600x400)",
	).Send()
	assert.NoError(err)
}

func TestAtPhone(t *testing.T) {
	assert := assert.New(t)

	err := galarm.DingAla.Text("测试消息。At all").AtAll().Send()
	assert.NoError(err)
	err = galarm.DingAla.Text("测试消息", "at 手机号").AtPhones("19911856236").Send()
	assert.NoError(err)
}

func TestSingleBtu(t *testing.T) {
	assert := assert.New(t)
	err := galarm.DingAla.Action("卡片消息",
		"![screenshot](https://dummyimage.com/600x400)",
		"# 一级标题",
		"## 二级标题",
		"### 三级标题").
		SetButs(true, galarm.DingBtn{Title: "查看更多", URL: "http://www.baidu.com"}).
		Send()
	assert.NoError(err)
}

func TestMultiBtu(t *testing.T) {
	assert := assert.New(t)
	err := galarm.DingAla.Action("卡片消息，多按钮",
		"![screenshot](https://dummyimage.com/600x400)",
		"# 一级标题 ",
		"## 二级标题 ",
		"### 三级标题").
		SetButs(false,
			galarm.DingBtn{Title: "查看更多", URL: "http://www.baidu.com"},
			galarm.DingBtn{Title: "不再提醒", URL: "http://www.bilibili.com"}).
		Send()
	assert.NoError(err)
}

func TestFeedCards(t *testing.T) {
	assert := assert.New(t)
	err := galarm.DingAla.FeedCard(
		galarm.DingFeedCard{
			Title:  "标题一一一一一一一",
			PicURL: "https://dummyimage.com/600x400",
			MsgURL: "http://www.baidu.com"},
		galarm.DingFeedCard{
			Title:  "标题二二二二二二",
			PicURL: "https://dummyimage.com/400x400",
			MsgURL: "http://www.bilibili.com"},
		galarm.DingFeedCard{
			Title:  "标题三三三三三三三三",
			PicURL: "https://dummyimage.com/300x300",
			MsgURL: "http://www.bilibili.com"},
	).Send()
	assert.NoError(err)
}

func TestSendText(t *testing.T) {
	assert := assert.New(t)
	err := galarm.DingAla.SendText(
		"这是快捷普通消息",
		"1. 啊手动阀手动阀手动阀",
		"2. 氨基酸地方阿斯顿饭卡手动阀离开",
	)
	assert.NoError(err)
}
