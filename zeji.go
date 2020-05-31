package main

import (
	"fmt"
	"os"

	"github.com/nongli/ccal"
	"github.com/nongli/ganzhi"
	"github.com/nongli/solar"
	"github.com/nongli/zeji"
	"github.com/ying32/govcl/vcl"
)

func (f *TForm1) OnButton4Click(sender vcl.IObject) {
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
	}

	switch inputb {
	case true:
		if sxt := shengxiao(sx); sxt == false {
			vcl.ShowMessageFmt("\"生肖\"输入错误，系統退出\n")
			os.Exit(0)
		}
		s, l, g, jq := ccal.Input(y, m, d, h, sx, mb)

		//纪年信息
		solarinfo := fmt.Sprintf("阳历纪年: %d年-%d月-%d日-周%s-阳历时间范围:%s\n", s.SYear, s.SMonth, s.SDay, s.SWeek, s.SHour)
		lunarinfo := fmt.Sprintf("农历纪年: %d年-%d月(%s)-%d日-%d时(%s时) 本年是否有闰月:%t 闰%d月\n",
			l.LYear, l.LMonth, l.LYdxs, l.LDay, l.LHour, l.LaliasHour, l.Leapmb, l.LeapMonth)
		gzinfo := fmt.Sprintf("干支纪年: %s%s年-%s月-%s%s日-%s时\n\n",
			g.YearGanM, g.YearZhiM, g.MonthGanZhiM, g.DayGanM, g.DayZhiM, g.HourGanZhiM)

		//值宿信息
		iqs := zeji.ZhiSu(s, g)
		ws := iqs.Ws //值宿名称
		wn := iqs.Wn //值宿当日周几　0为周日
		winfo := fmt.Sprintf("周%s 值宿:\"%s\"\n", solar.Zhou[wn], ws)

		//择吉数字
		n1, n2, n3 := zeji.JiGan(l.LMonth)
		number := fmt.Sprintf("农历本月吉干数字:%d %d %d-->", n1, n2, n3)
		name := fmt.Sprintf("吉干:%s %s %s\n", ganzhi.Gan[n1], ganzhi.Gan[n2], ganzhi.Gan[n3])

		//农历月份吉干
		jg, _, aliasZhi := zeji.ListDay(jq, l, iqs)
		jgs := zeji.ShowJiGan(jq.Sx, jg, aliasZhi)

		//择吉结果
		yearZhi := g.YearZhi
		nx := zeji.AllNumber(yearZhi, m, d, h)
		n1b := nx.YiPan()
		n2b := nx.ErPan()
		n3b := nx.SanPan()
		result := zeji.ShowResult(n1b, n2b, n3b)

		//七煞
		_, qsday := zeji.QiShaDay(jq, l, iqs)

		//月份吉干列表
		listJg := fmt.Sprintf("本月吉干列表:\n%s\n", jgs)

		//信息显示到UI界面
		vcl.ShowMessage(solarinfo + lunarinfo + gzinfo + winfo + number + name + result + qsday + listJg)
	case false:
		vcl.ShowMessage("数字输入错误，系统退出\n")
		os.Exit(0)
	}

}
