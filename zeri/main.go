package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Aquarian-Age/nongli/ccal"
	"github.com/Aquarian-Age/nongli/lunar"
	ganzhi "github.com/Aquarian-Age/ts/gz"
	ts "github.com/Aquarian-Age/ts/tongshu"
	"github.com/mjl-/duit"
)

var (
	T      = time.Now().Local()
	status = &duit.Label{Text: "农历择日\n\n使用方法: 输入农历的年月日时辰闰月:\n比如农历2020年六月六日十一时\n比如2020060611f\n"}
	get    = &duit.Field{} //年
)

type Input struct {
	Y int  //年數字
	M int  //月數字
	D int  //日數字
	H int  //時辰數字　1子時　2丑時...１２亥時
	B bool //閏月判斷　f非閏月　t當月爲閏月
}

//////////////////////通书
type TongShu struct {
	syear  int       //阳历年
	smonth int       //阳历月
	sday   int       //阳历日
	sweek  string    //阳历周几
	stime  time.Time //当日阳历时间戳　精确到日

	lyear  int //农历年
	lmonth int //农历月
	lday   int //农历日

	aliasyg    string //农历年干
	aliasyz    string //年支
	aliasygz   string //年干支
	aliasmgz   string //月干支
	aliasday   string //日别名　一　二　廿二
	aliasmonth string //月别名 正　二　三
	dgz        string //日干支
	daygan     string //日干
	hourgz     string //时干支
	hourz      string //时支
	leapmb     bool   //闰月
	ydx        string //月大小数组
	aliasHour  string //时辰(地支)
	hour       int    //时子:1 丑:2．．．
	leapmonth  int    //闰月
}

var jz60 = ganzhi.MakeJZ60()
var i ts.ZRYLImpl
var tos *TongShu

////////////////////////

func init() {
	os.Setenv("font", "/opt/fonts/unifont/unifont.font")
}

func main() {
	runMain()
}

func runMain() {
	dui, err := duit.NewDUI("择日", nil)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	//纪年信息
	basicInfo := &duit.Button{
		Text:     "信息",
		Colorset: &dui.Primary,
		Click: func() (e duit.Event) {
			years := get.Text //string类型
			sx := "猴"

			//类型转换
			infos := conv(years)
			y := infos.Y
			m := infos.M
			d := infos.D
			h := infos.H
			b := infos.B

			//日期干支信息
			err, s, l, g, _ := ccal.Input(y, m, d, h, sx, b)
			if err != nil {
				log.Fatal("ccal:", err)
			}

			hgz := convHourZhi(g.HourGanZhiM)
			tos = &TongShu{
				syear:      s.SYear,
				smonth:     s.SMonth,
				sday:       s.SDay,
				sweek:      s.SWeek,
				stime:      s.SolarDayT,
				lyear:      l.LYear,
				lmonth:     l.LMonth,
				lday:       l.LDay,
				aliasyg:    g.YearGanM,
				aliasyz:    g.YearZhiM,
				aliasygz:   fmt.Sprintf("%s%s", g.YearGanM, g.YearZhiM),
				aliasmgz:   g.MonthGanZhiM,
				aliasday:   convRmc(l.LDay),
				aliasmonth: convYmc(l.LMonth),
				dgz:        fmt.Sprintf("%s%s", g.DayGanM, g.DayZhiM),
				daygan:     g.DayGanM,
				hourgz:     g.HourGanZhiM,
				hourz:      hgz,
				leapmb:     l.Leapmb,
				ydx:        l.LYdxs,
				aliasHour:  l.LaliasHour,
				hour:       l.LHour,
				leapmonth:  l.LeapMonth,
			}
			//
			jinian := tos.Lunar()              //基本纪年信息
			nayin := tos.NaYin()               //干支纳因
			yeartab := tos.YearTab()           //协纪辩方 年表
			djc, jcb := tos.JCM()              //协纪辩方 建除十二神煞(日)
			monthtab := tos.MonthTab(djc, jcb) //协纪辩方 月表
			daytab := tos.DayTab()             //协纪辩方 日表
			bianwei := tos.BianWei()           //协纪辩方 辩伪+其他
			///////////////////////
			//UI显示
			status.Text = jinian + nayin + yeartab + monthtab + daytab + bianwei
			dui.MarkLayout(status)
			return
		},
	}

	//ui
	dui.Top.UI = &duit.Box{
		Padding: duit.SpaceXY(6, 4), //从窗口插入
		Margin:  image.Pt(6, 4),     //此框中kids之间的空间
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
					get,
				),
			},
			basicInfo, //信息
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

//基础历法
func (tos *TongShu) Lunar() string {
	//闰月汉字表示
	var aliasM string
	if tos.leapmb == true {
		aliasM = "是"
	} else {
		aliasM = "否"
	}
	//基础信息
	fmt.Printf("阳历纪年: %d年-%d月-%d日-周%s\n", tos.syear, tos.smonth, tos.sday, tos.sweek)
	fmt.Printf("农历纪年: %d年%s月(%s)%s %s时(%d时)\n本年是否有闰月:%s 闰%d月\n",
		tos.lyear, lunar.Ymc[tos.lmonth-1], tos.ydx, lunar.Rmc[tos.lday-1],
		tos.aliasHour, tos.hour, aliasM, tos.leapmonth)
	fmt.Printf("干支纪年: %s%s年-%s月-%s日-%s时\n",
		tos.aliasyg, tos.aliasyz, tos.aliasmgz, tos.dgz, tos.hourgz)
	///gui
	s1 := fmt.Sprintf("阳历纪年: %d年-%d月-%d日-周%s\n", tos.syear, tos.smonth, tos.sday, tos.sweek)
	s2 := fmt.Sprintf("农历纪年: %d年%s月(%s)%s %s时(%d时)\n本年是否有闰月:%s 闰%d月\n",
		tos.lyear, lunar.Ymc[tos.lmonth-1], tos.ydx, lunar.Rmc[tos.lday-1],
		tos.aliasHour, tos.hour, aliasM, tos.leapmonth)
	s3 := fmt.Sprintf("干支纪年: %s%s年-%s月-%s日-%s时\n",
		tos.aliasyg, tos.aliasyz, tos.aliasmgz, tos.dgz, tos.hourgz)

	return s1 + s2 + s3
}

