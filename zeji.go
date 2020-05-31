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
		zhisuInfo := fmt.Sprintf(iqs.ZhiSuInfo) //当日值宿信息

		//判断当日是否为七煞日
		qsB := iqs.IsQiSha(s.SolarDayT, g.DayZhiM)
		_, qsNumber, _ := zeji.QiShaInfo(int(s.SolarDayT.Weekday()), g.DayZhiM)
		isQiSha := PrintQS(qsB, qsNumber)

		//择吉数字
		n1, n2, n3 := zeji.JiGan(l.LMonth)
		number := fmt.Sprintf("\n农历本月吉干数字:%d %d %d-->", n1, n2, n3)
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
		result := zeji.ShowResult(n1b, n2b, n3b, qsB)

		//四绝日 3:立春...9:立夏...15:立秋...22:立冬
		//四大绝日 6:春分...12:夏至...19:秋分...25:冬至
		var jlr string
		_, jueRit := solar.ShowJieqi24(jq.Jqt, jq.Jq11t)

		lichunT := s.SolarDayT.Equal(jueRit[0]) //立春
		lixiaT := s.SolarDayT.Equal(jueRit[1])
		liqiuT := s.SolarDayT.Equal(jueRit[2])
		lidongT := s.SolarDayT.Equal(jueRit[3])

		chunfenT := s.SolarDayT.Equal(jueRit[4])
		xiazhiT := s.SolarDayT.Equal(jueRit[5])
		qiufenT := s.SolarDayT.Equal(jueRit[6])
		dongzhiT := s.SolarDayT.Equal(jueRit[7])

		if lichunT == true || lixiaT == true || liqiuT == true || lidongT == true ||
			chunfenT == true || xiazhiT == true || qiufenT == true || dongzhiT == true {
			jlr = fmt.Sprintf("[此日为四绝(离)日]\n绝(离)日: 立春 立夏 立秋 立冬 春分 夏至 秋分 冬至\n\n")
		}

		//月份吉干列表
		listJg := fmt.Sprintf("本月吉干列表:\n%s\n", jgs)

		//信息显示到UI界面
		vcl.ShowMessage(solarinfo + lunarinfo + gzinfo + winfo + zhisuInfo + isQiSha + number + name + result + jlr + listJg)
	case false:
		vcl.ShowMessage("数字输入错误，系统退出\n")
		os.Exit(0)
	}

}

//打印当时是否为七煞日
func PrintQS(qsb bool, qsn int) (info string) {

	if qsb == true {
		info = QiSha(qsn)
	}
	return
}

//显示七煞
func QiSha(wn int) (s string) {
	if wn == 0 ||
		wn == 1 ||
		wn == 8 ||
		wn == 14 ||
		wn == 22 ||
		wn == 23 ||
		wn == 24 {
		s = fmt.Sprintf("\"%s\"为七煞之一\n", zeji.XingSu28[wn])
	}
	return
}
