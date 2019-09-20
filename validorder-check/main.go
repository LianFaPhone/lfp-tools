package main


import (
	"go.uber.org/zap"
	"runtime/debug"
	"fmt"
	"LianFaPhone/lfp-base/config"
	. "LianFaPhone/lfp-base/log/zap"
	. "LianFaPhone/lfp-tools/validorder-check/config"

	"time"
	"github.com/tealeg/xlsx"
	"strings"
	"os"
)
//
//func main(){
//	//读取命令行
//	//读取配置
//
//	//读取excel
//	//调用接口
//	//输出excel
//
//
//}

var AllExcel    *Excel
var SuccExcel   *Excel
var FailExcel   *Excel
var timeNowStr  string

func PanicPrint() {
	if err := recover(); err != nil {
		ZapLog().With(zap.Any("error", err)).Error(string(debug.Stack()))
	}
}

func main() {
	laxFlag := config.NewLaxFlagDefault()
	cfgPath := laxFlag.String("conf_path", "config.yaml", "config path")
	logPath := laxFlag.String("log_path", "zap.conf", "log conf path")
	laxFlag.LaxParseDefault()
	fmt.Printf("command param: conf_path=%s, log_path=%s\n", *cfgPath, *logPath)
	LoadConfig(*cfgPath)
	LoadZapConfig(*logPath)
	ZapLog().Sugar().Infof("Config Content[%v]", GConfig)
	defer func(){
		time.Sleep(time.Second*10)
	}()
	defer ZapClose()
	defer PanicPrint()

	timeNowStr = time.Now().Format("2006-01-02-15-04-05")
	if err := Init(); err != nil {
		//老文件文件名改成日期
		ZapLog().Sugar().Errorf("Init err %v,  so exist",err)
		return
	}

	ZapLog().Sugar().Infof("*****************Start Read input excel and time[%v]........", timeNowStr)
	colnames, sheet, err := ReadInput()
	if err != nil {
		ZapLog().Sugar().Errorf("ReadInput err %v,  so exist",err)
		return
	}
	outcolnames := make([]string, 0)
	outcolnames = append(outcolnames, colnames...)
	outcolnames = append(outcolnames,  "状态", "错误信息")

	if err = AllExcel.AddHeader(outcolnames); err != nil {
		ZapLog().Sugar().Errorf("AllExcel.AddHeader err %v,  so exist",err)
		return
	}
	if err = SuccExcel.AddHeader(outcolnames); err != nil {
		ZapLog().Sugar().Errorf("SuccExcel.AddHeader err %v,  so exist",err)
		return
	}
	if err = FailExcel.AddHeader(outcolnames); err != nil {
		ZapLog().Sugar().Errorf("FailExcel.AddHeader err %v,  so exist",err)
		return
	}
	//if err = new(InputIndex).Parse(outcolnames); err != nil {
	//	ZapLog().Sugar().Errorf("InputIndex).Parse err %v,  so exist",err)
	//	return
	//}
	ZapLog().Sugar().Infof("**Start Read excel record [[%d]], and sleep 5 second to continue", len(sheet.Rows)-1)
	time.Sleep(time.Second * 5)

	//执行操作

	if err := Work(outcolnames, sheet); err != nil {
		ZapLog().Sugar().Errorf("work err %v so exist", err)
		return
	}
	//输出文件

	ZapLog().Sugar().Infof("\n\n**************Success done*****************")
}

func Init() error {
	err := os.MkdirAll(GConfig.Server.OutputPath, os.ModePerm)
	//cmd := exec.Command("mkdir", "-p", GConfig.Server.OutputPath)
	//err := cmd.Run()
	if err != nil {
		ZapLog().Sugar().Errorf("mkdir %v err %v", GConfig.Server.OutputPath, err)
		return err
	}

	AllExcel, err = NewExcel("", GConfig.Server.OutputPath + "/all-"+timeNowStr+".xlsx")
	if err != nil {
		ZapLog().Sugar().Errorf("NewExcel %v err %v", GConfig.Server.OutputPath + "/all-"+timeNowStr+".xlsx", err)
		return err
	}

	SuccExcel, err = NewExcel("", GConfig.Server.OutputPath + "/success-"+timeNowStr+".xlsx")
	if err != nil {
		ZapLog().Sugar().Errorf("NewExcel %v err %v", GConfig.Server.OutputPath + "\\success-"+timeNowStr+".xlsx", err)
		return err
	}
	FailExcel, err = NewExcel("", GConfig.Server.OutputPath + "/fail-"+timeNowStr+".xlsx")
	if err != nil {
		ZapLog().Sugar().Errorf("NewExcel %v err %v", GConfig.Server.OutputPath + "/fail-"+timeNowStr+".xlsx", err)
		return err
	}

	return nil
}

