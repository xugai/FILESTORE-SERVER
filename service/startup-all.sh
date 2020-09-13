#!/bin/bash

#检查service进程
# ps aux: 以用户的格式输出所有在终端机下面运行的程序，不以终端机来区分
# grep -v: 不显示包含-v参数后面内容的行或者文件。-v，取反的意思
# echo -e: 处理特殊字符，类似\n \a \t等等，不把他们当做普通字符处理
check_process() {
  sleep 1
  # shellcheck disable=SC2006
  res=`ps aux | grep -v grep | grep "service/bin" | grep $1`
  if [[ $res != '' ]]; then
    echo -e "\033[32m 已启动服务 \033[0m" "$1_service"
    return 1
  else
    echo -e "\033[31m 启动服务失败 \033[0m" "$1_service"
    return 0
  fi
}

# 编译service可执行文件
build_service() {
  go build -o service/bin/$1 service/$1/main.go
  resbin=`ls service/bin/ | grep $1`
  echo -e "\033[32m 编译完成: \033[0m service/bin/$resbin"
}

# 启动service
# 2>&1 把错误输出重定向到标准输出通道
run_service() {
  echo -e "检查日志文件$1.log是否存在..."
  if [[ ! -f $logpath/$1.log ]]; then
    echo -e "日志文件$1.log不存在，创建该日志文件."
    touch $logpath/$1.log
  else
    echo -e "日志文件$1.log存在."
  fi
  nohup ./service/bin/$1 >> $logpath/$1.log 2>&1 &
  sleep 1
  check_process $1
}

# 创建运行日志目录
logpath=./service/log
#mkdir -p $logpath

# 切换到工程根目录
cd /Users/behe/Desktop/work_station/FILESTORE-SERVER

# 微服务可以用supervisor做进程管理工具；
# 或者也可以通过docker/k8s进行部署
services="
apigw
account
dbproxy
upload
download
"

# 执行编译service
for service in $services
do
  build_service $service
done

# 编译完后，执行启动service
for service in $services
do
  run_service $service
done


echo "微服务启动完毕."