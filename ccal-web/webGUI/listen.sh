#!/bin/bash
# 进程监听并重启web服务
pgrep -f ccal-web | wc -l
if [ $? -eq 0 ]; #$?代表上面的返回值
then
        echo "程序没有运行 开始启动..."
        cd /home/xuan/ccal/web/v0.0.6
        echo " "> ./web.log #清空上次log
        echo "`date` restart" >> ./web.log 
        ./ccal-web & >> ./web.log 2>&1  #进程启动命令
fi

