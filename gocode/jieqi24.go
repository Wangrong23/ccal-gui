package main

import (
	"fmt"
	"strconv"

	"github.com/nongli/ccal"
	"github.com/nongli/solar"
	"github.com/ying32/govcl/vcl"
)

//24节气按钮　只需要年份数据即可
func (f *TForm1) OnButton2Click(sender vcl.IObject) {
	year := f.Edit1.Text()
	if year == "" {
		fmt.Printf("输入年份数字\n")
		f.Edit1.SetFocus()
		return
	}

	if year != "" {
		_y, _ := strconv.ParseInt(year, 10, 32)
		y := int(_y)
		_, _, _, jq := ccal.Input(y, 1, 1, 1, "猴", false)
		jq24 := solar.ShowJieqi24(jq.Jqt, jq.Jq11t)

		//信息显示到UI界面
		vcl.ShowMessage(jq24[0] + jq24[1] + jq24[2] + jq24[3] + jq24[4] + jq24[5] + jq24[6] + jq24[7] + jq24[8] +
			jq24[9] + jq24[10] + jq24[11] + jq24[12] + jq24[13] + jq24[14] + jq24[15] + jq24[15] + jq24[16] +
			jq24[17] + jq24[18] + jq24[19] + jq24[20] + jq24[21] + jq24[22] + jq24[23] + jq24[24] + jq24[25] + jq24[26])
	}
}

/* //生成点击事件(获取输入数据)
func (f *TForm1) OnFormKeyPressBt2(sender vcl.IObject, key *types.Char) {
	fmt.Println("key:", *key)
	if *key == keys.VkReturn {
		f.Button2.Click()
	}
} */
