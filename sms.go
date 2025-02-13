// 短信转发通知
package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gookit/slog"
	"github.com/wanghuiyt/ding"
	"sort"
	"time"
)

// FetchAllSMS 分页请求所有短信数据
func FetchAllSMS(cache *Cache, cookieManager *CookieManager) bool {
	const (
		pageSize = 10
	)
	start := 1
	end := pageSize
	totalCount := 0

	for {
		requestData := fmt.Sprintf(`{"start":%d,"end":%d,"smsbox":1}`, start, end)
		responseBody, _ := PostRequestWithCookie("sms_get_sms_list请求", Yaml.R106.SmsListURL, requestData,
			cookieManager)

		// 解析 JSON 响应
		var response map[string]interface{}
		if err := json.Unmarshal(responseBody, &response); err != nil {
			slog.Errorf("解析JSON异常：%v", err)
			return false
		}

		// 提取 totalCount 和 smsList
		data, ok := response["data"].(map[string]interface{})
		if !ok {
			slog.Errorf("无效响应结构：缺少“data”字段。")
			return false
		}

		if totalCount == 0 {
			// 仅在第一次请求时提取 totalCount
			totalCount, _ = intFromInterface(data["totalCount"])
		}

		smsList, ok := data["smsList"].([]interface{})
		if !ok {
			slog.Errorf("无效响应结构：缺少“smsList”字段。")
			return false
		}

		// 转换 smsList 为 []map[string]interface{}
		var parsedSMSList []map[string]interface{}
		for _, sms := range smsList {
			if smsMap, ok := sms.(map[string]interface{}); ok {
				parsedSMSList = append(parsedSMSList, smsMap)
			}
		}
		// 处理当前页的短信
		ProcessSMS(parsedSMSList, cache)

		if end >= totalCount {
			break
		}

		start += pageSize
		end += pageSize
		if end > totalCount {
			end = totalCount
		}
	}
	return true
}

// intFromInterface 将接口转换为整数
func intFromInterface(value interface{}) (int, bool) {
	switch v := value.(type) {
	case float64:
		return int(v), true
	case int:
		return v, true
	default:
		return 0, false
	}
}

// ProcessSMS 提取短信的 Phone、Content 和 Datetime 并发送消息
func ProcessSMS(smsList []map[string]interface{}, cache *Cache) {
	// 将 smsList 按 datetime 升序排序
	sort.Slice(smsList, func(i, j int) bool {
		timeI, _ := smsList[i]["datetime"].(float64)
		timeJ, _ := smsList[j]["datetime"].(float64)
		return timeI < timeJ
	})

	// 处理排序后的短信列表
	for _, sms := range smsList {
		// 解析 datetime 并检查是否小于 CacheTTL
		if timestamp, ok := sms["datetime"].(float64); ok {
			if time.Now().Unix()-int64(timestamp) < int64(Yaml.Sms.CacheTTL.Seconds()) {
				phone := sms["phone"].(string)
				content := sms["content"].(string)
				datetime := time.Unix(int64(timestamp), 0).Format("2025-01-01 15:04:05") // 时间戳转换为格式化时间
				text := fmt.Sprintf("%s\n%s\n\n%s\nR106", phone, content, datetime)

				// 生成缓存的 MD5 键
				hash := md5.Sum([]byte(phone + content + datetime))
				cacheKey := hex.EncodeToString(hash[:])

				// 检查是否已缓存
				if !cache.Exists(cacheKey) {
					cache.Add(cacheKey)
					slog.Infof("转发内容：%s", text)
					DingDingNotify(text)
				}
			}
		}
	}
}

// DingDingNotify 钉钉通知
func DingDingNotify(text string) {
	d := ding.Webhook{
		AccessToken: Yaml.DingDing.AccessToken,
		Secret:      Yaml.DingDing.Secret,
	}
	_ = d.SendMessageText(text)
	//fmt.Println(text)
}
