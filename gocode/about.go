package main

import (
	"fmt"

	"github.com/ying32/govcl/vcl"
)

//关于软件　关于作者
func (f *TForm1) OnButton5Click(sender vcl.IObject) {
	me := fmt.Sprintf("邮箱:xiaoyaoke7630@sina.com\n")
	ccal := fmt.Sprintf("中国农历择吉\n核心代码为Go,UI部分由Lazarus生成(govcl)\n")
	zeji := fmt.Sprintf("择吉部分算法依据道家小六壬择法卷\n")
	jinji := fmt.Sprintf("禁忌部分为来源民间\n")

	vcl.ShowMessageFmt(ccal + me + zeji + jinji)
}
