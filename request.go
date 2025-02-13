// 网络请求
package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gookit/slog"
	"io/ioutil"
	"net/http"
	"strings"
)

// LoginAndGetCookie 登录并获取新的Cookie
func LoginAndGetCookie(cookieManager *CookieManager) string {
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	userName := HmacMD5(Yaml.R106.UserName)
	password := HmacMD5(Yaml.R106.Password)
	loginData := fmt.Sprintf(`{"username":"%s","password":"%s"}`, userName, password)
	resp, respByte, err := PostRequest(Yaml.R106.LoginURL, headers, loginData)
	if err != nil {
		slog.Errorf("请求异常：%v", err)
		return ""
	}
	var result struct {
		Retcode    int `json:"retcode"`
		RemainSecs int `json:"remain_secs"`
	}
	json.Unmarshal(respByte, &result)
	if result.Retcode != 0 {
		slog.Warnf("请求cookie失败，返回值：%s", string(respByte))
		msg := ""
		secs := result.RemainSecs
		if secs > 0 {
			msg = fmt.Sprintf("，请%d秒后重试", secs)
		}
		slog.Errorf("用户名或密码错误%s", msg)
		return ""
	}
	// 提取 Set-Cookie 响应头
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if strings.Contains(cookie.Name, "-webs-session-") {
			slog.Infof("登录获取新的Cookie：%s", cookies)
			DingDingNotify("R106 Cookie失效，重新获取")
			newCookie := fmt.Sprintf("%s=%s", cookie.Name, cookie.Value)
			if newCookie != "" {
				// 保存新的 Cookie 到本地文件
				err := cookieManager.Save(newCookie)
				if err != nil {
					slog.Errorf("新的Cookie保存到文件异常：%v", err)
				}
			} else {
				slog.Error("未能从登录中获得新的cookie")
				return ""
			}
			return newCookie
		}
	}
	return ""
}

// PostRequestWithCookie POST请求带cookie
func PostRequestWithCookie(msg string, url string, requestBody interface{}, cookieManager *CookieManager) ([]byte, error) {
	// 从本地文件加载 Cookie
	cookie, err := cookieManager.Load()
	if err != nil {
		return nil, fmt.Errorf("加载cookie失败: %w", err)
	}

	// 尝试请求，最多重试2次
	for retry := 0; retry < 3; retry++ {
		headers := map[string]string{
			"Cookie":       cookie,
			"Content-Type": "application/json",
		}
		_, respBody, err := PostRequest(url, headers, requestBody)
		respBodyStr := string(respBody)
		if err != nil {
			slog.Errorf("请求异常：%v", err)
			return nil, fmt.Errorf("请求异常: %s", respBodyStr)
		}
		// 检查响应内容是否包含特定的 HTML 格式，表示 Cookie 失效
		if strings.Contains(respBodyStr, "<html>") {
			if retry < 2 {
				slog.Warnf("%s，cookie失效，%s，返回值：\n %v \n", msg, cookie, respBodyStr)
				cookie = LoginAndGetCookie(cookieManager)
				slog.Warnf("第%d次尝试登录获取新的Cookie：%s", retry+1, cookie)
				continue
			} else {
				return nil, fmt.Errorf("已经重试过，仍然失败: %s", respBodyStr)
			}
		}
		return respBody, nil
	}
	return nil, fmt.Errorf("重试后请求失败")
}

// PostRequest POST请求
func PostRequest(url string, headers map[string]string, requestBody interface{}) (*http.Response, []byte, error) {
	// 将 requestBody 转换为 []byte
	var body []byte
	var err error
	switch v := requestBody.(type) {
	case []byte:
		body = v
	case string:
		body = []byte(v)
	case int, int8, int16, int32, int64, float32, float64:
		body = []byte(fmt.Sprintf("%v", v))
	case map[string]interface{}, map[string]string, []interface{}, []string: // 如果是 JSON 结构
		body, err = json.Marshal(v)
		if err != nil {
			return nil, nil, fmt.Errorf("序列化请求体失败: %w", err)
		}
	default:
		// 尝试对自定义结构体类型进行 JSON 序列化
		body, err = json.Marshal(v)
		if err != nil {
			return nil, nil, fmt.Errorf("序列化自定义结构体失败: %w", err)
		}
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}
	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	// 发送请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("请求发送失败: %w", err)
	}
	// 确保响应体被正确关闭
	respBody, err := func() ([]byte, error) {
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("读取响应体失败: %w", err)
		}
		return respBody, nil
	}()
	if err != nil {
		return nil, nil, err
	}
	return resp, respBody, nil
}

// HmacMD5 生成 HMAC-MD5 的十六进制字符串
func HmacMD5(data string) string {
	// 使用 HMAC-MD5 加密
	key := "0123456789"
	h := hmac.New(md5.New, []byte(key))
	h.Write([]byte(data))
	// 返回十六进制字符串
	return hex.EncodeToString(h.Sum(nil))
}
