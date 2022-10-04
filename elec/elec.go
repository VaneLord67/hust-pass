package elec

import (
	"github.com/tebeka/selenium"
	"hust-pass/config"
	"hust-pass/sms"
	"hust-pass/spider"
	"log"
	"strconv"
	"time"
)

func SpiderGetElecResult(wd selenium.WebDriver) (string, error) {
	log.Println("进入电费查询网站")
	err := wd.Get("http://sdhq.hust.edu.cn/icbs/hust/html/index.html")
	if err != nil {
		return "", err
	}
	var elecResult string
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		elecValueElement, err := wd.FindElement(selenium.ByCSSSelector, ".AmValue")
		if err != nil {
			return false, nil
		}
		elecResult, err = elecValueElement.Text()
		if err != nil {
			return false, nil
		}
		return elecResult != "", nil
	}, time.Second*30)
	if err != nil {
		return "", err
	}
	return elecResult, nil
}

func GetElecResult() (string, error) {
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
	err = spider.RetryLogin(wd)
	if err != nil {
		log.Println("出错:", err)
		return "", err
	}
	elecResult, err := SpiderGetElecResult(wd)
	if err != nil {
		log.Println("出错:", err)
		return "", err
	}
	log.Println("电费:", elecResult)
	log.Println("运行时长:", time.Now().Sub(now).Seconds())
	return elecResult, nil
}

func Task() error {
	elec, err := GetElecResult()
	if err != nil {
		return err
	}
	elecNumber, err := strconv.ParseFloat(elec, 64)
	if err != nil {
		log.Println(err)
		return err
	}
	if elecNumber < config.GlobalConfig.ElecThreshold {
		client := sms.InitClient()
		err = client.Send(config.GlobalConfig.PhoneNumber,
			[]string{config.GlobalConfig.RoomID, elec})
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
