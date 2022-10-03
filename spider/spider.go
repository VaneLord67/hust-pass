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
		err = wd.Get("http://one.hust.edu.cn")
		if err != nil {
			return "", err
		}
		err = wd.Get("http://sdhq.hust.edu.cn/icbs/hust/html/index.html")
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
		}, time.Second*5)
		if err != nil {
			continue
		}
		return elecResult, nil
	}
	return "", fmt.Errorf("登录失败次数过多")
}

func Login(wd selenium.WebDriver, ocrResult string) error {
	var usernameElement selenium.WebElement
	_ = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		element, err := wd.FindElement(selenium.ByCSSSelector, "#un")
		if err == nil {
			usernameElement = element
			return true, nil
		}
		return false, nil
	}, time.Second*5)
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
