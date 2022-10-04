package spider

import (
	"fmt"
	"github.com/tebeka/selenium/chrome"
	"hust-pass/config"
	"hust-pass/ocr"
	"log"
	"time"

	"github.com/tebeka/selenium"
)

func RetryLogin(wd selenium.WebDriver) error {
	for maxRetryCnt := 3; maxRetryCnt > 0; maxRetryCnt-- {
		err := Login(wd)
		if err != nil {
			continue
		}
		err = wd.Get("http://one.hust.edu.cn")
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("登录失败次数过多")
}

func Login(wd selenium.WebDriver) error {
	ocrResult := ""
	for maxRetryCnt := 5; len(ocrResult) != 4 && maxRetryCnt > 0; maxRetryCnt-- {
		jsid, ipPool, err := GetJSIDAndIPPool(wd)
		ocrResult, err = ocr.OCR(jsid, ipPool)
		if err != nil {
			return err
		}
	}
	if len(ocrResult) != 4 {
		log.Fatal("OCR出错次数过多")
	}
	var usernameElement selenium.WebElement
	err := wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		element, err := wd.FindElement(selenium.ByCSSSelector, "#un")
		if err == nil {
			usernameElement = element
			return true, nil
		}
		return false, nil
	}, time.Second*30)
	if err != nil {
		return err
	}
	err = usernameElement.SendKeys(config.GlobalConfig.Username)
	if err != nil {
		return err
	}
	passwordInput, err := wd.FindElement(selenium.ByCSSSelector, "#pd")
	if err != nil {
		return err
	}
	err = passwordInput.SendKeys(config.GlobalConfig.Password)
	if err != nil {
		return err
	}
	codeInput, err := wd.FindElement(selenium.ByID, "code")
	if err != nil {
		return err
	}
	err = codeInput.SendKeys(ocrResult)
	if err != nil {
		return err
	}
	loginBtn, err := wd.FindElement(selenium.ByID, "index_login_btn")
	if err != nil {
		return err
	}
	err = loginBtn.Click()
	if err != nil {
		return err
	}
	return nil
}

func GetJSIDAndIPPool(wd selenium.WebDriver) (string, string, error) {
	err := wd.Get("http://pass.hust.edu.cn/cas/login")
	if err != nil {
		return "", "", err
	}
	cookies, err := wd.GetCookies()
	if err != nil {
		return "", "", err
	}
	var ipPool string
	var jsid string
	for _, cookie := range cookies {
		switch cookie.Name {
		case "JSESSIONID":
			jsid = cookie.Value
		case "BIGipServerpool-icdc-cas2":
			ipPool = cookie.Value
		}
	}
	return jsid, ipPool, nil
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
