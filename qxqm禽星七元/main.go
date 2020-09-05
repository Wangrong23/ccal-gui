package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mjl-/duit"
	"github.com/nongli/ccal"
	"github.com/nongli/lunar"
	"github.com/qxqm/禽"
)

var (
	T      = time.Now().Local()
	status = &duit.Label{Text: "禽星七元\n计算方法依据古本禽星\n使用方法\n输入农历的年月日时辰闰月:\n年4位数 月 日 时辰 都是两位数字 f输入月份非闰月 t输入月份为闰月\n比如2020060611f\n"}
	get    = &duit.Field{}
)

func init() {
	os.Setenv("font", "/opt/fonts/unifont/unifont.font")
}

func main() {
	runMain()
}

func runMain() {
	dui, err := duit.NewDUI("禽星七元", nil)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	//信息
	basicInfo := &duit.Button{
		Text:     "确定",
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
			//日期信息
			err, s, l, g, jq := ccal.Input(y, m, d, h, sx, b)
			if err != nil {
				log.Fatal("ccal:", err)
			}
			yg := g.YearGanM
			yz := g.YearZhiM
			ygz := fmt.Sprintf("%s%s", yg, yz) //年干支
			mgz := g.MonthGanZhiM              //月干支
			dg := g.DayGanM                    //日干
			dz := g.DayZhiM                    //日支
			hgz := g.HourGanZhiM               //时干支
			dgz := fmt.Sprintf("%s%s", dg, dz) //日干支
			纪年 := fmt.Sprintf("纪年信息:%s年-%s月-%s日-%s时\n", ygz, mgz, dgz, hgz)

			//年分元信息
			x元 := 禽.Y年元信息(y)
			info := 禽.Q年禽(x元)
			yuanstar := info.Find()              //本元数组
			star, sgz := info.Y年禽(ygz, yuanstar) //年份对应的禽星干支

			元, _ := 禽.Conv(x元 - 1) //转文字
			info元 := fmt.Sprintf("%d年: %s\n", y, 元)
			info年禽 := fmt.Sprintf("年禽: %s 干支:%s\n", star, sgz)

			//月禽
			xyuan := 禽.X元信息(ygz, yuanstar)
			mself := xyuan.X月禽(m)
			mgzx := 禽.X禽干支(mself.M禽名, yuanstar)
			info月禽 := fmt.Sprintf("%d月禽星: %s　属性: %s 禽干支: %s\n", m, mself.M禽名, mself.M禽属性, mgzx)

			//日禽
			xd := info.X元信息()
			xdinfo := xd.D日禽(dgz)
			name := xdinfo.D名称
			xdself := xdinfo.D属性
			info日禽 := fmt.Sprintf("%s日禽星: %s 属性: %s\n", dgz, name, xdself)

			//时禽
			xhinfo := xdinfo.H时禽(h)
			info时禽 := fmt.Sprintf("时禽名称:%s 时禽属性:%s\n", xhinfo.H禽名, xhinfo.H禽属性)

			//十二宮克應
			qh28 := 禽.QH二十八宿十二宮克應()
			qh28info := qh28.QH二十八宿十二宮克應方法(star, h)
			info十二宫 := fmt.Sprintf("十二宮克應:%s\n", qh28info)

			//日禽十二宮克應
			qd := 禽.QD日禽十二時辰七元克應()
			infodh := qd.QDH日禽時辰克應(h, x元)
			info日禽十二宫 := fmt.Sprintf("日禽十二時辰克應:%s", infodh)

			//特殊吉凶日
			var (
				info出师吉凶 string
				info干克支  string
				info伐日   string
				info猖鬼败亡 string
				info八专   string
				info五不归  string
				info八绝日  string
				info用兵吉日 string
				info十恶大败 string
				info攻取吉日 string
				info天败凶日 string
			)

			jx := 禽.NewJX吉凶日()
			chushib, isgood := jx.JX出師吉日er(dgz)
			if chushib == true {
				info出师吉凶 = fmt.Sprintf("%s\n", isgood)
			}
			gkzb, 干克支 := jx.Jx干克支er(dgz)
			if gkzb == true {
				info干克支 = fmt.Sprintf("%s\n", 干克支)
			}
			farib, 支克干 := jx.Jx伐日er(dgz)
			if farib == true {
				info伐日 = fmt.Sprintf("%s\n", 支克干)
			}
			changguib, 猖鬼 := jx.Jx猖鬼败亡er(dgz)
			if changguib == true {
				info猖鬼败亡 = fmt.Sprintf("%s\n", 猖鬼)
			}
			bazhuanb, 八专 := jx.Jx八专er(dgz)
			if bazhuanb == true {
				info八专 = fmt.Sprintf("%s\n", 八专)
			}
			wubuguib, 五不归 := jx.Jx五不归er(dgz)
			if wubuguib == true {
				info五不归 = fmt.Sprintf("%s\n", 五不归)
			}
			bajuerib, 八绝日 := jx.Jx八绝日er(dgz)
			if bajuerib == true {
				info八绝日 = fmt.Sprintf("%s\n", 八绝日)
			}
			yongb, 用兵吉日 := jx.Jx用兵吉日er(lunar.Rmc[l.LDay-1])
			if yongb == true {
				info用兵吉日 = fmt.Sprintf("%s\n", 用兵吉日)
			}
			//十恶大败
			seb, 十恶大败 := jx.Jx十恶大败er(yg, lunar.Ymc[l.LMonth-1], dgz)
			if seb == true {
				info十恶大败 = fmt.Sprintf("%s\n", 十恶大败)
			} else if seb == false {
				info十恶大败 = fmt.Sprintf("%s\n", 十恶大败)
			}
			gqb, 攻取吉日 := jx.Jx攻取吉日er(dgz, farib, bazhuanb, changguib, wubuguib, seb, bajuerib)
			if gqb == true {
				info攻取吉日 = fmt.Sprintf("%s\n", 攻取吉日)
			}
			//天败凶日
			tbb, 天败凶日 := jx.Jx天败凶日er(s.SolarDayT, jq.Jqt, dz)
			if tbb == true {
				info天败凶日 = fmt.Sprintf("%s\n", 天败凶日)
			}
			//UI显示
			status.Text = 纪年 + info元 + info年禽 + info月禽 + info日禽 + info时禽 +
				info十二宫 + info日禽十二宫 +
				info出师吉凶 + info干克支 + info伐日 + info猖鬼败亡 + info八专 + info五不归 + info八绝日 +
				info用兵吉日 + info十恶大败 + info攻取吉日 + info天败凶日
			dui.MarkLayout(status)
			return
		},
	}

	//ui
	dui.Top.UI = &duit.Box{
		//Width:   160,
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
					&duit.Label{Text: "输入"},
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

type Input struct {
	Y int  //年數字
	M int  //月數字
	D int  //日數字
	H int  //時辰數字　1子時　2丑時...１２亥時
	B bool //閏月判斷　f非閏月　t當月爲閏月
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
