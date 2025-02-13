#!/bin/bash

# 检查是否提供了参数
if [ "$1" == "1" ]; then
    # 如果参数是1，则关闭无线网络接口
    return_value=$(curl -s http://127.0.0.1:8686/disable_wifi)
else
    # 如果没有提供参数或参数不是1，则开启无线网络接口
    return_value=$(curl -s http://127.0.0.1:8686/enable_wifi)
fi

echo -e "\n$return_value\n"