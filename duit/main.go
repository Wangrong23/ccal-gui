package main

import "os"

func init() {
	os.Setenv("font", "./font/unifont.font")
}

func main() {
	run()
}
