package notify

import (
	"os"
	"testing"
)

func TestSignReq(t *testing.T) {
	cli := &dingTalkRobotClient{
		accessToken: os.Getenv("access_token"),
		signSecret:  os.Getenv("sign_secret"),
	}
	err := cli.SendMarkdown("asdasdsa", "dsadsadsa")
	if err != nil {
		t.Error(err)
	}
}
