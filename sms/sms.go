package sms

import (
	"encoding/json"
	"fmt"
	"hust-pass/config"

	"github.com/imroc/req"
)

type Client struct {
	AppID      string
	AppSecret  string
	TemplateID string
}

func InitClient() *Client {
	sender := Client{
		AppID:      config.GlobalConfig.SMSAppID,
		AppSecret:  config.GlobalConfig.SMSAppSecret,
		TemplateID: config.GlobalConfig.SMSTemplateID,
	}
	return &sender
}

func (c *Client) Send(phoneNumber string, templateParams []string) error {
	bytes, err := json.Marshal(&templateParams)
	if err != nil {
		return err
	}
	r := req.New()
	r.EnableInsecureTLS(true)
	postParam := req.Param{
		"appId":          c.AppID,
		"appSecret":      c.AppSecret,
		"templateId":     c.TemplateID,
		"number":         phoneNumber,
		"templateParams": string(bytes),
	}
	sendResp, err := r.Post(
		"https://sms_developer.zhenzikj.com/sms/v2/send.do",
		req.Header{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		postParam,
	)
	if err != nil {
		return err
	}
	var sendRespBody SendRespBody
	err = sendResp.ToJSON(&sendRespBody)
	if err != nil {
		return err
	}
	if sendRespBody.Code != 0 {
		return fmt.Errorf("发送短信失败:%s", sendRespBody.Data)
	}
	fmt.Println("发送短信成功")
	return nil
}

type SendRespBody struct {
	Code int64  `json:"code"`
	Data string `json:"data"`
}
