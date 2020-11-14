#!/bin/bash
# 进程监听并重启web服务
/usr/bin/kill `pgrep ccal-web`
cd /home/xuan/ccal/web
./ccal-web & > ./web.log 2>&1  #进程启动命令
