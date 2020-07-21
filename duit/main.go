package main

import "os"

func init() {
	os.Setenv("font", "./font/unifont.font")
	//	os.Setenv("font", "/opt/fonts/unifont/unifont.font")
}
func main() {
	run()
}
