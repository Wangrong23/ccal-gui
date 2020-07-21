#!/bin/bash

go build -o ccal-govcl-linux64-v0.6.0 -ldflags="-s -w" Form1.go basic.go jieqi24.go about.go earthMother.go auspicious.go listDay.go rules.go taboo.go today.go yg13.go  qimen.go main.go
GOOS=windows GOARCH=amd64 go build -o ccal-govcl-win64-v0.6.0.exe -ldflags="-H windowsgui -s -w" Form1.go basic.go jieqi24.go about.go earthMother.go auspicious.go listDay.go rules.go taboo.go today.go yg13.go qimen.go main.go
