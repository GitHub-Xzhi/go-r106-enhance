// yaml配置
package main

import (
	"fmt"
	"github.com/gookit/slog"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
	"strings"
	"time"
)

// Config 配置文件的结构
type Config struct {
	ApiPort string `yaml:"ApiPort"` // API端口
	// 短信相关配置
	Sms struct {
		CacheTTL          time.Duration `yaml:"CacheTTL"`          // 缓存有效时间
		FileCleanupPeriod time.Duration `yaml:"FileCleanupPeriod"` // 定时检查cache.json文件中过期记录并清理
		RequestInterval   time.Duration `yaml:"RequestInterval"`   // 定时请求短信列表间隔时间
	} `yaml:"Sms"`
	// 电池相关配置
	Battery struct {
		CheckInterval  time.Duration `yaml:"CheckInterval"`  // 电量检测时间间隔
		Low            int           `yaml:"Low"`            // 低电量阈值
		NotifyInterval time.Duration `yaml:"NotifyInterval"` // 低电量通知时间间隔
	} `yaml:"Battery"`
	// 外部服务相关配置
	R106 struct {
		UserName string `yaml:"UserName"` // 用户名
		Password string `yaml:"Password"` // 密码
	} `yaml:"R106"`
	// 钉钉相关配置
	DingDing struct {
		AccessToken string `yaml:"AccessToken"` // 钉钉访问令牌
		Secret      string `yaml:"Secret"`      // 钉钉密钥
	} `yaml:"DingDing"`
}

var Yaml *Config // 全局配置变量

// LoadConfig 从指定的 YAML 文件加载配置，并解析占位符
func LoadConfig() {
	fileContent, err := os.ReadFile(fmt.Sprintf("%s/config/config.yaml", GetProgramDir()))
	if err != nil {
		slog.Errorf("打开配置文件出错: %v", err)
	}
	content := string(fileContent)

	placeholders := extractPlaceholders(content)
	content = replacePlaceholders(content, placeholders)

	// 将处理后的内容解析为 Config 结构体
	config := &Config{}
	decoder := yaml.NewDecoder(strings.NewReader(content))
	if err := decoder.Decode(config); err != nil {
		slog.Errorf("解析配置文件出错: %v", err)
	}
	Yaml = config
}

// extractPlaceholders 从 YAML 文件内容中提取所有键值对作为占位符
func extractPlaceholders(content string) map[string]string {
	placeholders := make(map[string]string)

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		// 跳过空行和注释
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析键值对（假设 YAML 中的键值对格式为 key: value）
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// 去掉 value 两侧的引号（如果有）
			value = strings.Trim(value, `"'`)
			placeholders[key] = value
		}
	}

	return placeholders
}

// replacePlaceholders 替换字符串中的所有占位符
func replacePlaceholders(content string, placeholders map[string]string) string {
	// 使用正则表达式匹配 ${key} 格式的占位符
	re := regexp.MustCompile(`\$\{([a-zA-Z0-9_]+)\}`)

	return re.ReplaceAllStringFunc(content, func(match string) string {
		// 提取占位符中的键名
		key := match[2 : len(match)-1] // 去掉 ${ 和 } 部分
		if value, exists := placeholders[key]; exists {
			return value
		}
		return match // 如果未找到对应值，保留原占位符
	})
}
