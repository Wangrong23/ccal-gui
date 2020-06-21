package main

import (
	"fmt"
	"log"

	"github.com/nongli/ccal"
	"github.com/nongli/solar"
	"github.com/ying32/govcl/vcl"
)

//24节气按钮　只需要年份数据即可
func (f *TForm1) OnButton2Click(sender vcl.IObject) {
	year := f.Edit1.Text()
	y, m, d, h, inputb, err := String2Int(year, "1", "1", "1")
	if err != nil {
		s := fmt.Sprint(err)
		vcl.ShowMessage(s)
	}

	switch inputb {
	case true:
		err, _, _, _, jq := ccal.Input(y, m, d, h, "猴", false)
		if err != nil {
			log.Fatal(err)
		}
		jq24, _ := solar.ShowJieqi24(jq.Jqt, jq.Jq11t)

		//信息显示到UI界面
		vcl.ShowMessage(jq24[0] + jq24[1] + jq24[2] + jq24[3] + jq24[4] + jq24[5] + jq24[6] + jq24[7] + jq24[8] +
			jq24[9] + jq24[10] + jq24[11] + jq24[12] + jq24[13] + jq24[14] + jq24[15] + jq24[16] + jq24[17] +
			jq24[18] + jq24[19] + jq24[20] + jq24[21] + jq24[22] + jq24[23] + jq24[24] + jq24[25] + jq24[26])
	case false:
		vcl.ShowMessage("年份数字输入错误\n")
	}
}
