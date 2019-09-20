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
	ReIdCardCheck struct{

	}

	ResIdCardCheck struct{
		Code   string   `json:"code"`
		Message  string   `json:"message"`
		Data    *IdCardCheck   `json:"data"`
	}

	IdCardCheck struct{
		OrderNo  string    `json:"orderNo"`
		HandleTime string    `json:"handleTime"`
		Province   string     `json:"province"`
		City   string          `json:"city"`
		Country  string         `json:"country"`
		Birthday  string         `json:"birthday"`
		Age    string      `json:"age"`
		Gender   string    `json:"gender"`
		Remark   string     `json:"remark"`
		Result    string    `json:"result"`

	}

)

func (this * ReIdCardCheck) Send(name, idcard string) (*IdCardCheck, error) {

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
	err = w.WriteField("name", name)
	if err != nil {
		return nil, err
	}
	err = w.WriteField("idNum", idcard)
	if err != nil {
		return nil,err
	}
	err = w.Close()
	if err != nil {
		return nil,err
	}

	heads := map[string]string{
		"Content-Type": w.FormDataContentType(),
		"Accept-Charset": "utf-8",
	}

	resData,err := common.HttpSend2(GConfig.ChuangLan.Host+"/open/idcard/id-card-auth", buf, "POST", heads)
	if err != nil {
		return nil,err
	}
	res:=new(ResIdCardCheck)
	if err = json.Unmarshal(resData, res); err != nil {
		return nil,err
	}
	if res.Code != "200000" {
		return nil,fmt.Errorf("%v-%v", res.Code, res.Message)
	}
	if res.Data == nil  {
		return nil,fmt.Errorf("nil res.Data")
	}
	//fmt.Println("res= ",res)
	//fmt.Println("redata= ",res.Data)
	return res.Data,nil
}
