#!/bin/bash

go build -o ccal-gui -ldflags="-s -w" Form1.go Form1Impl.go jieqi24.go about.go dimu.go zeji.go listDay.go main.go
