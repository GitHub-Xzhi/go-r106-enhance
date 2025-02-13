#!/bin/sh
# r106_enhance程序监控脚本
PROGRAM="/home/root/r106-enhance/r106_enhance_linux"
LOG_FILE="/home/root/r106-enhance/logs/r106-enhance-monitor.log"
MAX_LOG_SIZE=2 # 单位M

chmod 755 $PROGRAM

# 日志文件大小检查函数
check_log_size() {
  LOG_SIZE=$(du -m "$LOG_FILE" | cut -f1)
  if [ "$LOG_SIZE" -gt "$MAX_LOG_SIZE" ]; then
    # 清空日志文件
    echo "$(date): $LOG_FILE 大于 ${MAX_LOG_SIZE}M, 开始清理..." >"$LOG_FILE"
  fi
}

sleep 5
while true; do
  # 检查程序是否正在运行
  if ! pgrep -f "$PROGRAM" >/dev/null; then
    # 启动程序并获取PID
    $PROGRAM >/dev/null 2>&1 &
    PROGRAM_PID=$!
    echo "$(date): $PROGRAM程序没有在运行，开始启动(PID: $PROGRAM_PID)..." >>"$LOG_FILE"
  fi

  # 检查日志文件大小
  check_log_size
  sleep 5
done
