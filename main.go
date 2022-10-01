package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"hust-pass/config"
	"hust-pass/spider"
	"log"
	"time"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}
	now := time.Now()
	//在后台启动一个ChromeDriver实例
	service, err := selenium.NewChromeDriverService(config.GlobalConfig.ChromeDriverPath,
		config.GlobalConfig.ChromeDriverServicePort, []selenium.ServiceOption{}...)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer service.Stop()
	wd, err := InitWebDriver()
	if err != nil {
		log.Fatal(err)
	}
	defer wd.Quit()
	elec, err := spider.LoginGetElec(wd)
	if err != nil {
		log.Println("出错:", err)
		return
	}
	fmt.Println("电费:", elec)
	fmt.Println("运行时长:", time.Now().Sub(now).Seconds())
}

func InitWebDriver() (selenium.WebDriver, error) {
	//连接到本地运行的 WebDriver 实例
	//这里的map键值只能为browserName，源码里需要获取这个键的值，来确定连接的是哪个浏览器
	caps := selenium.Capabilities{"browserName": "chrome"}
	//禁止图片加载，加快渲染速度
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}
	//设置实验谷歌浏览器驱动的参数
	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Args: []string{
			"--headless", //设置Chrome无头模式
		},
	}
	//添加浏览器设置参数
	caps.AddChrome(chromeCaps)
	//NewRemote 创建新的远程客户端，这也将启动一个新会话。 urlPrefix 是 Selenium 服务器的 URL，必须以协议 (http, https, ...) 为前缀。
	//为urlPrefix提供空字符串会导致使用 DefaultURLPrefix,默认访问4444端口，
	//所以最好自定义，避免端口已经被抢占。后面的路由还是照旧DefaultURLPrefix写
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub",
		config.GlobalConfig.ChromeDriverServicePort))
	if err != nil {
		return nil, err
	}
	return wd, nil
}
