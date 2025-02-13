// WiFi控制
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// WiFiConfig 定义 WiFi 配置的结构
type WiFiConfig struct {
	WifiState     string `json:"wifi_state_0"`
	WifiSSID      string `json:"wifi_ssid_0"`
	WifiSecurity  string `json:"wifi_security_0"`
	WifiPSK       string `json:"wifi_psk_0"`
	WifiBroadcast string `json:"wifi_broadcast_ssid_0"`
}

// sendRequest 发送 WiFi 开关请求
func sendRequest(state string, cookieManager *CookieManager) error {
	// 构造请求数据
	requestBody := WiFiConfig{
		WifiState:     state,
		WifiSSID:      "NBB",
		WifiSecurity:  "wpa_wpa2",
		WifiPSK:       "Sz6z-Xzhi-%e7@iK",
		WifiBroadcast: "visible",
	}
	_, err := PostRequestWithCookie("WiFi控制", Yaml.R106.WifiSetBasicParamsURL, requestBody, cookieManager)
	if err != nil {
		return fmt.Errorf("请求异常：%v", err)
	}
	return nil
}

func wifiIsEnable(cookieManager *CookieManager) bool {
	body := `{
		"keys": [
			"wifi_state_0"
		]
	}`
	request, _ := PostRequestWithCookie("获取WiFi状态", Yaml.R106.GetMgdbParamsURL, body, cookieManager)
	// 解析JSON响应
	var result struct {
		Retcode int `json:"retcode"`
		Data    struct {
			WifiState0 string `json:"wifi_state_0"`
		} `json:"data"`
	}
	json.Unmarshal(request, &result)
	if result.Data.WifiState0 == "ap_enable" {
		return true
	}
	return false
}

// enableWiFiHandler 开启 WiFi 操作
func enableWiFiHandler(w http.ResponseWriter, cookieManager *CookieManager) {
	if wifiIsEnable(cookieManager) {
		w.Write([]byte("WiFi已开启，无需重复操作"))
		return
	}
	if err := sendRequest("ap_enable", cookieManager); err != nil {
		http.Error(w, fmt.Sprintf("开启WiFi失败: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("远程开启WiFi"))
}

// disableWiFiHandler 关闭 WiFi 操作
func disableWiFiHandler(w http.ResponseWriter, cookieManager *CookieManager) {
	if !wifiIsEnable(cookieManager) {
		w.Write([]byte("WiFi已关闭，无需重复操作"))
		return
	}
	if err := sendRequest("ap_disable", cookieManager); err != nil {
		http.Error(w, fmt.Sprintf("关闭WiFi失败: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("远程关闭WiFi"))
}

// WiFiControlAPI 对外提供WiFi控制的API
func WiFiControlAPI(cookieManager *CookieManager) {
	http.HandleFunc("/wifi_state", func(w http.ResponseWriter, r *http.Request) {
		msg := "关闭"
		if wifiIsEnable(cookieManager) {
			msg = "开启"
		}
		w.Write([]byte(msg))
	})
	http.HandleFunc("/enable_wifi", func(w http.ResponseWriter, r *http.Request) {
		enableWiFiHandler(w, cookieManager)
	})
	http.HandleFunc("/disable_wifi", func(w http.ResponseWriter, r *http.Request) {
		disableWiFiHandler(w, cookieManager)
	})
	http.ListenAndServe(Yaml.ApiPort, nil)
}
