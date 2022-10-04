package elec

import (
	"fmt"
	"hust-pass/config"
	"log"
	"testing"

	"github.com/imroc/req"
)

func TestElec(t *testing.T) {
	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	resp, err := req.Get(
		fmt.Sprintf("http://sdhq.hust.edu.cn/icbs/PurchaseWebService.asmx/getReserveHKAM?AmMeter_ID=%s",
			config.GlobalConfig.AmMeterID),
		req.Header{
			"Cookie": fmt.Sprintf("Authentication=%s", config.GlobalConfig.Username),
		},
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	type ElecResp struct {
		RemainPower string `xml:"remainPower"`
	}
	elecResp := ElecResp{}
	err = resp.ToXML(&elecResp)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(elecResp)
}
