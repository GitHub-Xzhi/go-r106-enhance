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

//// ResponseData 定义响应结构体
//type ResponseData struct {
//	DeviceBatteryChargeStatus string `json:"device_battery_charge_status"`
//	DeviceBatteryLevelPercent string `json:"device_battery_level_percent"`
//}
//
//type ResponsePayload struct {
//	Retcode int          `json:"retcode"`
//	Data    ResponseData `json:"data"`
//}
//
//// StartBatteryCheck 启动电量检测任务
//func StartBatteryCheck(cookie string, cookieManager *CookieManager) {
//	go func() {
//		for {
//			// 调用统一方法获取并解析电量状态
//			batteryLevel, isLow, err := fetchBatteryStatus(cookie, cookieManager)
//			if err != nil {
//				continue
//			}
//
//			var interval = Yaml.BatteryCheckInterval
//			if isLow {
//				text := fmt.Sprintf("R106随身WiFi电量%d%%过低，请充电", batteryLevel)
//				slog.Warn(text)
//				DingDingNotify(text)
//				interval = 10 * time.Minute
//			} else {
//				interval = Yaml.BatteryCheckInterval
//			}
//			time.Sleep(interval)
//		}
//	}()
//}
//
//// Deprecated
//func sendPostRequest(url string, cookie string, body []byte) ([]byte, error) {
//	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
//	if err != nil {
//		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
//	}
//
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Cookie", cookie)
//
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
//	}
//	defer resp.Body.Close()
//
//	respBody, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, fmt.Errorf("failed to read response body: %w", err)
//	}
//
//	return respBody, nil
//}
//
//// Deprecated
//// 解析响应数据
//func parseResponse(body []byte) (ResponsePayload, error) {
//	var responsePayload ResponsePayload
//	err := json.Unmarshal(body, &responsePayload)
//	if err != nil {
//		return ResponsePayload{}, fmt.Errorf("failed to unmarshal response: %w", err)
//	}
//	return responsePayload, nil
//}
//
//// Deprecated
//// 检查电池状态是否低
//func isBatteryLow(responsePayload ResponsePayload) (int, bool, error) {
//	if responsePayload.Retcode != 0 {
//		return -1, false, fmt.Errorf("error in response: retcode = %d", responsePayload.Retcode)
//	}
//
//	chargeStatus := responsePayload.Data.DeviceBatteryChargeStatus
//	batteryLevel := responsePayload.Data.DeviceBatteryLevelPercent
//
//	batteryLevelInt, err := strconv.Atoi(batteryLevel)
//	if err != nil {
//		return -1, false, fmt.Errorf("failed to parse battery level: %v", err)
//	}
//
//	// 检查条件
//	if chargeStatus == "none" && batteryLevelInt <= Yaml.LowBatteryValue {
//		return batteryLevelInt, true, nil
//	}
//
//	return batteryLevelInt, false, nil
//}
//
//// StartBatteryStatusAPI 对外提供查询电量状态的API
//func StartBatteryStatusAPI(cookie string, cookieManager *CookieManager) {
//	http.HandleFunc("/api/battery-status", func(w http.ResponseWriter, r *http.Request) {
//		if r.Method != http.MethodGet {
//			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//			return
//		}
//		batteryLevel, _, err := fetchBatteryStatus(cookie, cookieManager)
//		if err != nil {
//			return
//		}
//		w.Header().Set("Content-Type", "text/plain")
//		fmt.Fprintf(w, "%d", batteryLevel)
//	})
//	http.ListenAndServe(":8686", nil)
//}
//
//// Deprecated
//// fetchBatteryStatus 获取并解析电量状态数据
//func fetchBatteryStatus(cookie string, cookieManager *CookieManager) (int, bool, error) {
//	body := []byte(`{
//		"keys": [
//			"device_battery_charge_status",
//			"device_battery_level_percent"
//		]
//	}`)
//	respBody, err := sendPostRequest(Yaml.GetMgdbParamsURL, cookie, body)
//
//	// 检查响应是否有效
//	if !json.Valid(respBody) {
//		time.Sleep(5 * time.Second)
//		cookieTmp := cookie
//		cookie, _ = cookieManager.Load()
//		slog.Warnf("电量检测，cookie失效，%s，尝试从本地文件加载Cookie：%s", cookieTmp, cookie)
//		// 重新发送请求
//		respBody, _ = sendPostRequest(Yaml.GetMgdbParamsURL, cookie, body)
//	}
//	if err != nil {
//		return -1, false, err
//	}
//
//	// 解析并校验响应
//	responsePayload, err := parseResponse(respBody)
//	if err != nil {
//		return -1, false, err
//	}
//
//	// 检查电量状态
//	batteryLevel, isLow, err := isBatteryLow(responsePayload)
//	if err != nil {
//		return -1, false, err
//	}
//
//	return batteryLevel, isLow, nil
//}
