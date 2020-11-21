package main

import (
	"errors"
	"fmt"
	"image"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mjl-/duit"
	"github.com/nongli/ccal"
	"github.com/nongli/dimu"
	"github.com/nongli/ganzhi"
	"github.com/nongli/lunar"
	"github.com/nongli/solar"
	"github.com/nongli/today"
	"github.com/nongli/utils"
	xlr "github.com/nongli/zeji"
)

var (
	T      = func() time.Time { return time.Now().Local() }
	status = &duit.Label{Text: "计算时间范围:1601~3498\n使用方法:\n输入农历对应年份数字\n生肖可输入拼音\n闰月:y表示当月是闰月 n表示当月非闰月\n"}
	Ly     = &duit.Field{} //年
	Lm     = &duit.Field{} //月
	Ld     = &duit.Field{} //日
	Lh     = &duit.Field{} //时辰
	Lsx    = &duit.Field{} //生肖
	Lmb    = &duit.Field{} //闰月
	aliasM string          //是否闰月别名
	//生肖
	aliasshu, aliasniu, aliashu, aliastu, aliaslong, aliasshe string
	aliasma, aliasyang, aliashou, aliasji, aliasgou, aliaszhu string
	aliaslongf, aliasmaf, aliasjif, aliaszhuf                 string
	//绝日
	jlr string
)

func run() {
	dui, err := duit.NewDUI("小六壬农历择吉", nil)
	if err != nil {
		log.Fatalf("择吉: %s\n", err)
	}

	//纪年信息
	basicInfo := &duit.Button{
		Text:     "基础信息",
		Colorset: &dui.Primary,
		Click: func() (e duit.Event) {
			year := Ly.Text //string类型
			month := Lm.Text
			day := Ld.Text
			hour := Lh.Text
			sx := Lsx.Text
			leapm := Lmb.Text

			//类型转换
			y, m, d, h, inputb, err := String2Int(year, month, day, hour)
			if err != nil {
				log.Fatal(err)
			}
			mb, err := leapBool(leapm)
			if err != nil {
				log.Fatal(err)
			}
			info := ymd(y, m, d, h, sx, inputb, mb)
			status.Text = info
			dui.MarkLayout(status)
			return
		},
	}
	//地母经
	motherEarth := &duit.Button{
		Text: "地母经",
		Click: func() (e duit.Event) {
			year := Ly.Text
			y, _, _, _, _, err := String2Int(year, "6", "9", "3")
			if err != nil {
				log.Fatal(err)
			}
			err, _, _, g, _ := ccal.Input(y, 3, 1, 1, "猴", false) //这里月份要在立春之后
			if err != nil {
				log.Fatal(err)
			}
			dmg := g.YearGan
			dmz := g.YearZhi
			infodmj := dimu.DimuInfo(dmg, dmz)
			info := fmt.Sprintf("%s%s年地母经:\n%s\n", ganzhi.Gan[dmg], ganzhi.Zhi[dmz], infodmj)
			status.Text = info
			dui.MarkLayout(status)
			return
		},
	}
	//择吉信息
	auspicious := &duit.Button{
		Text: "择吉信息",
		Click: func() (e duit.Event) {
			year := Ly.Text //string类型
			month := Lm.Text
			day := Ld.Text
			hour := Lh.Text
			sx := Lsx.Text
			leapm := Lmb.Text

			y, m, d, h, inputb, err := String2Int(year, month, day, hour)
			if err != nil {
				log.Fatal(err)
			}

			mb, err := leapBool(leapm)
			if err != nil {
				log.Fatal(err)
			}
			info := aus(y, m, d, h, sx, inputb, mb)
			status.Text = info
			dui.MarkLayout(status)
			return
		},
	}
	//今日信息
	todayInfo := &duit.Button{
		Text: "今日信息",
		Click: func() (e duit.Event) {
			info := day()
			status.Text = info
			dui.MarkLayout(status)
			return
		},
	}
	//二十四节气列表
	j24Info := &duit.Button{
		Text: "二十四节气",
		Click: func() (e duit.Event) {
			year := Ly.Text
			y, m, d, h, inputb, err := String2Int(year, "3", "6", "9")
			if err != nil {
				log.Fatal(err)
			}
			info := j24(y, m, d, h, inputb)
			status.Text = info
			dui.MarkLayout(status)
			return
		},
	}
	//农历月历表
	listDayInfo := &duit.Button{
		Text: "农历月历表",
		Click: func() (e duit.Event) {
			year := Ly.Text //string类型
			month := Lm.Text
			day := Ld.Text
			hour := Lh.Text
			sx := Lsx.Text
			leapm := Lmb.Text

			y, m, d, h, inputb, err := String2Int(year, month, day, hour)
			if err != nil {
				log.Fatal(err)
			}

			mb, err := leapBool(leapm)
			if err != nil {
				log.Fatal(err)
			}
			info := listDay(y, m, d, h, sx, inputb, mb)
			status.Text = info
			dui.MarkLayout(status)
			return
		},
	}
	//主界面
	dui.Top.UI = &duit.Box{
		//MaxWidth: 800, //输入框宽度
		//Width:   870,
		Padding: duit.SpaceXY(6, 4),
		Margin:  image.Pt(6, 4),
		Kids: duit.NewKids(
			status,
			&duit.Grid{
				Columns: 2,
				Padding: []duit.Space{
					{Right: 6, Top: 4, Bottom: 4},
					{Left: 6, Top: 4, Bottom: 4},
				},
				Valign: []duit.Valign{duit.ValignMiddle, duit.ValignMiddle},
				Kids: duit.NewKids(
					&duit.Label{Text: "农历年"},
					Ly,
					&duit.Label{Text: "农历月"},
					Lm,
					&duit.Label{Text: "农历日"},
					Ld,
					&duit.Label{Text: "时辰"},
					Lh,
					&duit.Label{Text: "生肖"},
					Lsx,
					&duit.Label{Text: "闰月y/n"},
					Lmb,
				),
			},
			basicInfo,   //纪年信息
			motherEarth, //地母经信息
			auspicious,  //择吉信息
			todayInfo,   //今日信息
			j24Info,     //二十四节气
			listDayInfo, //月历表
		),
	}
	//第一次绘制整个用户界面
	dui.Render()

	//主循环
	for {
		//监听两个chan
		select {
		case e := <-dui.Inputs:
			dui.Input(e)
		case warn, ok := <-dui.Error:
			if !ok {
				return
			}
			log.Printf("duit: %s\n", warn)
		}
	}
}

