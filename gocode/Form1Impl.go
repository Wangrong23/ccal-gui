// 由res2go自动生成。
// 在这里写你的事件。

package main

import (
	"fmt"
	"strconv"

	"github.com/nongli/ccal"
	_ "github.com/ying32/govcl/pkgs/winappres"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"github.com/ying32/govcl/vcl/types/keys"
)

//::private::
type TForm1Fields struct {
}

//获取输入信息
func (f *TForm1) OnButton1Click(object vcl.IObject) {

	year := f.Edit1.Text()
	month := f.Edit2.Text()
	day := f.Edit3.Text()
	hour := f.Edit4.Text()
	sx := f.Edit5.Text()
	leapm := f.Edit6.Text()
	if year == "" || month == "" || day == "" || hour == "" || sx == "" || leapm == "" {
		fmt.Printf("输入相应信息\n")
		f.Edit1.SetFocus()
		f.Edit2.SetFocus()
		f.Edit3.SetFocus()
		f.Edit4.SetFocus()
		f.Edit5.SetFocus()
		f.Edit6.SetFocus()
		return
	}
	if year != "" {
		y, m, d, h, mb := convString2Int(year, month, day, hour, leapm)
		s, l, g, _ := ccal.Input(y, m, d, h, sx, mb)
		solarinfo := fmt.Sprintf("阳历纪年: %d年-%d月-%d日-周%s-阳历时间范围:%s\n", s.SYear, s.SMonth, s.SDay, s.SWeek, s.SHour)
		lunarinfo := fmt.Sprintf("农历纪年: %d年-%d月(%s)-%d日-%d时(%s时) 本年是否有闰月:%t 闰%d月\n",
			l.LYear, l.LMonth, l.LYdxs, l.LDay, l.LHour, l.LaliasHour, l.Leapmb, l.LeapMonth)
		gzinfo := fmt.Sprintf("干支纪年: %s%s年-%s月-%s%s日-%s时\n\n",
			g.YearGanM, g.YearZhiM, g.MonthGanZhiM, g.DayGanM, g.DayZhiM, g.HourGanZhiM)
		//信息显示到UI界面
		vcl.ShowMessage(solarinfo + lunarinfo + gzinfo)
	}
}

//生成点击事件(获取输入数据)
func (f *TForm1) OnFormKeyPress(sender vcl.IObject, key *types.Char) {
	fmt.Println("key:", *key)
	if *key == keys.VkReturn {
		f.Button1.Click() //基础信息
		f.Button2.Click() //24节气
		f.Button3.Click() //月历表
		f.Button4.Click() //择吉
		f.Button6.Click() //地母经
		f.Button5.Click() //关于
	}
}

//类型转换
func convString2Int(year, month, day, hour, leapm string) (y, m, d, h int, mb bool) {

	_y, _ := strconv.ParseInt(year, 10, 32)
	_m, _ := strconv.ParseInt(month, 10, 32)
	_d, _ := strconv.ParseInt(day, 10, 32)
	_h, _ := strconv.ParseInt(hour, 10, 32)
	y, m, d, h = int(_y), int(_m), int(_d), int(_h)

	if leapm == "y" || leapm == "Y" {
		mb = true
	} else {
		mb = false
	}

	return
}

func (f *TForm1) OnEdit2Change(sender vcl.IObject) {

}
