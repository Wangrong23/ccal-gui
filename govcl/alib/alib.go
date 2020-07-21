package alib

import (
	"runtime"

	"github.com/ying32/govcl/pkgs/libname"
)

func init() {
	os := runtime.GOOS
	if os == "linux" {
		libname.LibName = "/usr/local/lib/govcl/liblcl.so"
	}
	if os == "windows" {
		libname.LibName = "C:\\govcl"
	}
}
