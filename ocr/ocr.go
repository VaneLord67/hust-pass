package ocr

import (
	"encoding/base64"
	"fmt"
	"github.com/gocolly/colly/v2"
	"hust-pass/config"
	"image/gif"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/imroc/req"
)

const GIFPath = "./captcha.gif"
const JPEGPath = "./captcha.jpeg"

var ocrResult string

func WriteJPEG(gifPath string, JPEGPath string) error {
	gifFile, err := os.Open(gifPath)
	if err != nil {
		return err
	}
	defer gifFile.Close()
	images, err := gif.DecodeAll(gifFile)
	if err != nil {
		return err
	}
	if len(images.Image) < 2 {
		return fmt.Errorf("gif length error")
	}
	writeFile, err := os.OpenFile(JPEGPath, os.O_SYNC|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer writeFile.Close()
	err = jpeg.Encode(writeFile, images.Image[1], nil)
	if err != nil {
		return err
	}
	return nil
}

func GetAccessToken(ak, sk string) (string, error) {
	resp, err := req.Post("https://aip.baidubce.com/oauth/2.0/token", req.QueryParam{
		"grant_type":    "client_credentials",
		"client_id":     ak,
		"client_secret": sk,
	})
	if err != nil {
		return "", err
	}
	type AccessTokenResp struct {
		AccessToken string `json:"access_token"`
	}
	var accessTokenResp AccessTokenResp
	err = resp.ToJSON(&accessTokenResp)
	if err != nil {
		return "", err
	}
	return accessTokenResp.AccessToken, nil
}

func DigitalOCR(accessToken string, imagePath string) (string, error) {
	imageString, err := GetImageString(imagePath)
	if err != nil {
		return "", err
	}
	resp, err := req.Post(
		"https://aip.baidubce.com/rest/2.0/ocr/v1/numbers",
		req.Param{"access_token": accessToken, "image": imageString},
		req.Header{"Content-Type": "application/x-www-form-urlencoded"},
	)
	if err != nil {
		return "", err
	}
	type WordsResult struct {
		Words string `json:"words"`
	}
	type OCRResp struct {
		WordsResults []WordsResult `json:"words_result"`
	}
	ocrResp := OCRResp{}
	err = resp.ToJSON(&ocrResp)
	if err != nil {
		return "", err
	}
	if len(ocrResp.WordsResults) == 0 {
		return "", fmt.Errorf("无法识别")
	}
	return ocrResp.WordsResults[0].Words, nil
}

func GetImageString(imagePath string) (string, error) {
	imageFile, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return "", err
	}
	const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var coder = base64.NewEncoding(base64Table)
	imgString := coder.EncodeToString(imageFile)
	return imgString, nil
}

func OCR(jsid, ipPool string) (string, error) {
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
		return "", err
	}
	c.OnResponse(OCRCallback)
	err = c.Visit("https://pass.hust.edu.cn/cas/code")
	if err != nil {
		return "", err
	}
	return ocrResult, nil
}

func OCRCallback(response *colly.Response) {
	writeFile, err := os.OpenFile(GIFPath, os.O_SYNC|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(writeFile *os.File) {
		err := writeFile.Close()
		if err != nil {
			log.Println(err)
		}
	}(writeFile)
	_, err = writeFile.Write(response.Body)
	if err != nil {
		log.Println(err)
		return
	}
	accessToken, err := GetAccessToken(config.GlobalConfig.AK, config.GlobalConfig.SK)
	if err != nil {
		log.Println(err)
		return
	}
	err = WriteJPEG(GIFPath, JPEGPath)
	if err != nil {
		log.Println(err)
		return
	}
	ocrResult, err = DigitalOCR(accessToken, JPEGPath)
	log.Println("识别出验证码:" + ocrResult)
	if err != nil {
		log.Println(err)
		return
	}
}