//基础函数
func ymd(y, m, d, h int, sx string, inputb, mb bool) (Text string) {
	switch inputb {
	case true:
		sx = shengxiao(sx)
		err, s, l, g, _ := ccal.Input(y, m, d, h, sx, mb)
		if err != nil {
			//log.Fatal(err)
			Text = fmt.Sprintf("%v\n", err)
		}

		if l.Leapmb == true {
			aliasM = "是"
		} else {
			aliasM = "否"
		}
		solarinfo := fmt.Sprintf("阳历纪年: %d年-%d月-%d日-周%s\n", s.SYear, s.SMonth, s.SDay, s.SWeek)
		lunarinfo := fmt.Sprintf("农历纪年: %d年%s月(%s)%s %s时(%d时)\n本年是否有闰月:%s 闰%d月\n",
			l.LYear, lunar.Ymc[l.LMonth-1], l.LYdxs, lunar.Rmc[l.LDay-1], l.LaliasHour, l.LHour, aliasM, l.LeapMonth)
		gzinfo := fmt.Sprintf("干支纪年: %s%s年-%s月-%s%s日-%s时\n\n",
			g.YearGanM, g.YearZhiM, g.MonthGanZhiM, g.DayGanM, g.DayZhiM, g.HourGanZhiM)

		//杨公祭日
		yginfo := yg13(l.LMonth, l.LDay)

		Text = solarinfo + lunarinfo + gzinfo + yginfo
	case false:
		//log.Fatal("数字输入错\n")
		Text = fmt.Sprintf("数字输入错\n")
		os.Exit(1)
	}

	return
}

