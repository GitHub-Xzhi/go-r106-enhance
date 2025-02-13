// 全局常量
package main

const (
	// BaseURL 基础 URL，用于拼接其他接口的完整 URL
	BaseURL = "http://192.168.43.1"
	// SmsListURL 短信列表接口
	SmsListURL = BaseURL + "/action/sms_get_sms_list"
	// LoginURL 登录接口
	LoginURL = BaseURL + "/goform/login"
	// GetMgdbParamsURL 获取 MGDB 参数接口
	GetMgdbParamsURL = BaseURL + "/action/get_mgdb_params"
	// WifiSetBasicParamsURL 设置 Wi-Fi 基本参数接口
	WifiSetBasicParamsURL = BaseURL + "/action/wifi_set_basic_params"
)
