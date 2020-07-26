package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/nongli/ccal"
	"github.com/qxqm/禽"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

func main() {
	rect := sciter.Rect{Left: 0, Top: 0, Right: 640, Bottom: 480}
	w, err := window.New(
		sciter.SW_TITLEBAR|
			sciter.SW_RESIZEABLE|
			sciter.SW_CONTROLS|
			sciter.SW_MAIN|
			sciter.SW_ENABLE_DEBUG,
		//设置窗口大小
		&rect)
	if err != nil {
		log.Fatal("window:", err)
	}
	//从html文件载入内容
	/* 	//fp, err := filepath.Abs("ccal-sciter.html")
	   	if err != nil {
	   		log.Fatal("file:", err)
	   	}

		   w.LoadFile(fp) */
	//载入代码html内容
	w.LoadHtml(html, "")
	//前端交互
	w.DefineFunction("yearinfo", yearinfo)

	w.Show()
	w.Run()
}

//接受前端按钮
func year(root *sciter.Element) {
	info, err := root.SelectById("btny")
	if err != nil {
		log.Fatal(err)
	}
	info.DefineMethod("yearinfo", yearinfo)
}

//获取前端输入的内容
func yearinfo(args ...*sciter.Value) *sciter.Value {
	if args[0].IsString() == true {
		log.Printf("输入的信息：%v\n", args[0])
	}
	//////
	var s string
	for _, arg := range args {
		s = arg.String()
	}

	infos := conv(s)
	y := infos.Y
	m := infos.M
	d := infos.D
	h := infos.H
	b := infos.B
	sx := "猴"

	//日期干支信息
	err, _, _, g, _ := ccal.Input(y, m, d, h, sx, b)
	if err != nil {
		log.Fatal("ccal:", err)
	}
	yg := g.YearGanM
	yz := g.YearZhiM
	ygz := fmt.Sprintf("%s%s", yg, yz) //年干支
	mgz := g.MonthGanZhiM              //月干支
	dg := g.DayGanM
	dz := g.DayZhiM
	hgz := g.HourGanZhiM               //时干支
	dgz := fmt.Sprintf("%s%s", dg, dz) //日干支
	纪年 := fmt.Sprintf("纪年信息:%s年-%s月-%s日-%s时\n", ygz, mgz, dgz, hgz)

	//年禽
	x元 := 禽.Y年元信息(y)
	info := 禽.Q年禽(x元)
	yuanstar := info.Find()              //本元数组
	star, sgz := info.Y年禽(ygz, yuanstar) //年份对应的禽星干支
	元, _ := 禽.Conv(x元 - 1)               //转文字
	年元 := fmt.Sprintf("%d年: %s\n", y, 元)
	年禽 := fmt.Sprintf("年禽: %s 干支:%s\n", star, sgz)

	//月禽
	xyuan := 禽.X元信息(ygz, yuanstar)
	mself := xyuan.X月禽(m)
	mxgz := 禽.X禽干支(mself.M禽名, yuanstar)
	月禽 := fmt.Sprintf("\n月禽: %s　属性: %s 禽干支: %s\n", mself.M禽名, mself.M禽属性, mxgz)
	//("\n%d月禽星: %s　属性: %s 禽干支: %s\n", m, mself.M禽名, mself.M禽属性, mxgz)

	//日禽
	xd := info.X元信息()
	xdinfo := xd.D日禽(dgz)
	name := xdinfo.D名称
	xdself := xdinfo.D属性
	日禽 := fmt.Sprintf("\n日禽: %s 属性: %s\n", name, xdself)
	//("%s日禽星: %s 属性: %s\n", dgz, name, xdself)

	//时禽
	xhinfo := xdinfo.H时禽(h)
	时禽 := fmt.Sprintf("\n时禽:%s 属性:%s\n", xhinfo.H禽名, xhinfo.H禽属性)

	//十二宮克應
	qh28 := 禽.QH二十八宿十二宮克應()
	qh28info := qh28.QH二十八宿十二宮克應方法(star, h)
	十二宫克应 := fmt.Sprintf("\n十二宮克應:%s\n", qh28info)

	//日禽十二宮克應
	qd := 禽.QD日禽十二時辰七元克應()
	infodh := qd.QDH日禽時辰克應(h, x元)
	日禽十二宫克应 := fmt.Sprintf("\n日禽十二時辰克應:%s\n", infodh)

	return sciter.NewValue(纪年 + 年元 + 年禽 + 月禽 + 日禽 + 时禽 + 十二宫克应 + 日禽十二宫克应) //返回到ui页面
}

//---------
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

	//fmt.Printf("input: y:%v m:%v d:%v h:%v b:%t\n", y, m, d, h, b)
	return &Input{
		Y: y,
		D: d,
		M: m,
		H: h,
		B: b,
	}
}

var html = `

<html>
    <head>
        <meta charset="UTF-8"> 
        <TItle>七元禽星</TItle>
    </head>
    <body>
        <!--输入年-->
        <h1>七元禽星古本</h1>
        <input |text #year></input>

        <!-- 按钮 -->
        <button #btny>确定</button>

        <script type="text/tiscript">
            //另一种方式w.DefineFunction取UI输入的信息
            event click $(#btny){
            // view.yearinfo($(#year).value)//传递输入信息到Go
            var inputGet = view.yearinfo($(#year).value)//Go处理之后的结果

            //显示输入的内容
            $(#输入信息).value= $(#year).value
                
            //Go处理之后的返回值显示到#getting的地方
            $(#getting).value=inputGet
            }
        </script>

        <!-- 指定信息显示位置 -->
        <h4 #输入信息></h4>
        <p #getting></p>

        <!-- 阳历日历 -->
        <h4>阳历日历</h4>
        <input type="date"></input>
        <!-- <button #scal></button>
        <script type="text/tisctipt">
        event click $(#scal){
            today()
        } 
        </script>-->

    </body>
</html>

`
