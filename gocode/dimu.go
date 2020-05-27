package main

import (
	"fmt"
	"strconv"

	"github.com/nongli/ccal"
	"github.com/nongli/dimu"
	"github.com/nongli/ganzhi"
	"github.com/ying32/govcl/vcl"
)

func (f *TForm1) OnButton6Click(sender vcl.IObject) {
	year := f.Edit1.Text()
	if year == "" {
		fmt.Printf("输入年份数字\n")
		f.Edit1.SetFocus()
		return
	}

	if year != "" {
		_y, _ := strconv.ParseInt(year, 10, 32)
		y := int(_y)
		_, _, g, _ := ccal.Input(y, 3, 1, 1, "猴", false) //这里月份要在立春之后
		dmg := g.YearGan
		dmz := g.YearZhi
		infodmj := dimu.DimuInfo(dmg, dmz)
		//fmt.Printf("%d %v-%v\n", y, dmg, dmz)
		info := fmt.Sprintf("%s%s年地母经:\n%s\n", ganzhi.Gan[dmg], ganzhi.Zhi[dmz], infodmj)
		vcl.ShowMessageFmt(info)
	}
}
