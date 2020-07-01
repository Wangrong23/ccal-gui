#!/bin/bash

go build -o ccal-gui-linux64-v0.3.9 -ldflags="-s -w" Form1.go basic.go jieqi24.go about.go earthMother.go auspicious.go listDay.go rules.go taboo.go today.go yg13.go  main.go
GOOS=windows GOARCH=amd64 go build -o ccal-gui-win64-v0.3.9.exe -ldflags="-H windowsgui -s -w" Form1.go basic.go jieqi24.go about.go earthMother.go auspicious.go listDay.go rules.go taboo.go today.go yg13.go main.go
