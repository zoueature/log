package notify

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// https://oapi.dingtalk.com/robot/send?access_token=7389ed52ebd684b9c194237ff718161311630526a1e1631c3362084bae7ee8fe

const dingtalkNotifyURL = "https://oapi.dingtalk.com/robot/send"

type dingTalkRobotClient struct {
	apiURL string
}

func NewDingtalkNotifyClient(accessToken string) Notifier {
	return &dingTalkRobotClient{apiURL: dingtalkNotifyURL + "?access_token=" + accessToken}
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

const msgtype_markdown = "markdown"
const okcode = 0

func (cli *dingTalkRobotClient) SendMarkdown(title, content string, atMobile ...string) error {
	httpClient := http.Client{}
	msg := markdownMsgBody{}
	msg.Msgtype = msgtype_markdown
	msg.Markdown.Title = title
	msg.Markdown.Text = content
	msg.At.AtMobiles = atMobile
	jsonStr, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(jsonStr)
	req, err := http.NewRequest(http.MethodPost, cli.apiURL, reader)
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
	if respData.Errcode != okcode {
		return errors.New(respData.Errmsg)
	}
	return nil
}
