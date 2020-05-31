// 由res2go自动生成。
// 在这里写你的事件。

package main

import (
	"fmt"
	"os"

	"github.com/nongli/ccal"
	_ "github.com/ying32/govcl/pkgs/winappres"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"github.com/ying32/govcl/vcl/types/keys"
)

//::private::
type TForm1Fields struct {
}

//基础纪年信息
func (f *TForm1) OnButton1Click(object vcl.IObject) {

	year := f.Edit1.Text()
	month := f.Edit2.Text()
	day := f.Edit3.Text()
	hour := f.Edit4.Text()
	sx := f.Edit5.Text()
	leapm := f.Edit6.Text()

	y, m, d, h, inputb := String2Int(year, month, day, hour)
	mb, err := leapBool(leapm)
	if err != nil {
		s := fmt.Sprintf(err.Error())
		vcl.ShowMessage(s)
		os.Exit(0)
	}

	switch inputb {
	case true: //数字输入正确
		//生肖判断
		if sxt := shengxiao(sx); sxt == false {
			vcl.ShowMessageFmt("生肖输入错误\n")
			//os.Exit(0)
		}

		s, l, g, _ := ccal.Input(y, m, d, h, sx, mb)
		solarinfo := fmt.Sprintf("阳历纪年: %d年-%d月-%d日-周%s-阳历时间范围:%s\n", s.SYear, s.SMonth, s.SDay, s.SWeek, s.SHour)
		lunarinfo := fmt.Sprintf("农历纪年: %d年-%d月(%s)-%d日-%d时(%s时) 本年是否有闰月:%t 闰%d月\n",
			l.LYear, l.LMonth, l.LYdxs, l.LDay, l.LHour, l.LaliasHour, l.Leapmb, l.LeapMonth)
		gzinfo := fmt.Sprintf("干支纪年: %s%s年-%s月-%s%s日-%s时\n\n",
			g.YearGanM, g.YearZhiM, g.MonthGanZhiM, g.DayGanM, g.DayZhiM, g.HourGanZhiM)

		//信息显示到UI界面
		vcl.ShowMessage(solarinfo + lunarinfo + gzinfo)
	case false:
		vcl.ShowMessage("数字输入错误\n")
		//os.Exit(0)
	}
}

//生成可重复的点击事件(获取输入数据)
func (f *TForm1) OnFormKeyPress(sender vcl.IObject, key *types.Char) {
	fmt.Println("key:", *key)
	for {
		go func() {
			if *key == keys.VkReturn {
				f.Button1.Click() //基础信息
				f.Button2.Click() //24节气
				f.Button3.Click() //月历表
				f.Button4.Click() //择吉
				f.Button6.Click() //地母经
				f.Button5.Click() //关于
				f.Button7.Click() //禁忌
				f.Button8.Click() //其他
			}
		}()
	}
}
