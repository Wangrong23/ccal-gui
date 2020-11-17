package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/Aquarian-Age/nongli/ccal"
	"github.com/Aquarian-Age/nongli/dimu"
	"github.com/Aquarian-Age/nongli/lunar"
	"github.com/Aquarian-Age/nongli/solar"
	"github.com/Aquarian-Age/nongli/today"
	"github.com/Aquarian-Age/nongli/utils"
	"github.com/Aquarian-Age/nongli/zeji"
	ganzhi "github.com/Aquarian-Age/ts/gz"
	ts "github.com/Aquarian-Age/ts/tongshu"
)

var (
	T    = time.Now().Local() //当前时间
	jz60 = ganzhi.MakeJZ60()  //六十甲子
	i    ts.ZRYLImpl          //协纪辩方择日方法
	zr   *ZR
)

//择日信息
type ZR struct {
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

//协纪辩方书
type XJBF struct {
	NB string `json:"xjbfsNB"` //协纪辩方书 年表
	YB string `json:"xjbfsYB"` //协纪辩方书 月表
	RB string `json:"xjbfsRB"` //协纪辩方书 日表
	BW string `json:"xjbfsBW"` //协纪辩方书 辩伪
	RS string `json:"tsRSZL"`  //通书 日时总览
}

//今日信息
type Today struct {
	Ti string `json:"todayInfo"` //今日纪年信息
	About
}

//关于...
type About struct {
	Ccal string
	Data string
	Xlr  string
	Xjbf string
	Ck   string
	Me   string
}

//月将
type YJ struct {
	YjName   string `json:"yjName"`   //月将名称
	StarName string `json:"starName"` //十二星宫
}

//应答数据
type Resp struct {
	//JN
	JiNian   string   `json:"jinianInfo"`   //纪年信息
	Dmj      string   `json:"dmInfo"`       //地母经
	Jq       []string `json:"jqInfo"`       //24节气
	ListDay  []string `json:"listdayInfo"`  //农历月历表(农历初一开始)
	StarName string   `json:"starnameInfo"` //当日值宿名称(28宿)
	StarInfo string   `json:"starInfo"`     //值宿信息
	Zeji     string   `json:"zejiInfo"`     //当日择吉信息
	XJBF              //协纪辩方书
	//Today    string   `json:"todayInfo"` //阳历今日信息
	YJ //月将信息
}

//宜忌
func ccalyj(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("ccal.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("home.html")
		t.Execute(w, nil)
	} else {
		//解析表单
		r.ParseForm()
		//农历年
		ly, err := strconv.Atoi(r.Form["ly"][0])
		if err != nil {
			log.Fatalln("ly:", err)
		}
		//农历月
		lm, err := strconv.Atoi(r.Form["lm"][0])
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("农历: %d月\n", lm)
		ld, err := strconv.Atoi(r.Form["ld"][0])
		if err != nil {
			log.Fatalln(err)
		}
		//时辰 子时1 丑时2 寅时3...
		lh, err := strconv.Atoi(r.Form["lh"][0])
		if err != nil {
			log.Fatalln("时辰异常:", err)
		}
		//生肖
		sxs := r.Form["la"][0]
		sxn, err := strconv.Atoi(sxs)
		if err != nil {
			log.Fatal("生肖异常:", err)
		}
		var sx string
		var zhi = []string{"err", "鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}
		for index := 0; index <= 13; index++ {
			if sxn == index {
				sx = zhi[index]
				break
			}
		}
		//闰月
		sb := r.Form["le"][0]
		var leapb bool
		if strings.EqualFold(sb, "t") {
			leapb = true
		} else if strings.EqualFold(sb, "f") {
			leapb = false
		}
		/////////////////////////////ccal农历基本纪年信息
		_, s, l, g, jq := ccal.Input(ly, lm, ld, lh, sx, leapb)
		var aliasM string
		if l.Leapmb == true {
			aliasM = "是"
		} else {
			aliasM = "否"
		}
		solarinfo := fmt.Sprintf("阳历纪年: %d年-%d月-%d日-周%s\n", s.SYear, s.SMonth, s.SDay, s.SWeek)
		lunarinfo := fmt.Sprintf("农历纪年: %d年%s月(%s)%s %s时(%d时)<br />本年是否有闰月:%s闰%d月\n", l.LYear, lunar.Ymc[l.LMonth-1], l.LYdxs, lunar.Rmc[l.LDay-1], l.LaliasHour, l.LHour, aliasM, l.LeapMonth)
		gzinfo := fmt.Sprintf("干支纪年: %s%s年-%s月-%s%s日-%s时\n",
			g.YearGanM, g.YearZhiM, g.MonthGanZhiM, g.DayGanM, g.DayZhiM,
			g.HourGanZhiM)
		//纳音
		aliasygz := fmt.Sprintf("%s%s", g.YearGanM, g.YearZhiM)
		aliasmgz := g.MonthGanZhiM
		dgz := fmt.Sprintf("%s%s", g.DayGanM, g.DayZhiM)
		hourgz := g.HourGanZhiM

		ygzny := ganzhi.GZ纳音(aliasygz)
		mgzny := ganzhi.GZ纳音(aliasmgz)
		dgzny := ganzhi.GZ纳音(dgz)
		hgzny := ganzhi.GZ纳音(hourgz)
		nyinfo := fmt.Sprintf("干支纳音: %s %s %s %s\n", ygzny[aliasygz], mgzny[aliasmgz], dgzny[dgz], hgzny[hourgz])
		jinianinfo := solarinfo + "<br />" + lunarinfo + "<br />" + gzinfo + "<br />" + nyinfo

		/////////////////地母经
		dmg := g.YearGan
		dmz := g.YearZhi
		infodmj := dimu.DimuInfo(dmg, dmz)
		dmjinfo := fmt.Sprintf("%s", infodmj)

		///////////////24节气
		jq24info := solar.ShowJieqi24(jq.Jqt, jq.Jq11t)

		////////////////农历月历表
		_, listday, _ := zeji.ListLunarDay(jq, l)

		///////////小六壬择吉
		iqs := zeji.ZhiSu(s, g)
		starName := iqs.StarNames        //值宿名称
		zhisus := fmt.Sprintf(iqs.ZhiSu) //当日值宿信息
		//fmt.Printf("二十八宿:\"%s\"\n %s\n", starName, zhisus)
		//七煞判断
		qsB := iqs.IsQiSha(s.SolarDayT, g.DayZhiM)
		nx := zeji.AllNumber(g.YearZhi, l.LMonth, l.LDay, l.LHour)
		n1b := nx.YiPan()
		n2b := nx.ErPan()
		n3b := nx.SanPan()
		zeji := zeji.ShowResult(n1b, n2b, n3b, qsB)
		//fmt.Printf("择吉结果: %s\n", zeji)

		//择日 协纪辩方书
		hgz := convHourZhi(g.HourGanZhiM)
		zr = &ZR{
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
		yeartab := zr.YearTab()           //协纪辩方 年表
		djc, jcb := zr.JCM()              //协纪辩方 建除十二神煞(日)
		monthtab := zr.MonthTab(djc, jcb) //协纪辩方 月表
		daytab := zr.DayTab()             //协纪辩方 日表
		bianwei := zr.BianWei()           //协纪辩方 辩伪+其他
		rszl := zr.RSZL()                 //通书日时总览
		xjbfs := XJBF{
			NB: yeartab,
			YB: monthtab,
			RB: daytab,
			BW: bianwei,
			RS: rszl,
		}
		////月将
		jqt := ts.JQT(ly)
		solarT := time.Date(s.SYear, time.Month(s.SMonth), s.SDay, 0, 0, 0, 0, time.UTC)
		yjs := ts.NEWZRYLYueJiang(solarT, jqt)
		yjname := fmt.Sprintf("月将: %s", yjs.Name)
		star := fmt.Sprintf("十二宫: %s", yjs.Star)
		yj := YJ{
			YjName:   yjname,
			StarName: star,
		}

		resp := Resp{
			JiNian:   jinianinfo,
			Dmj:      dmjinfo,
			Jq:       jq24info,
			ListDay:  listday,
			StarName: starName,
			StarInfo: zhisus,
			Zeji:     zeji,
			XJBF:     xjbfs,
			YJ:       yj,
		}
		json.NewEncoder(w).Encode(resp)
	}
}

func about() About {
	ccal := "农历 择吉 可计算时间范围:1601～3498"
	data := "农历数据来源: https://github.com/ytliu0/ChineseCalendar/raw/master/TDBtimes.txt"
	xlr := "小六壬择吉 依据道家小六壬择法卷"
	xjbf := "择日 依据协纪辩方书 时辰吉凶参考<<讲武全书兵占>>通书部分"
	ck := "农历编算参考: https://ytliu0.github.io/ChineseCalendar/index_simp.html"
	me := "作者 梁子: xiaoyaoke7630@sina.com"
	return About{
		Ccal: ccal,
		Data: data,
		Xlr:  xlr,
		Xjbf: xjbf,
		Ck:   ck,
		Me:   me,
	}
}
func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/ccal", ccalyj)
	http.HandleFunc("/today", todayInfo)
	http.HandleFunc("/jz60", selectlist)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//年月日干支对应的数字
type GZS struct {
	Ygz, Mgz, Dgz, Hgz string
	JC                 string //日建除
	Dhh                string //黄黑神煞(日)
	Hhh                string //黄黑神煞(时辰)
	Tscy               string //太岁出游日
	Tygr               string //太岁天乙贵人
	TYDH               string //天乙贵人 日时论
	Luy                string //岁禄
	Lum                string //建禄
	Lud                string ///专禄(坐禄)
	Luh                string //归禄
	GuGua              string //孤辰寡宿
	WuLu               string //无禄日（十恶大败日）
	ChongRi, FuRi      string //重日 复日
	YiJu               string //移居吉日 含满成开日
	TanBing            string //忌探病日
	XianChi            string //咸池 桃花
	ShangShuo          string //上朔日
}

//干支年月日下拉选择
func selectlist(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("jz60.html")
		if err != nil {
			log.Fatal("t-err:", err)
		}
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		y, err := strconv.Atoi(r.Form["jzy"][0]) //这里的jzy 是select的name值
		if err != nil {
			log.Fatal("y-err:", err)
		}
		m, err := strconv.Atoi(r.Form["jzm"][0])
		if err != nil {
			log.Fatal("m-err:", err)
		}
		d, err := strconv.Atoi(r.Form["jzd"][0])
		if err != nil {
			log.Fatal("d-err:", err)
		}
		h, err := strconv.Atoi(r.Form["jzh"][0])
		if err != nil {
			log.Fatal("h-err", err)
		}
		///fmt.Printf("y:%d m:%d d:%d h:%d \n", y, m, d, h)
		///太岁出游日
		var ygz, mgz, dgz, hgz string
		for i := 0; i < 60; i++ {
			if i == y {
				ygz = jz60[i]
			}
			if i == m {
				mgz = jz60[i]
			}
			if i == d {
				dgz = jz60[i]
				//	break
			}
			if i == h {
				hgz = jz60[i]
			}
		}

		ymc := gztoYmc(mgz) //月名称 正 二 三...十一 十二
		fmt.Println(ygz, mgz, dgz, hgz, ymc)
		//////
		日建除, 日黄黑, 时辰黄黑 := jchh(ymc, mgz, dgz, hgz)
		太岁出游 := ts.XJBF太岁出游(dgz, jz60)
		阳贵人, 阴贵人 := ts.XJBF太岁天乙贵人(ygz)
		天乙贵人 := ts.XJBF天乙贵人(dgz, hgz)
		岁禄 := ts.XJBF岁禄(dgz, ygz)
		建禄 := ts.XJBF建禄(dgz, mgz)
		专禄 := ts.XJBF专禄(dgz)
		归禄 := ts.XJBF归禄(dgz, hgz)
		孤辰, 寡宿 := ts.XJBF孤辰寡宿(ygz, mgz, dgz, hgz)
		无禄 := ts.XJBF无禄日(dgz)
		重日 := ts.XJBF重日(dgz)
		复日 := ts.XJBF复日(ymc, dgz)
		移居吉日 := ts.RSZL移居吉日(ymc, dgz)
		忌探病日 := ts.RSZL忌探病日(dgz)
		咸池 := ts.XCTH咸池桃花(ygz, mgz, dgz, hgz)
		上朔日 := ts.XJBF上朔(ygz, dgz)

		gzs := GZS{
			Ygz:       ygz,
			Mgz:       mgz,
			Dgz:       dgz,
			Hgz:       hgz,
			JC:        日建除,
			Dhh:       日黄黑,
			Hhh:       时辰黄黑,
			Tscy:      太岁出游,
			Tygr:      阳贵人 + " " + 阴贵人,
			TYDH:      天乙贵人,
			Luy:       岁禄,
			Lum:       建禄,
			Lud:       专禄,
			Luh:       归禄,
			GuGua:     孤辰 + " " + 寡宿,
			WuLu:      无禄,
			ChongRi:   重日,
			FuRi:      复日,
			YiJu:      移居吉日,
			TanBing:   忌探病日,
			XianChi:   咸池,
			ShangShuo: 上朔日,
		}
		/* 		js, err := json.Marshal(gzs)
		   		if err != nil {
		   			log.Fatal("js-err", err)
		   		}
		   		io.WriteString(w, string(js)) */
		json.NewEncoder(w).Encode(gzs)
	}
}

//日建除
func jchh(ymc, mgz, dgz, hgz string) (djc, dhh, hhh string) {
	yjc := ts.JC本月建除(ymc)
	djc = ts.JC本日建除(dgz, yjc) //日建除
	zhimap := ts.MakeHHMap()
	dhh = ts.XJBF黄黑起法(mgz, dgz, zhimap) //本日黄黑神煞
	hhh = ts.XJBF黄黑起法(dgz, hgz, zhimap) //时辰黄黑神煞
	return
}

//今日信息
func todayInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("today.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		//今日阳历信息
		todayinfo := todayccal()
		//关于
		about := about()
		ti := Today{
			Ti:    todayinfo,
			About: about,
		}
		json.NewEncoder(w).Encode(ti)
	}

}

//默认当日结果
func todayccal() (infoToday string) {
	expectInfo, err := today.NewExpectInfo()
	if err != nil {
		log.Fatal("今日时间解析异常:", err)
	}
	///
	var lm string
	leapY := expectInfo.LeapY
	leapM := expectInfo.LeapM
	expectLeapD := expectInfo.ExpectleapD
	leapB := expectInfo.LeapB

	normalY := expectInfo.NormalY
	normalM := expectInfo.NormalM
	expectD := expectInfo.ExpectD
	normalB := expectInfo.NormalB

	h24 := T.Hour()
	h := utils.Conv24Hto12H(h24)
	sx := "猴"

	if leapM != 0 && leapB == true { //闰月月份
		err, s, l, g, _ := ccal.Input(leapY, leapM, expectLeapD, h, sx, leapB)
		if err != nil {
			log.Fatal()
		}
		if l.Leapmb == true {
			lm = "是"
		} else {
			lm = "否"
		}
		solarinfo := fmt.Sprintf("阳历纪年: %d年-%d月-%d日-周%s", s.SYear, s.SMonth, s.SDay, s.SWeek)
		lunarinfo := fmt.Sprintf("农历纪年: %d年%s月(%s)%s%s时(%d时)",
			l.LYear, lunar.Ymc[l.LMonth-1], l.LYdxs, lunar.Rmc[l.LDay-1], l.LaliasHour, l.LHour)
		gzinfo := fmt.Sprintf("干支纪年: %s%s年-%s月-%s%s日-%s时",
			g.YearGanM, g.YearZhiM, g.MonthGanZhiM, g.DayGanM, g.DayZhiM, g.HourGanZhiM)
		leapminfo := fmt.Sprintf("本年是否有闰月:%s-->闰%d月", lm, l.LeapMonth)
		infoToday = solarinfo + "<br />" + lunarinfo + "<br />" + gzinfo + "<br />" + leapminfo

	} else if normalM != 0 && normalB == false { //非闰月月份
		err, s, l, g, _ := ccal.Input(normalY, normalM, expectD, h, sx, normalB)
		if err != nil {
			log.Fatal(err)
		}
		if l.Leapmb == true {
			lm = "是"
		} else {
			lm = "否"
		}
		solarinfo := fmt.Sprintf("阳历纪年: %d年-%d月-%d日-周%s", s.SYear, s.SMonth, s.SDay, s.SWeek)
		lunarinfo := fmt.Sprintf("农历纪年: %d年%s月(%s)%s%s时(%d时)",
			l.LYear, lunar.Ymc[l.LMonth-1], l.LYdxs, lunar.Rmc[l.LDay-1], l.LaliasHour, l.LHour)
		gzinfo := fmt.Sprintf("干支纪年: %s%s年-%s月-%s%s日-%s时",
			g.YearGanM, g.YearZhiM, g.MonthGanZhiM, g.DayGanM, g.DayZhiM, g.HourGanZhiM)
		leapminfo := fmt.Sprintf("本年是否有闰月:%s-->闰%d月", lm, l.LeapMonth)
		infoToday = solarinfo + "<br />" + lunarinfo + "<br />" + gzinfo + "<br />" + leapminfo
	}
	return
}

//协纪辩方 年表
func (xjbf *ZR) YearTab() string {
	dgz := xjbf.dgz
	ygz := xjbf.aliasygz
	mgz := xjbf.aliasmgz
	yz := xjbf.aliasyz
	aliasmonth := xjbf.aliasmonth
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
func (xjbf *ZR) JCM() (string, bool) {
	m := xjbf.aliasmonth //农历月别名
	dgz := xjbf.dgz      //日干支
	i = &ts.ZRYL{
		AliasMonth: m,
		DGZ:        dgz,
	}
	return i.JC12M()
}

//协纪辩方 月表
func (xjbf *ZR) MonthTab(djc string, jcb bool) string {
	ygz := xjbf.aliasygz
	m := xjbf.aliasmonth
	ly := xjbf.lyear
	st := xjbf.stime
	mgz := xjbf.aliasmgz
	dgz := xjbf.dgz
	lday := xjbf.lday
	sy := xjbf.syear
	sm := xjbf.smonth
	sd := xjbf.sday
	i = &ts.ZRYL{
		YGZ:        ygz,
		AliasMonth: m,
		Lyear:      ly,
		SolarT:     st,
		MGZ:        mgz,
		DGZ:        dgz,
		Lday:       lday,
		Syear:      sy,
		Smonth:     sm,
		Sday:       sd,
	}

	return i.XJBF月表(djc, jcb)
}

//协纪辩方 日表
func (xjbf *ZR) DayTab() string {
	dgz := xjbf.dgz
	hgz := xjbf.hourgz
	i = &ts.ZRYL{
		DGZ: dgz,
		HGZ: hgz,
	}
	return i.XJBF日表(jz60)
}

//协纪辩方书 辩伪+其他
func (xjbf *ZR) BianWei() string {
	ygz := xjbf.aliasygz
	mgz := xjbf.aliasmgz
	dgz := xjbf.dgz
	hgz := xjbf.hourgz
	var guc, gus, th string
	guchen, guasu := ts.XJBF孤辰寡宿(ygz, mgz, dgz, hgz)
	if guchen != "" {
		guc = guchen
	} else if guasu != "" {
		gus = guasu
	}
	taohua := ts.XCTH咸池桃花(ygz, mgz, dgz, hgz)
	if taohua != "" {
		th = taohua
	}
	return "<br />" + "孤辰寡宿 咸池桃花" + "<br />" +
		guc + " " + gus + " " + th
}

//通书 日时总览方法
func (tos *ZR) RSZL() string {
	m := tos.aliasmonth
	dgz := tos.dgz
	rmc := tos.aliasday
	aliasHour := tos.aliasHour
	return ts.RSZLResult(m, dgz, rmc, aliasHour)
}

//时辰地支
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

//日名称 这里用的是廿
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

//干支对应的月名称 正 二 三...十一 十二
func gztoYmc(gz string) (ym string) {
	ymc := []string{"正", "二", "三", "四", "五", "六", "七", "八", "九", "十", "十一", "十二"}
	zhi := []string{"寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥", "子", "丑"}
	for i := 0; i < len(zhi); i++ {
		if strings.ContainsAny(gz, zhi[i]) {
			switch i { //0正月
			case 0:
				ym = ymc[0]
			case 1:
				ym = ymc[1]
			case 2:
				ym = ymc[2]
			case 3:
				ym = ymc[3]
			case 4:
				ym = ymc[4]
			case 5:
				ym = ymc[5]
			case 6:
				ym = ymc[6]
			case 7:
				ym = ymc[7]
			case 8:
				ym = ymc[8]
			case 9:
				ym = ymc[9]
			case 10:
				ym = ymc[10]
			case 11:
				ym = ymc[11]
			}
			break
		}
	}
	return
}
