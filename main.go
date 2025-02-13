package main

import (
	"fmt"
	"github.com/gookit/slog"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// 初始化日志配置
	InitLogConfig()

	slog.Info("<<<程序启动>>>")

	// 创建tmp目录
	initTmpDir()

	// 加载配置文件
	LoadConfig()

	// 初始化缓存并加载本地文件
	cache := NewCache()
	cache.LoadFromFile()

	// 初始化 Cookie 管理器
	cookieManager := NewCookieManager()

	// 读取系统文件的电量
	StartBatteryCheckByCmd()

	// 提供WiFi控制的API
	go WiFiControlAPI(cookieManager)

	// 定时请求短信列表
	for {
		if !FetchAllSMS(cache, cookieManager) {
			slog.Error("未能获取短信数据。正在退出...")
			return
		}
		time.Sleep(Yaml.Sms.RequestInterval)
	}
}

// GetProgramDir 获取程序所在的目录
func GetProgramDir() string {
	// 获取程序的可执行文件路径
	exePath, _ := os.Executable()
	return filepath.Dir(exePath)
}

// 创建tmp目录
func initTmpDir() {
	// 定义要检查和创建的目录路径
	dirPath := fmt.Sprintf("%s/tmp", GetProgramDir())
	// 检查目录是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// 如果目录不存在，则创建它
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			slog.Errorf("创建目录失败: %v", err)
			return
		}
	}
}
