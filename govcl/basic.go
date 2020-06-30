// 由res2go自动生成。
// 在这里写你的事件。

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nongli/ccal"
	"github.com/nongli/lunar"

	//	_ "github.com/ying32/govcl/pkgs/winappres"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"github.com/ying32/govcl/vcl/types/keys"
)

//::private::
type TForm1Fields struct {
}

var lm string

//基础纪年信息
func (f *TForm1) OnButton1Click(object vcl.IObject) {

	year := f.Edit1.Text()
	month := f.Edit2.Text()
	day := f.Edit3.Text()
	hour := f.Edit4.Text()
	sx := f.Edit5.Text()
	leapm := f.Edit6.Text()

	y, m, d, h, inputb, err := String2Int(year, month, day, hour)
	if err != nil {
		s := fmt.Sprint(err)
		vcl.ShowMessage(s)
	}

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
			vcl.ShowMessageFmt("生肖输入错误")
			//os.Exit(0)
		}

		err, s, l, g, _ := ccal.Input(y, m, d, h, sx, mb)
		if err != nil {
			log.Fatal(err)
		}
		if l.Leapmb == true {
			lm = "是"
		} else {
			lm = "否"
		}

		solarinfo := fmt.Sprintf("阳历纪年: %d年-%d月-%d日-周%s-阳历时间:%d:%d\n", s.SYear, s.SMonth, s.SDay, s.SWeek, T.Hour(), T.Minute())
		lunarinfo := fmt.Sprintf("农历纪年: %d年%s月(%s)%s %s时(%d时)\n本年是否有闰月:%s 闰%d月\n",
			l.LYear, lunar.Ymc[l.LMonth-1], l.LYdxs, lunar.Rmc[l.LDay-1], l.LaliasHour, l.LHour, lm, l.LeapMonth)
		gzinfo := fmt.Sprintf("干支纪年: %s%s年-%s月-%s%s日-%s时\n\n",
			g.YearGanM, g.YearZhiM, g.MonthGanZhiM, g.DayGanM, g.DayZhiM, g.HourGanZhiM)

		//杨公祭日
		yginfo := yg13(l.LMonth, l.LDay)

		//信息显示到UI界面
		vcl.ShowMessage(solarinfo + lunarinfo + gzinfo + yginfo)
	case false:
		vcl.ShowMessage("数字输入错误")
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
				f.Button7.Click() //其他内容
				f.Button8.Click() //显示当日
			}
		}()
	}
}
