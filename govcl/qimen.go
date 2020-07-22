package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nongli/ccal"
	"github.com/nongli/ganzhi"
	"github.com/nongli/qimen"
	"github.com/nongli/solar"
)

//時家奇門(拆補法)
func qm(st time.Time, g *ccal.LunarGanZhiInfo, jq *ccal.JieQiInfo) (Text string) {

	fg, offg := qimen.FuTouGan(g.DayGan)
	//fmt.Printf("當日天幹數字和符頭天干的差值:%d\n", offg)

	fz := qimen.FuTouZhi(g.DayZhi, offg)
	//fmt.Printf("符頭天干數字:%d 符頭地支數字:%d\n", fg, fz)
	f符头 := fmt.Sprintf("符頭:%s%s\n", ganzhi.Gan[fg], ganzhi.Zhi[fz])

	yuan := qimen.FuTouYuan(fg, fz)
	jqt := qimen.AllJqt(jq.Jqt, jq.Jq11t)
	jmc := qimen.FestivalName(st, jqt)
	//fmt.Println("當日爲", jmc, yuan)

	bginfo, err := qimen.BaGongInfo(jmc)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("節氣對應的八宫信息:%v\n", bginfo)

	jie := qimen.ConvJie(jmc, solar.JMC)
	saninfo := bginfo.DingJiu(jie, yuan)
	//fmt.Printf("三元信息:%v\n", saninfo)
	//精確的節氣時間
	jieInfo, _ := qimen.J24H(st.Year(), jie)

	ju := saninfo.DingJu(yuan)
	info := fmt.Sprintf("拆補定局: %s %s 第%d天 %s遁%d局\n", jie, yuan, offg+1, bginfo.YinYang, ju)
	//旬首
	旬首 := qimen.X旬首(g.HourGanZhiM)
	x旬首 := fmt.Sprintf("旬首:%s\n\n", 旬首)

	Text = jieInfo + info + f符头 + x旬首
	return
}