func ReadInput() ([]string, *xlsx.Sheet, error) {
	excelFileName := GConfig.Server.InputPath
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		return nil,nil,err
	}
	sheet, ok := xlFile.Sheet["Sheet1"]
	if !ok {
		return nil,nil, fmt.Errorf("nofind Sheet1")
	}

	rowNames := make([]string, 0)
	for _, cell := range sheet.Rows[0].Cells {
		text := cell.String()
		text = strings.Replace(text, " ", "", -1)
		rowNames = append(rowNames, text)
	}
	if len(rowNames) <= 1 {
		return nil,nil, fmt.Errorf("in.xlsx format err")
	}
	return rowNames, sheet, nil
}


func Work(colnames []string, sheet *xlsx.Sheet) error {
	//读取
	//创建输出文件 总表，成功，失败， 二次


	AllRecords := len(sheet.Rows) - 1
	SuccessRecords := 0
	FailRecords := 0

	for i:=0; i< len(sheet.Rows); i++ {
		if i == 0 { //第一行标题栏
			continue
		}
		row := sheet.Rows[i]

		inIdCard := strings.Replace(row.Cells[2].String(), " ", "", -1)
		inXinMing := strings.Replace(row.Cells[1].String(), " ", "", -1)
		inPhone := strings.Replace(row.Cells[3].String(), " ", "", -1)
		if len(inIdCard) <= 0 && len(inXinMing) <= 0 && len(inPhone) <= 0{
			continue
		}

		excelArr := CellsToArr(row.Cells)
		fmt.Println("\n============================================")
		ZapLog().Sugar().Infof("************record= [%v]", excelArr)

		if GConfig.Server.PhoneCheckFlag == 1 {
			if len(inPhone) <= 0 {
				excelArr = append(excelArr,  "数据缺失", "手机数据缺失")
				if err2 := AllExcel.Append(excelArr); err2 != nil {
					ZapLog().Sugar().Errorf("AllExcel.Append error %v", err2)
				}
				if err2 := FailExcel.Append(excelArr); err2 != nil {
					ZapLog().Sugar().Errorf("SuccExcel.Append error %v", err2)
				}
				FailRecords++
				continue
			}
			if ok := VerifyMobile(inPhone, excelArr); !ok {
				FailRecords++
				continue
			}
			ZapLog().Sugar().Infof("VerifyMobile Ok")
		}
		if GConfig.Server.IdcardCheckFlag == 1 {
			if len(inIdCard) <= 0 || len(inXinMing) <= 0 {
				excelArr = append(excelArr,  "数据缺失", "姓名或身份证数据缺失")
				if err2 := AllExcel.Append(excelArr); err2 != nil {
					ZapLog().Sugar().Errorf("AllExcel.Append error %v", err2)
				}
				if err2 := FailExcel.Append(excelArr); err2 != nil {
					ZapLog().Sugar().Errorf("SuccExcel.Append error %v", err2)
				}
				FailRecords++
				continue
			}
			if ok := VerifyIdcard(inXinMing, inIdCard, excelArr); !ok {
				FailRecords++
				continue
			}
			ZapLog().Sugar().Infof("VerifyIdcard Ok")
		}

		excelArr = append(excelArr,  "成功", "成功")
		if err2 := AllExcel.Append(excelArr); err2 != nil {
				ZapLog().Sugar().Errorf("AllExcel.Append error %v", err2)
		}
		if err2 := SuccExcel.Append(excelArr); err2 != nil {
				ZapLog().Sugar().Errorf("SuccExcel.Append error %v", err2)
		}
		SuccessRecords++

		if GConfig.Server.Intvl > 0 {
			time.Sleep(time.Millisecond * time.Duration(GConfig.Server.Intvl))
		}
	}
	ZapLog().Sugar().Infof("****The End**********AllRecords[%v] SuccessRecords[%v] FailRecords[%v] ", AllRecords, SuccessRecords, FailRecords)
	return nil
}



