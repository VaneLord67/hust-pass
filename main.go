package main

import (
	"fmt"
	"hust-pass/config"
	"hust-pass/sms"
	"hust-pass/spider"
	"log"
	"strconv"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/tebeka/selenium"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	RetryTask()
	log.Println("开启定时任务...")
	crontab := cron.New()
	_, err = crontab.AddFunc("0 0,8,16 * * *", RetryTask)
	if err != nil {
		log.Fatal(err)
		return
	}
	// 启动定时器 定时任务是另起协程执行的
	crontab.Start()
	select {}
}

func ElecTask() error {
	now := time.Now()
	//在后台启动一个ChromeDriver实例
	service, err := selenium.NewChromeDriverService(config.GlobalConfig.ChromeDriverPath,
		config.GlobalConfig.ChromeDriverServicePort, []selenium.ServiceOption{}...)
	if err != nil {
		log.Fatal(err)
	}
	defer func(service *selenium.Service) {
		err := service.Stop()
		if err != nil {
			log.Fatal(err)
		}
	}(service)
	wd, err := spider.InitWebDriver()
	if err != nil {
		log.Fatal(err)
	}
	defer func(wd selenium.WebDriver) {
		err := wd.Quit()
		if err != nil {
			log.Fatal(err)
		}
	}(wd)
	elec, err := spider.LoginGetElec(wd)
	if err != nil {
		log.Println("出错:", err)
		return err
	}
	fmt.Println("电费:", elec)
	fmt.Println("运行时长:", time.Now().Sub(now).Seconds())
	elecNumber, err := strconv.ParseFloat(elec, 64)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if elecNumber < config.GlobalConfig.ElecThreshold {
		client := sms.InitClient()
		err = client.Send(config.GlobalConfig.PhoneNumber,
			[]string{config.GlobalConfig.RoomID, elec})
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func RetryTask() {
	for maxTryCnt := 5; maxTryCnt > 0; maxTryCnt-- {
		err := ElecTask()
		if err == nil {
			break
		}
		fmt.Println("出错:", err.Error())
		fmt.Println("重试...")
	}
}
