package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nongli/ccal"
	"github.com/nongli/ganzhi"
	"github.com/nongli/lunar"
	"github.com/nongli/solar"
	"github.com/nongli/today"
	"github.com/nongli/utils"
	"github.com/nongli/zeji"
	"github.com/ying32/govcl/vcl"
)

var T = time.Now().Local()

//自动显示阳历当日择吉内容
func (f *TForm1) OnButton8Click(sender vcl.IObject) {

	expectInfo, err := today.FindLunarMD()
	if err != nil {
		log.Fatal("时间异常\n", err)
	}

	//润月
	leapY := expectInfo.LeapY
	leapM := expectInfo.LeapM
	expectLeapD := expectInfo.ExpectleapD
	leapB := expectInfo.LeapB

	//正常月
	normalY := expectInfo.NormalY
	normalM := expectInfo.NormalM
	expectD := expectInfo.ExpectD
	normalB := expectInfo.NormalB

	h24 := T.Hour()
	h := utils.Conv24Hto12H(h24)
	sx := "猴"

	if leapM != 0 && leapB == true {
		err, s, l, g, jq := ccal.Input(leapY, leapM, expectLeapD, h, sx, leapB)
		if err != nil {
			log.Fatal(err)
		}
		if l.Leapmb == true {
			lm = "是"
		} else {
			lm = "否"
		}
		//纪年信息
		solarinfo := fmt.Sprintf("阳历纪年: %d年-%d月-%d日-周%s-阳历时间 %d:%d\n", s.SYear, s.SMonth, s.SDay, s.SWeek, T.Hour(), T.Minute())
		lunarinfo := fmt.Sprintf("农历纪年: %d年%s月(%s)%s %s时(%d时)\n本年是否有闰月: %s-->闰%d月\n",
			l.LYear, lunar.Ymc[l.LMonth-1], l.LYdxs, lunar.Rmc[l.LDay-1], l.LaliasHour, l.LHour, lm, l.LeapMonth)
		gzinfo := fmt.Sprintf("干支纪年: %s%s年-%s月-%s%s日-%s时\n\n",
			g.YearGanM, g.YearZhiM, g.MonthGanZhiM, g.DayGanM, g.DayZhiM, g.HourGanZhiM)

		//值宿信息
		iqs := zeji.ZhiSu(s, g)
		ws := iqs.StarNames //值宿名称
		wn := iqs.Week      //值宿当日周几　0为周日
		winfo := fmt.Sprintf("周%s 值宿:\"%s\"\n", solar.Zhou[wn], ws)
		zhisuInfo := fmt.Sprintf(iqs.ZhiSu) //当日值宿信息

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
		nx := zeji.AllNumber(yearZhi, leapM, expectLeapD, h)
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

		//奇門
		qmdj := qm(s.SolarDayT, g, jq)

		//信息显示到UI界面
		vcl.ShowMessage(solarinfo + lunarinfo + gzinfo + qmdj + winfo + zhisuInfo + isQiSha + number + name + result + jlr + listJg)
	} else if normalM != 0 && normalB == false {
		err, s, l, g, jq := ccal.Input(normalY, normalM, expectD, h, sx, normalB)
		if err != nil {
			log.Fatal(err)
		}
		if l.Leapmb == true {
			lm = "是"
		} else {
			lm = "否"
		}
		//纪年信息
		solarinfo := fmt.Sprintf("阳历纪年: %d年-%d月-%d日-周%s-阳历时间 %d:%d\n", s.SYear, s.SMonth, s.SDay, s.SWeek, T.Hour(), T.Minute())
		lunarinfo := fmt.Sprintf("农历纪年: %d年%s月(%s)%s %s时(%d时)\n本年是否有闰月: %s-->闰%d月\n",
			l.LYear, lunar.Ymc[l.LMonth-1], l.LYdxs, lunar.Rmc[l.LDay-1], l.LaliasHour, l.LHour, lm, l.LeapMonth)
		gzinfo := fmt.Sprintf("干支纪年: %s%s年-%s月-%s%s日-%s时\n\n",
			g.YearGanM, g.YearZhiM, g.MonthGanZhiM, g.DayGanM, g.DayZhiM, g.HourGanZhiM)

		//值宿信息
		iqs := zeji.ZhiSu(s, g)
		ws := iqs.StarNames //值宿名称
		wn := iqs.Week      //值宿当日周几　0为周日
		winfo := fmt.Sprintf("周%s 值宿:\"%s\"\n", solar.Zhou[wn], ws)
		zhisuInfo := fmt.Sprintf(iqs.ZhiSu) //当日值宿信息

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
		nx := zeji.AllNumber(yearZhi, normalM, expectD, h)
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

		//奇門
		qmdj := qm(s.SolarDayT, g, jq)

		//信息显示到UI界面
		vcl.ShowMessage(solarinfo + lunarinfo + gzinfo + qmdj + winfo + zhisuInfo + isQiSha + number + name + result + jlr + listJg)
	}
}