//择吉函数
func aus(y, m, d, h int, sx string, inputb, mb bool) (Text string) {
	switch inputb {
	case true:
		sx = shengxiao(sx)
		err, s, l, g, jq := ccal.Input(y, m, d, h, sx, mb)
		if err != nil {
			log.Fatal(err)
		}
		if l.Leapmb == true {
			aliasM = "是"
		} else {
			aliasM = "否"
		}
		//纪年信息
		solarinfo := fmt.Sprintf("阳历纪年: %d年-%d月-%d日-周%s\n",
			s.SYear, s.SMonth, s.SDay, s.SWeek)
		lunarinfo := fmt.Sprintf("农历纪年: %d年%s月(%s)%s %s时(%d时)\n本年是否有闰月: %s-->闰%d月\n",
			l.LYear, lunar.Ymc[l.LMonth-1], l.LYdxs, lunar.Rmc[l.LDay-1], l.LaliasHour, l.LHour, aliasM, l.LeapMonth)
		gzinfo := fmt.Sprintf("干支纪年: %s%s年-%s月-%s%s日-%s时\n\n",
			g.YearGanM, g.YearZhiM, g.MonthGanZhiM, g.DayGanM, g.DayZhiM, g.HourGanZhiM)

		//值宿信息
		dgz := fmt.Sprintf("%s%s", g.DayGanM, g.DayZhiM)
		wd := int(s.SolarDayT.Weekday()) //周几
		ygz := fmt.Sprintf("%s%s", g.YearGanM, g.YearZhiM)
		m := l.LMonth
		d := l.LDay
		h := l.LHour
		stars := xlr.NewStars(wd, dgz)
		ws := stars.Name //值宿名称
		wn := stars.Week //值宿当日周几　0为周日

		winfo := fmt.Sprintf("周%s 值宿:\"%s\"\n", solar.Zhou[wn], ws)
		info := stars.Info //当日值宿信息

		total := xlr.NewTotal(ygz, m, d, h)
		b1 := total.YiPan()
		b2 := total.ErPan()
		b3 := total.SanPan()
		zeji := stars.GoodNumberDay(b1, b2, b3)

		//本月吉干信息
		jg1, jg2, jg3 := xlr.JiGanNumber(m)
		name := fmt.Sprintf("本月吉干:%s-%s-%s", ganzhi.Gan[jg1], ganzhi.Gan[jg2], ganzhi.Gan[jg3])

		//判断当日是否为七煞日

		//农历月份吉干
		lunarmjd := jq.LunarmJd
		ydx := l.LYdx
		jgmap, qsarr := xlr.FindMonthJG(lunarmjd, m, ydx)
		//去除金神七煞的吉干map
		newjgmap := xlr.DelQs(jgmap, qsarr)
		//去除生肖相冲的吉干map
		good := xlr.DelSX(newjgmap, sx)
		//本月吉干数组
		goodarr := xlr.GoodJG(good)
		jg := fmt.Sprintln(goodarr)
		//信息显示到UI界面
		Text = (solarinfo + lunarinfo + gzinfo + winfo + info + zeji + name + "\n本月吉日:" + jg)
	case false:
		Text = ("数字输入错误，系统退出\n")
		os.Exit(0)
	}
	return
}

