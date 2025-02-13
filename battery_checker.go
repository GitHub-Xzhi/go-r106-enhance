// 电池检测
package main

import (
	"fmt"
	"github.com/gookit/slog"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// StartBatteryCheckByCmd 读取系统文件的电量
func StartBatteryCheckByCmd() {
	go func() {
		for {
			// 获取电量
			cmd := exec.Command("cat", "/sys/class/power_supply/battery/capacity")
			output, err := cmd.Output()
			if err != nil {
				return
			}
			capacityStr := strings.TrimSpace(string(output))
			capacity, err := strconv.Atoi(capacityStr)
			if err != nil {
				return
			}

			// 检查电池是否在充电，1：在充电，0：没有在充电
			cmdOnline := exec.Command("cat", "/sys/class/power_supply/battery/online")
			onlineOutput, err := cmdOnline.Output()
			if err != nil {
				return
			}
			onlineStr := strings.TrimSpace(string(onlineOutput))
			online, err := strconv.Atoi(onlineStr)
			if err != nil {
				return
			}

			var interval = Yaml.Battery.CheckInterval
			currentPower := capacity - 1
			if online == 0 && currentPower <= Yaml.Battery.Low {
				text := fmt.Sprintf("R106随身WiFi电量%d%%过低，请充电", currentPower)
				slog.Warn(text)
				DingDingNotify(text)
				interval = Yaml.Battery.NotifyInterval
			} else {
				interval = Yaml.Battery.CheckInterval
			}
			time.Sleep(interval)
		}
	}()
}
