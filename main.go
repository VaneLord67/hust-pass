package main

import (
	"github.com/robfig/cron/v3"
	"hust-pass/config"
	"hust-pass/elec"
	"log"
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

func RetryTask() {
	for maxTryCnt := 5; maxTryCnt > 0; maxTryCnt-- {
		err := elec.Task()
		if err == nil {
			break
		}
		log.Println("出错:", err.Error())
		log.Println("重试...")
	}
}