//检测姓名身份证，返回 是否成功
func VerifyIdcard(inXinMing, inIdCard string, excelArr []string) bool {
	res,err :=  new(ReIdCardCheck).Send(inXinMing, inIdCard)
	if err != nil {
		ZapLog().Sugar().Errorf("ReIdCardCheck.Send error %v", err)
		excelArr = append(excelArr, "系统问题", err.Error())
		if err2 := AllExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("AllExcel.Append error %v", err2)
		}
		if err2 := FailExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("FailExcel.Append error %v", err2)
		}
		return false
	}
	switch res.Result {
	case "01":
		return true
	case "02":
		excelArr = append(excelArr, "身份校验失败", res.Result +"-"+res.Remark)
		if err2 := AllExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("AllExcel.Append error %v", err2)
		}
		if err2 := FailExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("FailExcel.Append error %v", err2)
		}
		return false
	case "03":
		excelArr = append(excelArr, "身份校验不确定", res.Result +"-"+res.Remark)
		if err2 := AllExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("AllExcel.Append error %v", err2)
		}
		if err2 := FailExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("FailExcel.Append error %v", err2)
		}
		return false
	case "04":
		excelArr = append(excelArr, "身份校验不确定", res.Result +"-"+res.Remark)
		if err2 := AllExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("AllExcel.Append error %v", err2)
		}
		if err2 := FailExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("FailExcel.Append error %v", err2)
		}
		return false
	default:
		excelArr = append(excelArr, "身份校验不确定", res.Result +"-"+res.Remark)
		if err2 := AllExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("AllExcel.Append error %v", err2)
		}
		if err2 := FailExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("FailExcel.Append error %v", err2)
		}
		return false
	}
	return false
}

func VerifyMobile(mobile string, excelArr []string) bool {

	resArr,err := new(ReUnnCheck).Send(mobile)
	if err != nil {
		ZapLog().Sugar().Errorf("ReUnnCheck).Send err=%v", err)
		excelArr = append(excelArr, "系统问题", err.Error())
		if err2 := AllExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("AllExcel.Append error %v", err2)
		}
		if err2 := FailExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("FailExcel.Append error %v", err2)
		}
		return false
	}

	// 0：空号 1：实号 2：停机 3：库无 4：沉默号 5：风险号
	switch resArr[0].Status {
	case "0":
		ZapLog().Sugar().Errorf("mobile %v is 空号 %v ",mobile,resArr[0].Status)
		excelArr = append(excelArr, "手机空号", resArr[0].Status)
		if err2 := AllExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("AllExcel.Append error %v", err2)
		}
		if err2 := FailExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("FailExcel.Append error %v", err2)
		}
		return false
	case "1":
		return true
	case "2":
		ZapLog().Sugar().Errorf("mobile %v is 停机 %v ",mobile,resArr[0].Status)
		excelArr = append(excelArr, "手机停机", resArr[0].Status)
		if err2 := AllExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("AllExcel.Append error %v", err2)
		}
		if err2 := FailExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("FailExcel.Append error %v", err2)
		}
		return false
	case "3":
		ZapLog().Sugar().Errorf("mobile %v is 库无 %v ",mobile,resArr[0].Status)
		excelArr = append(excelArr, "手机无法确定", resArr[0].Status)
		if err2 := AllExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("AllExcel.Append error %v", err2)
		}
		if err2 := FailExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("FailExcel.Append error %v", err2)
		}
		return false
	case "4":
		ZapLog().Sugar().Errorf("mobile %v is 沉默号 %v",mobile,resArr[0].Status)
		excelArr = append(excelArr, "手机沉默号", resArr[0].Status)
		if err2 := AllExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("AllExcel.Append error %v", err2)
		}
		if err2 := FailExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("FailExcel.Append error %v", err2)
		}
		return false
	case "5":
		ZapLog().Sugar().Errorf("mobile %v is 风险号 %v",mobile,resArr[0].Status)
		excelArr = append(excelArr, "手机风险号", resArr[0].Status)
		if err2 := AllExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("AllExcel.Append error %v", err2)
		}
		if err2 := FailExcel.Append(excelArr); err2 != nil {
			ZapLog().Sugar().Errorf("FailExcel.Append error %v", err2)
		}
		return false
	}

	ZapLog().Sugar().Errorf("mobile %v is 不确定 %v",mobile,resArr[0].Status)
	excelArr = append(excelArr, "手机无法确定", resArr[0].Status)
	if err2 := AllExcel.Append(excelArr); err2 != nil {
		ZapLog().Sugar().Errorf("AllExcel.Append error %v", err2)
	}
	if err2 := FailExcel.Append(excelArr); err2 != nil {
		ZapLog().Sugar().Errorf("FailExcel.Append error %v", err2)
	}
	return false
}