package main

import (
	"fmt"
	"strconv"

	"github.com/nongli/ccal"

	"github.com/andlabs/ui"
)

func main() {
	ui.Main(w)
}

//主窗口
func w() {

	//生成主窗口
	win := ui.NewWindow("农历择吉", 800, 600, true)
	//添加鼠标点击之后关闭窗口
	win.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})

	//添加＂退出＂菜单时执行的函数功能
	ui.OnShouldQuit(func() bool {
		win.Destroy()
		return true
	})

	//设置控件
	tab := ui.NewTab()
	Info(tab, win)    //历法信息
	ListDay(tab, win) //农历月历表
	JiQi24(tab, win)  //24节气
	Dimu(tab, win)    //地母经
	About(tab, win)   //软件信息

	hbox := ui.NewHorizontalBox() //创建水平框架
	hbox.SetPadded(true)

	//显示
	win.Show()

}

//历法信息
//生成输入信息控件
func Info(tab *ui.Tab, win *ui.Window) {
	win.SetChild(tab)
	win.SetMargined(true)
	tab.Append("历法信息", TabBasicInfo())
	tab.SetMargined(0, true)
}

//获取输入内容
func TabBasicInfo() ui.Control {

	hbox := ui.NewHorizontalBox() //创建水平框架
	hbox.SetPadded(true)

	vbox := ui.NewVerticalBox() //创建一个新的垂直Box
	vbox.SetPadded(true)
	hbox.Append(vbox, false) //将给定的控件添加到Box的末尾。

	//创建一个新的DateTimePicker，它同时显示日期和时间
	//将给定的控件添加到Box的末尾
	vbox.Append(ui.NewDatePicker(), false)

	//生成输入框文本 农历年月日时
	year := ui.NewEntry()
	month := ui.NewEntry()
	day := ui.NewEntry()
	hour := ui.NewEntry()
	sx := ui.NewEntry() //生肖
	b := ui.NewEntry()  //闰月
	// 生成标签对应上面文本输入框
	laby := ui.NewLabel(``)  //year
	labm := ui.NewLabel(``)  //month
	labd := ui.NewLabel(``)  //day
	labh := ui.NewLabel(``)  //hour
	labsx := ui.NewLabel(``) //生肖
	labb := ui.NewLabel(``)  //输入月份是否为闰月

	// 生成按钮
	button := ui.NewButton(`查看`)

	// 设置按钮点击事件
	button.OnClicked(func(*ui.Button) {
		//lab.SetText(year.Text() + month.Text() + day.Text() + hour.Text())//一次性获取全部输入内容

		//获取标签内输入的内容
		laby.SetText(year.Text())
		labm.SetText(month.Text())
		labd.SetText(day.Text())
		labh.SetText(hour.Text())
		labsx.SetText(sx.Text())
		labb.SetText(b.Text())

		//获取输入的内容
		lsY := laby.Text()
		lsM := labm.Text()
		lsD := labd.Text()
		lsH := labh.Text()
		lsSx := labsx.Text()
		lsB := labb.Text()
		fmt.Printf("lsY:%s lsM:%s lsD:%s lsH:%s lsSx:%s lsB:%s\n", lsY, lsM, lsD, lsH, lsSx, lsB) //打印到终端
		//格式解析
		y, _ := strconv.ParseInt(lsY, 10, 32)
		m, _ := strconv.ParseInt(lsM, 10, 32)
		d, _ := strconv.ParseInt(lsD, 10, 32)
		h, _ := strconv.ParseInt(lsH, 10, 32)
		//转换闰月值为函数类型
		var b bool
		if lsB == "是" {
			b = true
		} else if lsB == "否" {
			b = false
		}

		s, l, g, jq := ccal.Input(int(y), int(m), int(d), int(h), lsSx, b)
		fmt.Printf("%v %v %v %v\n", s, l, g, jq)
	})

	//添加到竖列Ｂox
	vbox.Append(year, false)
	vbox.Append(month, false)
	vbox.Append(day, false)
	vbox.Append(hour, false)
	vbox.Append(sx, false)
	vbox.Append(b, false)

	vbox.Append(laby, false)
	vbox.Append(button, false)
	//ui.NewVerticalSeparator 创建一个新的垂直分隔符。
	//hbox.Append 将给定的控件添加到Box的末尾。
	hbox.Append(ui.NewVerticalSeparator(), false)

	//待用
	hbox.Append(ui.NewVerticalSeparator(), false)
	hbox.Append(ui.NewVerticalSeparator(), false)

	fmt.Printf("历法基础\n")
	return hbox
}

//显示月历表
func ListDay(tab *ui.Tab, win *ui.Window) {
	win.SetChild(tab)
	win.SetMargined(true)
	tab.Append("农历月历表", TabListDay())
	tab.SetMargined(1, true)
}

//月历表
func TabListDay() ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, false)

	fmt.Printf("农历月历表\n")
	return hbox
}

//节气信息
func JiQi24(tab *ui.Tab, win *ui.Window) {
	win.SetChild(tab)
	win.SetMargined(true)
	tab.Append("24节气", TabJiQi24())
	tab.SetMargined(2, true) //数字为显示顺序
}

//显示24节气标签
func TabJiQi24() ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, false)

	fmt.Printf("24jieqi\n")
	return hbox
}

//地母经内容
func Dimu(tab *ui.Tab, win *ui.Window) {
	win.SetChild(tab)
	win.SetMargined(true)
	tab.Append("地母经", TabDimu())
	tab.SetMargined(3, true)
}

//地母经
func TabDimu() ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, false)

	fmt.Printf("地母经内容\n")
	return hbox
}

//软件版本信息
func About(tab *ui.Tab, win *ui.Window) {
	win.SetChild(tab)
	win.SetMargined(true)
	tab.Append("关于", TabDimu())
	tab.SetMargined(3, true)
}

//软件信息
func TabAbout() ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, false)
	fmt.Printf("这里是软件版本的内容部分\n")
	return hbox
}
