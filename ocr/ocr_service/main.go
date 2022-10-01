package main

import (
	"encoding/json"
	"fmt"
	"hust-pass/config"
	"hust-pass/ocr"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func ocrHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("callerToken") != config.GlobalConfig.CallerToken {
		w.Write([]byte("callerToken错误"))
		return
	}
	file, fh, err := r.FormFile("file")
	if err != nil {
		w.Write([]byte("未读取到Form中的file"))
		return
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		w.Write([]byte("读取file错误"))
		return
	}
	gifPath := fmt.Sprintf("./%d-%s", time.Now().Unix(), fh.Filename)
	writeFile, err := os.OpenFile(gifPath, os.O_SYNC|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		w.Write([]byte("打开文件错误"))
		return
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Println(err)
		}
	}(gifPath)
	defer writeFile.Close()
	_, err = writeFile.Write(bytes)
	if err != nil {
		log.Println(err)
		return
	}
	jpegPath := fmt.Sprintf("./%d.jpeg", time.Now().Unix())
	err = ocr.WriteJPEG(gifPath, jpegPath)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Println(err)
		}
	}(jpegPath)
	accessToken, err := ocr.GetAccessToken(config.GlobalConfig.AK, config.GlobalConfig.SK)
	if err != nil {
		w.Write([]byte("accessToken获取错误"))
		return
	}
	ocrResult, err := ocr.DigitalOCR(accessToken, jpegPath)
	if err != nil {
		w.Write([]byte("ocr出错"))
		return
	}
	type OCRServiceResp struct {
		Words string `json:"words"`
	}
	ocrServiceResp := OCRServiceResp{Words: ocrResult}
	respBytes, err := json.Marshal(&ocrServiceResp)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(respBytes)
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}
	http.HandleFunc("/ocr", ocrHandler)
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.GlobalConfig.OCRServicePort), nil)
	if err != nil {
		log.Fatal(err)
	}
}
