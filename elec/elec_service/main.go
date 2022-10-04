package main

import (
	"encoding/json"
	"fmt"
	"hust-pass/config"
	"hust-pass/elec"
	"log"
	"net/http"
)

func elecQueryHandler(w http.ResponseWriter, r *http.Request) {
	elecResult, err := elec.GetElecResult()
	if err != nil {
		_, _ = w.Write([]byte("get elec result error"))
		log.Println(err)
		return
	}
	type ElecQueryResp struct {
		Elec string `json:"elec"`
	}
	elecQueryResp := ElecQueryResp{Elec: elecResult}
	respBytes, err := json.Marshal(&elecQueryResp)
	if err != nil {
		_, _ = w.Write([]byte("json Marshal error"))
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
	log.Println("HTTP服务启动...")
	http.HandleFunc("/elec_query", elecQueryHandler)
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.GlobalConfig.ElecQueryServicePort), nil)
	if err != nil {
		log.Fatal(err)
	}
}