//自动显示阳历当日择吉内容
func day() (Text string) {
	expectInfo, err := today.NewExpectInfo()
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

	h24 := T().Hour()
	h := utils.Conv24Hto12H(h24)
	sx := "猴"

	if leapM != 0 && leapB == true {
		err, s, l, g, jq := ccal.Input(leapY, leapM, expectLeapD, h, sx, leapB)
		if err != nil {
			log.Fatal(err)
		}
		if l.Leapmb == true {
			aliasM = "是"
		} else {
			aliasM = "否"
		}
		//纪年信息
		solarinfo := fmt.Sprintf("阳历纪年: %d年-%d月-%d日-周%s\n", s.SYear, s.SMonth, s.SDay, s.SWeek)
		lunarinfo := fmt.Sprintf("农历纪年: %d年%s月(%s)%s %s时(%d时)\n本年是否有闰月: %s-->闰%d月\n",
			l.LYear, lunar.Ymc[l.LMonth-1], l.LYdxs, lunar.Rmc[l.LDay-1], l.LaliasHour, l.LHour, aliasM, l.LeapMonth)
		gzinfo := fmt.Sprintf("干支纪年: %s%s年-%s月-%s%s日-%s时\n\n",
			g.YearGanM, g.YearZhiM, g.MonthGanZhiM, g.DayGanM, g.DayZhiM, g.HourGanZhiM)
		//值宿信息
		dgz := fmt.Sprintf("%s%s", g.DayGanM, g.DayZhiM)
		wd := int(s.SolarDayT.Weekday()) //周几
		ygz := fmt.Sprintf("%s%s", g.YearGanM, g.YearZhiM)
		m := l.LMonth
		d := l.LDay
		h := l.LHour
		stars := xlr.NewStars(wd, dgz)
		ws := stars.Name //值宿名称
		wn := stars.Week //值宿当日周几　0为周日

		winfo := fmt.Sprintf("周%s 值宿:\"%s\"\n", solar.Zhou[wn], ws)
		info := stars.Info //当日值宿信息

		total := xlr.NewTotal(ygz, m, d, h)
		b1 := total.YiPan()
		b2 := total.ErPan()
		b3 := total.SanPan()
		zeji := stars.GoodNumberDay(b1, b2, b3)

		//本月吉干信息
		jg1, jg2, jg3 := xlr.JiGanNumber(m)
		//number := fmt.Sprintf("吉干数字:%d-%d-%d\n", jg1, jg2, jg3)
		name := fmt.Sprintf("本月吉干:%s-%s-%s", ganzhi.Gan[jg1], ganzhi.Gan[jg2], ganzhi.Gan[jg3])

		//判断当日是否为七煞日

		//农历月份吉干
		lunarmjd := jq.LunarmJd
		ydx := l.LYdx
		jgmap, qsarr := xlr.FindMonthJG(lunarmjd, m, ydx)
		//去除金神七煞的吉干map
		newjgmap := xlr.DelQs(jgmap, qsarr)
		//去除生肖相冲的吉干map
		good := xlr.DelSX(newjgmap, sx)
		//本月吉干数组
		goodarr := xlr.GoodJG(good)
		jg := fmt.Sprintln(goodarr)
		//信息显示到UI界面
		Text = (solarinfo + lunarinfo + gzinfo + winfo + info + zeji + name + "\n本月吉日:" + jg)

	} else if normalM != 0 && normalB == false {
		err, s, l, g, jq := ccal.Input(normalY, normalM, expectD, h, sx, normalB)
		if err != nil {
			log.Fatal(err)
		}
		if l.Leapmb == true {
			aliasM = "是"
		} else {
			aliasM = "否"
		}
		//纪年信息
		solarinfo := fmt.Sprintf("阳历纪年: %d年-%d月-%d日-周%s\n", s.SYear, s.SMonth, s.SDay, s.SWeek)
		lunarinfo := fmt.Sprintf("农历纪年: %d年%s月(%s)%s %s时(%d时)\n本年是否有闰月: %s-->闰%d月\n",
			l.LYear, lunar.Ymc[l.LMonth-1], l.LYdxs, lunar.Rmc[l.LDay-1], l.LaliasHour, l.LHour, aliasM, l.LeapMonth)
		gzinfo := fmt.Sprintf("干支纪年: %s%s年-%s月-%s%s日-%s时\n\n",
			g.YearGanM, g.YearZhiM, g.MonthGanZhiM, g.DayGanM, g.DayZhiM, g.HourGanZhiM)

		//值宿信息
		dgz := fmt.Sprintf("%s%s", g.DayGanM, g.DayZhiM)
		wd := int(s.SolarDayT.Weekday()) //周几
		ygz := fmt.Sprintf("%s%s", g.YearGanM, g.YearZhiM)
		m := l.LMonth
		d := l.LDay
		h := l.LHour
		stars := xlr.NewStars(wd, dgz)
		ws := stars.Name //值宿名称
		wn := stars.Week //值宿当日周几　0为周日

		winfo := fmt.Sprintf("周%s 值宿:\"%s\"\n", solar.Zhou[wn], ws)
		info := stars.Info //当日值宿信息

		total := xlr.NewTotal(ygz, m, d, h)
		b1 := total.YiPan()
		b2 := total.ErPan()
		b3 := total.SanPan()
		zeji := stars.GoodNumberDay(b1, b2, b3)

		//本月吉干信息
		jg1, jg2, jg3 := xlr.JiGanNumber(m)
		//number := fmt.Sprintf("吉干数字:%d-%d-%d\n", jg1, jg2, jg3)
		name := fmt.Sprintf("本月吉干:%s-%s-%s", ganzhi.Gan[jg1], ganzhi.Gan[jg2], ganzhi.Gan[jg3])

		//判断当日是否为七煞日

		//农历月份吉干
		lunarmjd := jq.LunarmJd
		ydx := l.LYdx
		jgmap, qsarr := xlr.FindMonthJG(lunarmjd, m, ydx)
		//去除金神七煞的吉干map
		newjgmap := xlr.DelQs(jgmap, qsarr)
		//去除生肖相冲的吉干map
		good := xlr.DelSX(newjgmap, sx)
		//本月吉干数组
		goodarr := xlr.GoodJG(good)
		jg := fmt.Sprintln(goodarr)
		//信息显示到UI界面
		Text = (solarinfo + lunarinfo + gzinfo + winfo + info + zeji + name + "\n本月吉日:" + jg)
	}
	return
}

