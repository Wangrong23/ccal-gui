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
	"github.com/Aquarian-Age/nongli/lunar"
	ganzhi "github.com/Aquarian-Age/ts/gz"
)

func get(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		f(w, r)
	}
}

type Input struct {
	Ly   int    `json:"lunarY"`
	Lm   int    `json:"lunarM"`
	Ld   int    `json:"lunarD"`
	Lh   int    `json:"lunarH"` //时辰
	Sx   string `json:"lunarA"` //生肖
	Leap bool   `json:"lunarB"` //输入月份是否闰月
}

//应答数据
type Resp struct {
	Solar  string `json:"solarInfo"` //阳历纪年信息
	Lunar  string `json:"lunarInfo"` //农历纪年信息
	GanZhi string `json:"gzInfo"`    //干支信息
	NaYin  string `json:"nyInfo"`    //纳因信息
}

//表单提交
func forminput(w http.ResponseWriter, r *http.Request) {
	//请求
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		//t, _ := template.ParseFiles("input.gptl")
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, nil)
	} else {
		//解析表单
		r.ParseForm()
		ly, err := strconv.Atoi(r.Form["ly"][0])
		if err != nil {
			log.Fatalln(err)
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

		input := Input{
			Ly:   ly,
			Lm:   lm,
			Ld:   ld,
			Lh:   lh,
			Sx:   sx,
			Leap: leapb,
		}
		fmt.Println(input)
		//应答 (浏览器开发模式 Console看结果)
		//json.NewEncoder(w).Encode(input) //{lunarY: 1166, lunarM: 2266, lunarD: 3388}
		/////////////////////////////ccal农历基本纪年信息
		//sx := "猴"
		b := leapb
		_, s, l, g, _ := ccal.Input(ly, lm, ld, lh, sx, b)
		var aliasM string
		if l.Leapmb == true {
			aliasM = "是"
		} else {
			aliasM = "否"
		}
		solar := fmt.Sprintf("阳历纪年: %d年-%d月-%d日-周%s\n", s.SYear, s.SMonth, s.SDay, s.SWeek)
		lunar := fmt.Sprintf("农历纪年: %d年%s月(%s)%s %s时(%d时)<br />本年是否有闰月:%s闰%d月\n", l.LYear, lunar.Ymc[l.LMonth-1], l.LYdxs, lunar.Rmc[l.LDay-1], l.LaliasHour, l.LHour, aliasM, l.LeapMonth)
		gz := fmt.Sprintf("干支纪年: %s%s年-%s月-%s%s日-%s时\n",
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
		ny := fmt.Sprintf("干支纳音: %s %s %s %s\n", ygzny[aliasygz], mgzny[aliasmgz], dgzny[dgz], hgzny[hourgz])
		resp := Resp{
			Solar:  solar,
			Lunar:  lunar,
			GanZhi: gz,
			NaYin:  ny,
		}
		fmt.Println(resp)
		json.NewEncoder(w).Encode(resp)
	}
}

func main() {
	http.HandleFunc("/forminput", forminput) //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
