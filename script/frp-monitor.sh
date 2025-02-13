#!/bin/sh
# frp程序监控脚本：指定时间段实时开启和关闭
PROGRAM_PATH="/home/root/frp"
PROGRAM="$PROGRAM_PATH/frpc"
LOG_FILE="$PROGRAM_PATH/frp-monitor.log"
CONFIG_FILE="$PROGRAM_PATH/frpc.toml"
MAX_LOG_SIZE=2 # 单位M

# 定义时间范围（小时）
START_CLOSE=1
END_CLOSE=7
START_OPEN=9
END_OPEN=22
# 检查间隔，单位秒
CHECK_INTERVAL=60

# 获取程序PID的函数
get_program_pid() {
  pidof "$PROGRAM"
}

# 日志文件大小检查函数
check_log_size() {
  LOG_SIZE=$(du -m "$LOG_FILE" | cut -f1)
  if [ "$LOG_SIZE" -gt "$MAX_LOG_SIZE" ]; then
    # 清空日志文件
    echo "$(date): $LOG_FILE 大于 ${MAX_LOG_SIZE}M, 开始清理..." >"$LOG_FILE"
  fi
}

while true; do
  # 获取当前小时
  CURRENT_HOUR=$(date +"%H")

  # 检查当前时间是否在关闭程序的时间范围内
  if [ "$CURRENT_HOUR" -ge "$START_CLOSE" ] && [ "$CURRENT_HOUR" -lt "$END_CLOSE" ]; then
    PROGRAM_PID=$(get_program_pid)
    # 检查程序是否正在运行
    if [ -n "$PROGRAM_PID" ]; then
      # 关闭程序
      echo "$(date): $START_CLOSE点~$END_CLOSE点时间段关闭$PROGRAM(PID: $PROGRAM_PID)..." >>"$LOG_FILE"
      pkill "$PROGRAM"
    fi
  # 检查当前时间是否在启动程序的时间范围内
  elif [ "$CURRENT_HOUR" -ge "$START_OPEN" ] && [ "$CURRENT_HOUR" -le "$END_OPEN" ]; then
    # 检查程序是否未运行
    if ! get_program_pid >/dev/null; then
      # 启动程序并获取PID
      $PROGRAM -c $CONFIG_FILE >/dev/null 2>&1 &
      PROGRAM_PID=$!
      echo "$(date): $START_OPEN点~$END_OPEN点时间段启动$PROGRAM(PID: $PROGRAM_PID)..." >>"$LOG_FILE"
    fi
  fi

  # 检查日志文件大小
  check_log_size
  sleep "$CHECK_INTERVAL"
done