//24节气
func j24(y, m, d, h int, inputb bool) (Text string) {

	switch inputb {
	case true:
		err, _, _, _, jq := ccal.Input(y, m, d, h, "猴", false)
		if err != nil {
			log.Fatal(err)
		}
		jq24 := solar.ShowJieqi24(jq.Jqt, jq.Jq11t)
		//信息显示到UI界面
		Text = (jq24[0] + jq24[1] + jq24[2] + jq24[3] + jq24[4] + jq24[5] +
			jq24[6] + jq24[7] + jq24[8] + jq24[9] + jq24[10] + jq24[11] +
			jq24[12] + jq24[13] + jq24[14] + jq24[15] + jq24[16] + jq24[17] +
			jq24[18] + jq24[19] + jq24[20] + jq24[21] + jq24[22] + jq24[23] +
			jq24[24] + jq24[25] + jq24[26])
	case false:
		Text = "年份数字输入错误\n"
	}
	return
}

//农历月历表
func listDay(y, m, d, h int, sx string, inputb, mb bool) (Text string) {
	switch inputb {
	case true:
		err, _, l, _, jq := ccal.Input(y, m, d, h, sx, mb)
		if err != nil {
			Text = fmt.Sprintf("%v\n", err)
		}
		lunarmjd := jq.LunarmJd
		ydx := l.LYdx
		days := lunar.ListLunarDay(lunarmjd, ydx)
		x := len(days)

		n := fmt.Sprintf("\n") //自动换行
		if x == 29 {
			Text = (n + days[0] + days[1] + days[2] + days[3] + days[4] + days[5] + days[6] + n +
				n + days[7] + days[8] + days[9] + days[10] + days[11] + days[12] + days[13] + n +
				n + days[14] + days[15] + days[16] + days[17] + days[18] + days[19] + days[20] + n +
				n + days[21] + days[22] + days[23] + days[24] + days[2] + days[26] + days[27] + n +
				n + days[28])

		}
		if x == 30 {
			Text = (n + days[0] + days[1] + days[2] + days[3] + days[4] + days[5] + days[6] + n +
				n + days[7] + days[8] + days[9] + days[10] + days[11] + days[12] + days[13] + n +
				n + days[14] + days[15] + days[16] + days[17] + days[18] + days[19] + days[20] + n +
				n + days[21] + days[22] + days[23] + days[24] + days[2] + days[26] + days[27] + n +
				n + days[28] + days[29])
		}
	case false:
		Text = "数字输入错误\n"
	}
	return
}

//字符串类型转int　返回值为真表示输入数字正常
func String2Int(year, month, day, hour string) (y, m, d, h int, inputb bool, err error) {
	y, _ = strconv.Atoi(year)
	m, _ = strconv.Atoi(month)
	d, _ = strconv.Atoi(day)
	h, _ = strconv.Atoi(hour)
	inputb, err = dateBool(y, m, d, h)

	return
}

//判断输入的数字
func dateBool(year, month, day, hour int) (dateB bool, err error) {

	if (year > 1600 && year < 3499) &&
		(month >= 1 && month <= 12) &&
		(day >= 1 && day <= 30) &&
		(hour >= 1 && hour <= 12) {
		dateB = true
	} else {
		err = errors.New("年份时间范围1601到3498")
		dateB = false
	}
	return
}

