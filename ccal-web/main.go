package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/Aquarian-Age/nongli/ccal"
	"github.com/Aquarian-Age/nongli/dimu"
	"github.com/Aquarian-Age/nongli/lunar"
	"github.com/Aquarian-Age/nongli/solar"
	"github.com/Aquarian-Age/nongli/zeji"
	ganzhi "github.com/Aquarian-Age/ts/gz"
)

func get(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		f(w, r)
	}
}

//应答数据
type Resp struct {
	Solar    string   `json:"solarInfo"`    //阳历纪年信息
	Lunar    string   `json:"lunarInfo"`    //农历纪年信息
	GanZhi   string   `json:"gzInfo"`       //干支信息
	NaYin    string   `json:"nyInfo"`       //纳因信息
	Dmj      string   `json:"dmInfo"`       //地母经
	Jq       []string `json:"jqInfo"`       //24节气
	ListDay  []string `json:"listdayInfo"`  //农历月历表(农历初一开始)
	StarName string   `json:"starnameInfo"` //当日值宿名称(28宿)
	StarInfo string   `json:"starInfo"`     //值宿信息
	Zeji     string   `json:"zejiInfo"`     //当日择吉信息
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
		/////////////////地母经
		dmg := g.YearGan
		dmz := g.YearZhi
		infodmj := dimu.DimuInfo(dmg, dmz)
		dmjinfo := fmt.Sprintf("%s", infodmj)
		///////////////24节气
		jq24info := solar.ShowJieqi24(jq.Jqt, jq.Jq11t)
		////////////////农历月历表
		_, listday, _ := zeji.ListLunarDay(jq, l)
		/* 		listday := func(s []string) []string {
			var day string
			var days []string
			for i := 0; i < len(s); i++ {
				day = s[i] + "<br />"
				days = append(days, day)
			}
			return days
		}(list) */
		///////////小六壬择吉
		//择吉部分
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

		resp := Resp{
			Solar:    solarinfo,
			Lunar:    lunarinfo,
			GanZhi:   gzinfo,
			NaYin:    nyinfo,
			Dmj:      dmjinfo,
			Jq:       jq24info,
			ListDay:  listday,
			StarName: starName,
			StarInfo: zhisus,
			Zeji:     zeji,
		}

		/////////
		fmt.Println(resp)
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
