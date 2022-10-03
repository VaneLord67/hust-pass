package spider

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/tebeka/selenium"
	"hust-pass/config"
	"hust-pass/ocr"
	"log"
	"net/http"
	"os"
	"time"
)

var ocrResult string

const GIFPath = "./captcha.gif"
const JPEGPath = "./captcha.jpeg"

func OCRCallback(response *colly.Response) {
	writeFile, err := os.OpenFile(GIFPath, os.O_SYNC|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer writeFile.Close()
	_, err = writeFile.Write(response.Body)
	if err != nil {
		return
	}
	accessToken, err := ocr.GetAccessToken(config.GlobalConfig.AK, config.GlobalConfig.SK)
	if err != nil {
		return
	}
	err = ocr.WriteJPEG(GIFPath, JPEGPath)
	if err != nil {
		return
	}
	ocrResult, err = ocr.DigitalOCR(accessToken, JPEGPath)
	fmt.Println("识别出验证码:" + ocrResult)
	if err != nil {
		return
	}
}

func LoginGetElec(wd selenium.WebDriver) (string, error) {
	maxRetryCnt := 3
	for ; maxRetryCnt > 0; maxRetryCnt-- {
		for maxRetryCnt := 5; len(ocrResult) != 4 && maxRetryCnt > 0; maxRetryCnt-- {
			jsid, ipPool, err := GetJSIDAndIPPool(wd)
			err = OCR(jsid, ipPool)
			if err != nil {
				return "", err
			}
		}
		if len(ocrResult) != 4 {
			log.Fatal("OCR出错次数过多")
		}
		err := Login(wd)
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

func Login(wd selenium.WebDriver) error {
	elem, err := wd.FindElement(selenium.ByCSSSelector, "#un")
	if err != nil {
		return err
	}
	err = elem.SendKeys(config.GlobalConfig.Username)
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

func OCR(jsid, ipPool string) error {
	c := colly.NewCollector()
	err := c.SetCookies("https://pass.hust.edu.cn", []*http.Cookie{
		{
			Name:   "JSESSIONID",
			Value:  jsid,
			Path:   "/",
			Domain: "pass.hust.edu.cn",
		},
		{
			Name:   "BIGipServerpool-icdc-cas2",
			Value:  ipPool,
			Path:   "/",
			Domain: "pass.hust.edu.cn",
		},
	})
	if err != nil {
		return err
	}
	c.OnResponse(OCRCallback)
	err = c.Visit("https://pass.hust.edu.cn/cas/code")
	if err != nil {
		return err
	}
	return nil
}
