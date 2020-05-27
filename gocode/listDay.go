package main

import (
	"fmt"

	"github.com/nongli/ccal"
	"github.com/nongli/zeji"
	"github.com/ying32/govcl/vcl"
)

//农历月历表
func (f *TForm1) OnButton3Click(sensor vcl.IObject) {
	year := f.Edit1.Text()
	month := f.Edit2.Text()
	day := f.Edit3.Text()
	hour := f.Edit4.Text()
	//sx := f.Edit5.Text()
	leapm := f.Edit6.Text()
	sx := "猴"

	if year == "" || month == "" || day == "" || hour == "" || sx == "" || leapm == "" {
		fmt.Printf("输入相应信息\n")
		f.Edit1.SetFocus()
		f.Edit2.SetFocus()
		f.Edit3.SetFocus()
		f.Edit4.SetFocus()
		//f.Edit5.SetFocus()
		f.Edit6.SetFocus()
		return
	}

	if year != "" {

		y, m, d, h, mb := convString2Int(year, month, day, hour, leapm)
		s, l, g, jq := ccal.Input(y, m, d, h, sx, mb)
		iqs := zeji.ZhiSu(s, g)
		x, days, _ := zeji.ListLunarDay(jq, l, iqs)

		if x == 29 {
			vcl.ShowMessageFmt(days[0] + days[1] + days[2] + days[3] + days[4] + days[5] + days[6] +
				days[7] + days[8] + days[9] + days[10] + days[11] + days[12] + days[13] +
				days[14] + days[15] + days[16] + days[17] + days[18] + days[19] + days[20] +
				days[21] + days[22] + days[23] + days[24] + days[2] + days[26] + days[27] +
				days[28])

		} else if x == 30 {
			vcl.ShowMessage(days[0] + days[1] + days[2] + days[3] + days[4] + days[5] + days[6] +
				days[7] + days[8] + days[9] + days[10] + days[11] + days[12] + days[13] +
				days[14] + days[15] + days[16] + days[17] + days[18] + days[19] + days[20] +
				days[21] + days[22] + days[23] + days[24] + days[2] + days[26] + days[27] +
				days[28] + days[29])
		}
	}
}
