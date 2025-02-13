// Cookie 管理器
package main

import (
	"fmt"
	"github.com/gookit/slog"
	"io/ioutil"
	"os"
	"strings"
)

// CookieManager 管理 Cookie 的结构
type CookieManager struct {
	file string
}

// NewCookieManager 创建一个新的 Cookie 管理器
func NewCookieManager() *CookieManager {
	cookieFile := fmt.Sprintf("%s/tmp/cookie.txt", GetProgramDir())
	return &CookieManager{file: cookieFile}
}

// Save 保存 Cookie 到本地文件
func (cm *CookieManager) Save(cookie string) error {
	err := ioutil.WriteFile(cm.file, []byte(cookie), 0644)
	if err != nil {
		slog.Errorf("将Cookie保存到文件时出错:%v", err)
	}
	return err
}

// Load 从本地文件加载 Cookie
func (cm *CookieManager) Load() (string, error) {
	data, err := ioutil.ReadFile(cm.file)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在的情况不认为是错误，返回空字符串
			return "", nil
		}
		slog.Errorf("从文件加载cookie出错:%v", err)
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
