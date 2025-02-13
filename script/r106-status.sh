#!/bin/sh
# 状态查询脚本
free_output=$(free)

# 使用 awk 提取 Mem 行的 total、used 和 free，并以 M 为单位输出
echo "$free_output" | awk '/^Mem:/ {
    total=$2 / 1024;  # 将总内存转换为 MB
    used=$3 / 1024;   # 将已用内存转换为 MB
    free=$4 / 1024;   # 将空闲内存转换为 MB
    printf "\n总内存: %.2f M，使用: %.2f M，剩余: %.2f M\n\n", total, used, free
}'

# 检查 r106_enhance 进程
r106_enhance_pid=$(pgrep r106_enhance)
if [ -n "$r106_enhance_pid" ]; then
    r106_enhance_mem=$(cat /proc/$r106_enhance_pid/status | grep -E 'VmRSS' | awk '{printf "%.2f M\n", $2/1024}')
    echo -e "r106_enhance内存: $r106_enhance_mem\n"
else
    echo -e "r106_enhance程序停止\n"
fi

# 检查 frpc 进程
frpc_pid=$(pgrep frpc)
if [ -n "$frpc_pid" ]; then
    frpc_mem=$(cat /proc/$frpc_pid/status | grep -E 'VmRSS' | awk '{printf "%.2f M\n", $2/1024}')
    echo -e "frpc内存: $frpc_mem\n"
else
    echo -e "frpc程序停止\n"
fi

# WiFi状态
wifi_state=$(curl -s http://127.0.0.1:8686/wifi_state)
echo -e "WiFi状态：$wifi_state\n"

# 电池状态
charge_status="否"
battery_online=$(cat /sys/class/power_supply/battery/online)
if [ "$battery_online" -eq 1 ]; then
    charge_status="是"
fi
battery_status=$(cat /sys/class/power_supply/battery/capacity)
echo -e "当前设备电量: $((battery_status - 1))，是否在充电：$charge_status\n"
