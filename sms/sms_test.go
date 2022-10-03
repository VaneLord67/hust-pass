package sms

import (
	"hust-pass/config"
	"testing"
)

func TestSMSSend(t *testing.T) {
	err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	client := InitClient()
	err = client.Send("手机号", []string{"参数1", "参数2"})
	if err != nil {
		t.Fatal(err)
	}
}
