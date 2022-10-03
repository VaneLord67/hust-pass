package spider

import (
	"fmt"
	"hust-pass/config"
	"hust-pass/ocr"
	"log"
	"time"

	"github.com/tebeka/selenium"
)

func LoginGetElec(wd selenium.WebDriver) (string, error) {
	maxRetryCnt := 3
	ocrResult := ""
	for ; maxRetryCnt > 0; maxRetryCnt-- {
		for maxRetryCnt := 5; len(ocrResult) != 4 && maxRetryCnt > 0; maxRetryCnt-- {
			jsid, ipPool, err := GetJSIDAndIPPool(wd)
			ocrResult, err = ocr.OCR(jsid, ipPool)
			if err != nil {
				return "", err
			}
		}
		if len(ocrResult) != 4 {
			log.Fatal("OCR出错次数过多")
		}
		err := Login(wd, ocrResult)
		if err != nil {
			return "", err
		}
		err = wd.Get("http://sdhq.hust.edu.cn/icbs/hust/html/index.html")
		if err != nil {
			return "", err
		}
		title, err := wd.Title()
		if err != nil {
			return "", err
		}
		if title == "水电收费平台" {
			break
		}
		fmt.Println("登录失败,重新登录...")
		ocrResult = ""
	}
	if maxRetryCnt <= 0 {
		return "", fmt.Errorf("登录失败次数过多")
	}
	log.Println("开启go routine,获取电费信息...")
	ch := make(chan string)
	go func() {
		for {
			elecValueElement, err := wd.FindElement(selenium.ByCSSSelector, ".AmValue")
			if err != nil {
				continue
			}
			elecResult, err := elecValueElement.Text()
			if err != nil {
				continue
			}
			ch <- elecResult
		}
	}()
	select {
	case <-time.After(time.Second * 5):
		return "", fmt.Errorf("获取电费信息超时")
	case elecResult := <-ch:
		return elecResult, nil
	}
}

func Login(wd selenium.WebDriver, ocrResult string) error {
	ch := make(chan selenium.WebElement)
	go func() {
		var usernameElement selenium.WebElement
		for {
			element, err := wd.FindElement(selenium.ByCSSSelector, "#un")
			if err == nil {
				usernameElement = element
				break
			}
			fmt.Println("找不到#un,重试...")
		}
		fmt.Println("找到#un")
		ch <- usernameElement
	}()
	fmt.Println("开启go routine,等待登录页面加载...")
	var usernameElement selenium.WebElement
	select {
	case <-time.After(time.Second * 5):
		return fmt.Errorf("寻找#un超时")
	case usernameElement = <-ch:
	}
	err := usernameElement.SendKeys(config.GlobalConfig.Username)
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
