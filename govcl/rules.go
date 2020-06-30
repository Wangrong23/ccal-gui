package main

import (
	"errors"
	"strconv"
	"strings"
)

var (
	aliasshu, aliasniu, aliashu, aliastu, aliaslong, aliasshe string
	aliasma, aliasyang, aliashou, aliasji, aliasgou, aliaszhu string
	aliaslongf, aliasmaf, aliasjif, aliaszhuf                 string
)

//字符串类型转int　返回值为真表示输入数字正常
func String2Int(year, month, day, hour string) (y, m, d, h int, inputb bool, err error) {
	y, _ = strconv.Atoi(year)
	m, _ = strconv.Atoi(month)
	d, _ = strconv.Atoi(day)
	h, _ = strconv.Atoi(hour)
	inputb, err = dateBool(y, m, d, h)

	return
}

//判断输入的数字
func dateBool(year, month, day, hour int) (dateB bool, err error) {

	if (year > 1600 && year < 3499) &&
		(month >= 1 && month <= 12) &&
		(day >= 1 && day <= 30) &&
		(hour >= 1 && hour <= 12) {
		dateB = true
	} else {
		err = errors.New("年份时间范围1600到3499")
		dateB = false
	}
	return
}

//判断输入是不是闰月
func leapBool(leapm string) (lt bool, err error) {

	sl := strings.ToLower(leapm)
	slby := strings.EqualFold(sl, "y")
	slby1 := strings.EqualFold(sl, "yes")
	slbn := strings.EqualFold(sl, "n")
	slbn1 := strings.EqualFold(sl, "no")

	if slby == true || slby1 == true {
		lt = true
	} else if slbn == false || slbn1 == false {
		lt = false
	}

	if leapm != "yes" && leapm != "y" &&
		leapm != "no" && leapm != "n" {
		err = errors.New("闰月判断值输入错误软件自动关闭...")
	}
	return
}

//生肖判断　可以输入生肖的拼音或者汉字支持繁体
func shengxiao(s string) (t bool) {
	lows := strings.ToLower(s) //转为小写

	aliasshu, aliasniu, aliashu, aliastu, aliaslong, aliasshe = "鼠", "牛", "虎", "兔", "龙", "蛇"
	aliasma, aliasyang, aliashou, aliasji, aliasgou, aliaszhu = "马", "羊", "猴", "鸡", "狗", "猪"
	aliaslongf, aliasmaf, aliasjif, aliaszhuf = "龍", "馬", "雞", "豬"

	//简体部分
	shub := strings.EqualFold(s, aliasshu)
	niub := strings.EqualFold(s, aliasniu)
	hub := strings.EqualFold(s, aliashu)
	tub := strings.EqualFold(s, aliastu)
	longb := strings.EqualFold(s, aliaslong)
	sheb := strings.EqualFold(s, aliasshe)
	mab := strings.EqualFold(s, aliasma)
	yangb := strings.EqualFold(s, aliasyang)
	houb := strings.EqualFold(s, aliashou)
	jib := strings.EqualFold(s, aliasji)
	goub := strings.EqualFold(s, aliasgou)
	zhub := strings.EqualFold(s, aliaszhu)

	//繁體部分
	longfb := strings.EqualFold(s, aliaslongf)
	mafb := strings.EqualFold(s, aliasmaf)
	jifb := strings.EqualFold(s, aliasjif)
	zhufb := strings.EqualFold(s, aliaszhuf)

	//拼音部分
	shuB := strings.EqualFold(lows, "shu")
	niuB := strings.EqualFold(lows, "niu")
	huB := strings.EqualFold(lows, "hu")
	tuB := strings.EqualFold(lows, "tu")
	longB := strings.EqualFold(lows, "long")
	sheB := strings.EqualFold(lows, "she")
	maB := strings.EqualFold(lows, "ma")
	yangB := strings.EqualFold(lows, "yang")
	houB := strings.EqualFold(lows, "hou")
	jiB := strings.EqualFold(lows, "ji")
	gouB := strings.EqualFold(lows, "gou")
	zhuB := strings.EqualFold(lows, "zhu")

	if (shub == false && shuB == false) &&
		(niub == false && niuB == false) &&
		(hub == false && huB == false) &&
		(tub == false && tuB == false) &&
		(longb == false && longB == false && longfb == false) &&
		(sheb == false && sheB == false) &&
		(mab == false && maB == false && mafb == false) &&
		(yangb == false && yangB == false) &&
		(houb == false && houB == false) &&
		(jib == false && jiB == false && jifb == false) &&
		(goub == false && gouB == false) &&
		(zhub == false && zhuB == false && zhufb == false) {

		t = false
	} else {
		t = true
	}
	return
}
