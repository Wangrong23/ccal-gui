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

//应答数据
type Resp struct {
	JiNian   string   `json:"jinianInfo"`   //纪年信息
	Dmj      string   `json:"dmInfo"`       //地母经
	Jq       []string `json:"jqInfo"`       //24节气
	ListDay  []string `json:"listdayInfo"`  //农历月历表(农历初一开始)
	StarName string   `json:"starnameInfo"` //当日值宿名称(28宿)
	StarInfo string   `json:"starInfo"`     //值宿信息
	Zeji     string   `json:"zejiInfo"`     //当日择吉信息
	XJBF
}

//表单提交
func forminput(w http.ResponseWriter, r *http.Request) {
	//请求
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, nil)
	} else {
		//解析表单
		r.ParseForm()

		ly, err := strconv.Atoi(r.Form["ly"][0])
		if err != nil {
			log.Fatalln("ly:", err)
		}
		lm, err := strconv.Atoi(r.Form["lm"][0])
		if err != nil {
			log.Fatalln(err)
		}
		ld, err := strconv.Atoi(r.Form["ld"][0])
		if err != nil {
			log.Fatalln(err)
		}
		lh, err := strconv.Atoi(r.Form["lh"][0])
		if err != nil {
			log.Fatalln(err)
		}
		sx := r.Form["la"][0]
		sb := r.Form["le"][0]
		var leapb bool
		if strings.EqualFold(sb, "t") {
			leapb = true
		} else if strings.EqualFold(sb, "f") {
			leapb = false
		}

		//应答 (浏览器开发模式 Console看结果)
		/////////////////////////////ccal农历基本纪年信息
		b := leapb
		_, s, l, g, jq := ccal.Input(ly, lm, ld, lh, sx, b)
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
		fmt.Printf("二十八宿:\"%s\"\n %s\n", starName, zhisus)
		//七煞判断
		qsB := iqs.IsQiSha(s.SolarDayT, g.DayZhiM)
		nx := zeji.AllNumber(g.YearZhi, l.LMonth, l.LDay, l.LHour)
		n1b := nx.YiPan()
		n2b := nx.ErPan()
		n3b := nx.SanPan()
		zeji := zeji.ShowResult(n1b, n2b, n3b, qsB)
		fmt.Printf("择吉结果: %s\n", zeji)

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

		resp := Resp{
			JiNian:   jinianinfo,
			Dmj:      dmjinfo,
			Jq:       jq24info,
			ListDay:  listday,
			StarName: starName,
			StarInfo: zhisus,
			Zeji:     zeji,
			XJBF:     xjbfs,
		}

		//fmt.Println(resp)
		json.NewEncoder(w).Encode(resp)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", "农历 协纪辩方 通书 择日")
	if r.Method == "GET" {
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
	}
}
func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/forminput", forminput) //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
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
	return "\n----孤辰寡宿 咸池桃花----\n" +
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
