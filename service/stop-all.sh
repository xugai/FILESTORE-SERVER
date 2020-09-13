#!/bin/bash

stop_process() {
  sleep 1
  pid=`ps aux | grep -v grep | grep "service/bin" | grep $1 | awk '{print $2}'`
  if [[ pid != '' ]]; then
    ps aux | grep -v grep | grep "service/bin" | grep $1 | awk '{print $2}' | xargs kill
    echo -e "\033[32m已关闭服务: \033[0m" "$1_service"
    return 1
  else
    echo -e "\033[31m并未启动服务: \033[0m" "$1_service"
    return 0
  fi
}

services="
apigw
account
dbproxy
upload
download
"

# 关闭service
for service in $services
do
  stop_process $service
done

echo "微服务已全部关闭."
