package galarm

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var DingAla *DingAlarm

func InitDingAla(webHook, secret string) {
	DingAla = DingAlarmNew(webHook, secret)
}

type DingAlarm struct {
	webHook   string
	secret    string
	sign      string
	timestamp string
	Msg       *DingMsg
}

type DingMsg struct {
	Msgtype string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text,omitempty"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown,omitempty"`
	At struct {
		AtMobiles []string `json:"atMobiles"`
		AtUserIds []string `json:"atUserIds"`
		IsAtAll   bool     `json:"isAtAll"`
	} `json:"at,omitempty"`
	ActionCard struct {
		Title          string    `json:"title"`
		Text           string    `json:"text"`
		BtnOrientation int       `json:"btnOrientation"`
		SingleTitle    string    `json:"singleTitle"`
		SingleURL      string    `json:"singleURL"`
		Btns           []DingBtn `json:"btns"`
	} `json:"actionCard,omitempty"`
	FeedCard struct {
		Links []DingFeedCard `json:"links"`
	} `json:"feedCard,omitempty"`
}

type DingBtn struct {
	Title string `json:"title"`
	URL   string `json:"actionURL"`
}

type DingFeedCard struct {
	Title  string `json:"title"`
	MsgURL string `json:"messageURL"`
	PicURL string `json:"picURL"`
}

func DingAlarmNew(webHook, secret string) *DingAlarm {
	d := &DingAlarm{
		webHook: webHook,
		secret:  secret,
		Msg:     &DingMsg{},
	}
	return d
}

func (d *DingAlarm) signature() string {
	now := time.Now().Unix() * 1000
	d.timestamp = strconv.FormatInt(now, 10)
	h := hmac.New(sha256.New, []byte(d.secret))
	h.Write([]byte(d.timestamp + "\n" + d.secret))
	sign := base64.URLEncoding.EncodeToString(h.Sum(nil))
	sign = url.PathEscape(sign)
	sign = strings.Replace(sign, "-", "%2B", -1)
	sign = strings.Replace(sign, "_", "%2F", -1)
	d.sign = sign
	return sign
}

func (d *DingAlarm) Send() error {
	err := d.SendMsg(d.Msg)
	// ????????????????????????
	d.Msg = &DingMsg{}
	return err
}

func (d *DingAlarm) SendMsg(msg *DingMsg) error {
	sign := d.signature()
	url := d.webHook + "&timestamp=" + d.timestamp + "&sign=" + sign
	body, _ := json.Marshal(msg)
	resp, err := new(http.Client).Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	res, _ := ioutil.ReadAll(resp.Body)
	ress := make(map[string]interface{})
	json.Unmarshal(res, &ress)
	errcd, ok := ress["errcode"].(float64)
	if ok && errcd == 0 {
		return nil
	}
	return errors.New(string(res))
}

// ????????????
func (d *DingAlarm) Text(con ...string) *DingAlarm {
	d.Msg.Msgtype = "text"
	text := strings.Join(con, "\n")
	d.Msg.Text.Content = text
	return d
}

// markdown??????
func (d *DingAlarm) Markdown(title string, md ...string) *DingAlarm {
	d.Msg.Msgtype = "markdown"
	d.Msg.Markdown.Title = title
	mdStr := strings.Join(md, "\n\n")
	d.Msg.Markdown.Text = mdStr
	return d
}

// ????????????
func (d *DingAlarm) Action(title string, content ...string) *DingAlarm {
	d.Msg.Msgtype = "actionCard"
	d.Msg.ActionCard.Title = title
	text := strings.Join(content, "\n\n")
	d.Msg.ActionCard.Text = text
	return d
}

// ????????????
func (d *DingAlarm) SetButs(butVertical bool, buts ...DingBtn) *DingAlarm {
	if !butVertical {
		d.Msg.ActionCard.BtnOrientation = 1
	}
	if len(buts) <= 1 {
		d.Msg.ActionCard.SingleTitle = buts[0].Title
		d.Msg.ActionCard.SingleURL = buts[0].URL
	} else {
		d.Msg.ActionCard.Btns = buts
	}
	return d
}

// ????????????
func (d *DingAlarm) FeedCard(cards ...DingFeedCard) *DingAlarm {
	d.Msg.Msgtype = "feedCard"
	d.Msg.FeedCard.Links = cards
	return d
}

// at ?????????
func (d *DingAlarm) AtPhones(phone ...string) *DingAlarm {
	d.Msg.At.AtMobiles = phone
	return d
}

// at userID
func (d *DingAlarm) AtUsers(id ...string) *DingAlarm {
	d.Msg.At.AtUserIds = id
	return d
}

// at ????????????
func (d *DingAlarm) AtAll() *DingAlarm {
	d.Msg.At.IsAtAll = true
	return d
}

// ??????markdown??????
func (d *DingAlarm) SendMd(title, content string) error {
	msg := DingMsg{
		Msgtype: "markdown",
	}
	msg.Markdown.Title = title
	msg.Markdown.Text = content
	return d.SendMsg(&msg)
}

// ??????????????????
func (d *DingAlarm) SendText(con ...string) error {
	return d.Text(con...).Send()
}
