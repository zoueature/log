package notify

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// https://oapi.dingtalk.com/robot/send?access_token=7389ed52ebd684b9c194237ff718161311630526a1e1631c3362084bae7ee8fe

const dingtalkNotifyURL = "https://oapi.dingtalk.com/robot/send"

type dingTalkRobotClient struct {
	accessToken string
	signSecret  string
}

func NewDingtalkNotifyClient(accessToken string, signSecret string) Notifier {
	return &dingTalkRobotClient{
		accessToken: accessToken,
		signSecret:  signSecret,
	}
}

type markdownMsgBody struct {
	Msgtype  string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
	At struct {
		AtMobiles []string `json:"atMobiles,omitempty"`
		AtUserIds []string `json:"atUserIds,omitempty"`
		IsAtAll   bool     `json:"isAtAll"`
	} `json:"at"`
}

type respMsg struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

const msgTypeMarkdown = "markdown"
const okCode = 0

func (cli *dingTalkRobotClient) sign() (string, string) {
	timeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	s := timeStamp + "\n" + cli.signSecret
	hash := hmac.New(sha256.New, []byte(cli.signSecret))
	hash.Write([]byte(s))
	hashResult := hash.Sum(nil)

	return timeStamp, base64.StdEncoding.EncodeToString(hashResult)
}

func (cli *dingTalkRobotClient) SendMarkdown(title, content string, atMobile ...string) error {
	httpClient := http.Client{}
	msg := markdownMsgBody{}
	msg.Msgtype = msgTypeMarkdown
	msg.Markdown.Title = title
	msg.Markdown.Text = content
	msg.At.AtMobiles = atMobile
	jsonStr, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(jsonStr)
	param := url.Values{}
	timestamp, sign := cli.sign()
	param.Add("timestamp", timestamp)
	param.Add("sign", sign)
	param.Add("access_token", cli.accessToken)
	req, err := http.NewRequest(http.MethodPost, dingtalkNotifyURL+"?"+param.Encode(), reader)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	respData := new(respMsg)
	err = json.Unmarshal(body, respData)
	if err != nil {
		return err
	}
	if respData.Errcode != okCode {
		return errors.New(respData.Errmsg)
	}
	return nil
}
