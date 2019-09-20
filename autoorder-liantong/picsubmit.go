package main

import(
	"LianFaPhone/lfp-common"
	. "LianFaPhone/lfp-autoorder-liantong/config"
	"fmt"
	"encoding/json"
	"net/url"
)

type(
	RePicSubmit struct{

	}

	ResPicSubmit struct{
		Total   int   `json:"total"`
		Success  bool   `json:"success"`
		Message string `json:"message"`
	}
)

func (this * RePicSubmit) Send(oId string) error {

	heads := map[string]string{
		"Referer": GPreConfig.NowShopUrl.Url,
		"Content-Length": "6782",
	}
	formBody := make(url.Values)
	formBody.Add("previewImg_front", "data:image/jpeg;base64,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")
	formBody.Add("previewImg_back",  "data:image/jpeg;base64,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")
	formBody.Add("previewImg_person",  "data:image/jpeg;base64,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")


	data,err := common.HttpFormSend(GPreConfig.NowShopUrl.Host+"/viceCard/submitPicNow.do?action=M&orderId="+oId, formBody, "POST", heads)
	if err != nil {
		return err
	}
	res := new(ResPicSubmit)

	err = json.Unmarshal(data, res)
	if err != nil {
		return err
	}
	if !res.Success {
		return fmt.Errorf("%v-%v", res.Success, res.Message)
	}
	return nil
}