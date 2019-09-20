package main


import(
	"LianFaPhone/lfp-common"
	"fmt"
	"encoding/json"
	. "LianFaPhone/lfp-tools/validorder-check/config"
	"bytes"
	"mime/multipart"
)



type (
	ReUnnCheck struct{

	}

	ResUnnCheck struct{
		Code   string   `json:"code"`
		Message  string   `json:"message"`
		Data    []*UnnCheck   `json:"data"`
	}

	UnnCheck struct{
		Mobile  string    `json:"mobile"`
		LastTime string    `json:"lastTime"`
		Area   string     `json:"area"`
		NumberType   string          `json:"numberType"`
		ChargesStatus  string         `json:"chargesStatus"`
		Status  string         `json:"status"`
	}

)

func (this * ReUnnCheck) Send(mobiles string) ([]*UnnCheck, error) {

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	err := w.WriteField("appId", GConfig.ChuangLan.AppId)
	if err != nil {
		return nil, err
	}
	err = w.WriteField("appKey", GConfig.ChuangLan.AppKey)
	if err != nil {
		return nil, err
	}
	err = w.WriteField("mobiles", mobiles)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}

	heads := map[string]string{
		"Content-Type": w.FormDataContentType(),
		"Accept-Charset": "utf-8",
	}


	resData,err := common.HttpSend2(GConfig.ChuangLan.Host+"/open/unn/batch-ucheck", buf, "POST", heads)
	if err != nil {
		return nil,err
	}
	res:=new(ResUnnCheck)
	if err = json.Unmarshal(resData, res); err != nil {
		return nil,err
	}
	if res.Code != "200000" {
		return nil,fmt.Errorf("%v-%v", res.Code, res.Message)
	}
	if res.Data == nil  || len(res.Data) == 0{
		return nil,fmt.Errorf("nil res.Data")
	}
	return res.Data,nil
}
