package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var GlobalConfig *Config

func LoadConfig() error {
	bytes, readErr := ioutil.ReadFile("./config.json")
	if readErr != nil {
		return fmt.Errorf("没有找到配置文件")
	}
	var cfg Config
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return fmt.Errorf("解析配置文件异常，%s", err.Error())
	}
	GlobalConfig = &cfg
	return nil
}

type Config struct {
	AK                      string  `json:"ak"`
	SK                      string  `json:"sk"`
	Username                string  `json:"username"`
	Password                string  `json:"password"`
	ChromeDriverServicePort int     `json:"chrome_driver_service_port"`
	ChromeDriverPath        string  `json:"chrome_driver_path"`
	CallerToken             string  `json:"caller_token"`
	OCRServicePort          int     `json:"ocr_service_port"`
	AmMeterID               string  `json:"AmMeter_ID"`
	SMSAppID                string  `json:"sms_app_id"`
	SMSAppSecret            string  `json:"sms_app_secret"`
	SMSTemplateID           string  `json:"sms_template_id"`
	PhoneNumber             string  `json:"phone_number"`
	RoomID                  string  `json:"room_id"`
	ElecThreshold           float64 `json:"elec_threshold"`
}
