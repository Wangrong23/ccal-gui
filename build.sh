#!/bin/bash

go build -o ccal-gui-linux64-v0.3.6 -ldflags="-s -w" Form1.go Form1Impl.go jieqi24.go about.go dimu.go zeji.go listDay.go rules.go jinji.go today.go yanggongji.go main.go
GOOS=windows GOARCH=amd64 go build -o ccal-gui-win64-v-0.3.6.exe -ldflags="-H windowsgui -s -w" Form1.go Form1Impl.go jieqi24.go about.go dimu.go zeji.go listDay.go rules.go jinji.go today.go yanggongji.go main.go
#GOOS=darwin GOARCH=amd64 go build -o ccal-gui-darwin-today-fix_12h -ldflags="-s -w" Form1.go Form1Impl.go jieqi24.go about.go dimu.go zeji.go listDay.go rules.go jinji.go today.go yanggongji.go main.go