//判断输入是不是闰月
func leapBool(leapm string) (lt bool, err error) {

	sl := strings.ToLower(leapm)
	slby := strings.EqualFold(sl, "y")
	slby1 := strings.EqualFold(sl, "yes")
	slbn := strings.EqualFold(sl, "n")
	slbn1 := strings.EqualFold(sl, "no")

	if slby == true || slby1 == true {
		lt = true
	} else if slbn == false || slbn1 == false {
		lt = false
	}

	if leapm != "yes" && leapm != "y" &&
		leapm != "no" && leapm != "n" {
		err = errors.New("闰月判断值输入错误")
	}
	return
}

//生肖判断　可以输入生肖的拼音或者汉字支持繁体
func shengxiao(s string) (sxname string) {
	lows := strings.ToLower(s) //转为小写

	if strings.EqualFold(s, "鼠") || strings.EqualFold(lows, "shu") {
		sxname = "鼠"
	}
	if strings.EqualFold(s, "牛") || strings.EqualFold(lows, "niu") {
		sxname = "牛"
	}
	if strings.EqualFold(s, "虎") || strings.EqualFold(lows, "hu") {
		sxname = "虎"
	}
	if strings.EqualFold(s, "兔") || strings.EqualFold(lows, "tu") {
		sxname = "兔"
	}
	if strings.EqualFold(s, "龙") || strings.EqualFold(lows, "long") || strings.EqualFold(s, "龍") {
		sxname = "龙"
	}
	if strings.EqualFold(s, "蛇") || strings.EqualFold(lows, "she") {
		sxname = "蛇"
	}
	if strings.EqualFold(s, "马") || strings.EqualFold(lows, "ma") || strings.EqualFold(s, "馬") {
		sxname = "马"
	}
	if strings.EqualFold(s, "羊") || strings.EqualFold(lows, "yang") {
		sxname = "羊"
	}
	if strings.EqualFold(s, "猴") || strings.EqualFold(lows, "hou") {
		sxname = "猴"
	}
	if strings.EqualFold(s, "鸡") || strings.EqualFold(lows, "ji") || strings.EqualFold(s, "雞") {
		sxname = "鸡"
	}
	if strings.EqualFold(s, "狗") || strings.EqualFold(lows, "gou") {
		sxname = "狗"
	}
	if strings.EqualFold(s, "猪") || strings.EqualFold(lows, "zhu") || strings.EqualFold(s, "豬") {
		sxname = "猪"
	}
	return
}

//杨公十三祭
func yg13(m, d int) (info string) {

	if m == 1 && d == 13 {
		info = "正月十三杨公忌日\n"
	}
	if m == 2 && d == 11 {
		info = "二月十一杨公忌日\n"
	}
	if m == 3 && d == 9 {
		info = "三月初九杨公忌日\n"
	}
	if m == 4 && d == 7 {
		info = "四月初七杨公忌日\n"
	}
	if m == 5 && d == 5 {
		info = "五月初五杨公忌日\n"
	}
	if m == 6 && d == 3 {
		info = "六月初三杨公忌日\n"
	}
	if m == 7 && d == 1 {
		info = "七月初一杨公忌日\n"
	}
	if m == 7 && d == 29 {
		info = "七月二十九杨公忌日\n"
	}
	if m == 8 && d == 27 {
		info = "八月二十七杨公忌日\n"
	}
	if m == 9 && d == 25 {
		info = "九月二十五杨公忌日\n"
	}
	if m == 10 && d == 23 {
		info = "十月二十三杨公忌日\n"
	}
	if m == 11 && d == 21 {
		info = "十一月二十一杨公忌日\n"
	}
	if m == 12 && d == 19 {
		info = "十二月十九杨公忌日\n"
	}
	return
}

/*
//打印当时是否为七煞日
func PrintQS(qsb bool, qsn int) (info string) {

	if qsb == true {
		info = QiSha(qsn)
	}
	return
}
*/
/* //显示七煞
func QiSha(wn int) (s string) {
	if wn == 0 ||
		wn == 1 ||
		wn == 8 ||
		wn == 14 ||
		wn == 22 ||
		wn == 23 ||
		wn == 24 {
		s = fmt.Sprintf("\n\"%s\"为七煞之一\n", zeji.XingSu28[wn])
	}
	return
} */