//干支纳音
func (tos *TongShu) NaYin() string {
	ygzny := ganzhi.GZ纳音(tos.aliasygz)
	mgzny := ganzhi.GZ纳音(tos.aliasmgz)
	dgzny := ganzhi.GZ纳音(tos.dgz)
	hgzny := ganzhi.GZ纳音(tos.hourgz)
	return fmt.Sprintf("干支纳音: %s %s %s %s\n",
		ygzny[tos.aliasygz], mgzny[tos.aliasmgz], dgzny[tos.dgz], hgzny[tos.hourgz])
}

//协纪辩方 年表
func (tos *TongShu) YearTab() string {
	dgz := tos.dgz
	ygz := tos.aliasygz
	mgz := tos.aliasmgz
	yz := tos.aliasyz
	aliasmonth := tos.aliasmonth
	i = &ts.ZRYL{
		YGZ:          ygz,
		DGZ:          dgz,
		MGZ:          mgz,
		AliasYearZhi: yz,
		AliasMonth:   aliasmonth,
	}
	yeartab := i.XJBF年表(jz60)
	return yeartab
}

//协纪辩方 建除(月日论)
func (tos *TongShu) JCM() (string, bool) {
	m := tos.aliasmonth //农历月别名
	dgz := tos.dgz      //日干支
	i = &ts.ZRYL{
		AliasMonth: m,
		DGZ:        dgz,
	}
	return i.JC12M()
}

//协纪辩方 月表
func (tos *TongShu) MonthTab(djc string, jcb bool) string {
	ygz := tos.aliasygz
	m := tos.aliasmonth
	ly := tos.lyear
	st := tos.stime
	mgz := tos.aliasmgz
	dgz := tos.dgz
	lday := tos.lday
	i = &ts.ZRYL{
		YGZ:        ygz,
		AliasMonth: m,
		Lyear:      ly,
		SolarT:     st,
		MGZ:        mgz,
		DGZ:        dgz,
		Lday:       lday,
	}

	return i.XJBF月表(djc, jcb)
}

//协纪辩方 日表
func (tos *TongShu) DayTab() string {
	dgz := tos.dgz
	hgz := tos.hourgz
	i = &ts.ZRYL{
		DGZ: dgz,
		HGZ: hgz,
	}
	return i.XJBF日表(jz60)
}

//协纪辩方书 辩伪+其他
func (tos *TongShu) BianWei() string {
	ygz := tos.aliasygz
	mgz := tos.aliasmgz
	dgz := tos.dgz
	hgz := tos.hourgz
	guchen, guasu := ts.XJBF孤辰寡宿(ygz, mgz, dgz, hgz)
	taohua := ts.XCTH咸池桃花(ygz, mgz, dgz, hgz)
	return guchen + guasu + taohua
}

//字符串转数字
func conv(s string) *Input {
	rs := []rune(s)
	ys := string(rs[:4])  //年數字
	ms := string(rs[4:6]) //月數字
	ds := string(rs[6:8])
	hs := string(rs[8:10]) //時辰數字1子時　2丑時...１２亥時
	bs := string(rs[10:11])

	y, err := strconv.Atoi(ys)
	if err != nil {
		log.Fatal("年份時間解析:", err)
	}

	m, err := strconv.Atoi(ms)
	if err != nil {
		log.Fatal("月份時間解析:", err)
	}
	d, err := strconv.Atoi(ds)
	if err != nil {
		log.Fatal("日期時間解析:", err)
	}
	h, err := strconv.Atoi(hs)
	if err != nil {
		log.Fatal("時辰解析:", err)
	}
	b, err := strconv.ParseBool(bs)
	if err != nil {
		log.Fatal("閏月解析:", err)
	}

	//fmt.Sprintf("input: y:%v m:%v d:%v h:%v b:%t\n", y, m, d, h, b)
	return &Input{
		Y: y,
		D: d,
		M: m,
		H: h,
		B: b,
	}
}

//找出时辰地支
func convHourZhi(hourgz string) (alias string) {
	zhi := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
	for i := 0; i < len(zhi); i++ {
		if strings.ContainsAny(hourgz, zhi[i]) {
			alias = zhi[i]
			break
		}
	}
	return
}

//这里日名称用的是廿
func convRmc(n int) (alias string) {
	rmc := []string{
		"初一", "初二", "初三", "初四", "初五", "初六", "初七", "初八", "初九", "初十",
		"十一", "十二", "十三", "十四", "十五", "十六", "十七", "十八", "十九", "二十",
		"廿一", "廿二", "廿三", "廿四", "廿五", "廿六", "廿七", "廿八", "廿九", "三十"}
	for i := 0; i < len(rmc); i++ {
		if i+1 == n {
			alias = rmc[i]
			break
		}
	}
	return
}

//农历数字专汉字
func convYmc(n int) (alias string) {
	ymc := []string{"正", "二", "三", "四", "五", "六", "七", "八", "九", "十", "十一", "十二"}
	for i := 0; i < len(ymc); i++ {
		if i+1 == n {
			alias = ymc[i]
			break
		}
	}
	return
}
