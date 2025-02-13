// 缓存管理
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gookit/slog"
	"os"
	"sync"
	"time"
)

// Cache 缓存结构
type Cache struct {
	Data  map[string]time.Time
	Mutex sync.Mutex
	File  string
}

// NewCache 新建缓存对象
func NewCache() *Cache {
	cache := &Cache{
		Data: make(map[string]time.Time),
		File: fmt.Sprintf("%s/tmp/cache.json", GetProgramDir()),
	}
	go cache.StartFileCleanupTask() // 启动定时清理任务
	return cache
}

// Add 添加缓存记录并设置过期时间
func (c *Cache) Add(key string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Data[key] = time.Now().Add(Yaml.Sms.CacheTTL)
	c.SaveToFile()
}

// Exists 检查缓存记录是否存在且未过期
func (c *Cache) Exists(key string) bool {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	expiration, found := c.Data[key]
	if !found {
		return false
	}
	if time.Now().After(expiration) {
		delete(c.Data, key)
		c.SaveToFile()
		return false
	}
	return true
}

// SaveToFile 保存缓存到本地文件（JSON 格式）
func (c *Cache) SaveToFile() {
	file, err := os.Create(c.File)
	if err != nil {
		slog.Errorf("将缓存保存到文件时出错:%v", err)
		return
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(c.Data); err != nil {
		slog.Errorf("缓存数据编码错误:%v", err)
	}
}

// LoadFromFile 从本地文件加载缓存（JSON 格式）
func (c *Cache) LoadFromFile() {
	file, err := os.Open(c.File)
	if err != nil {
		if !os.IsNotExist(err) {
			slog.Errorf("从文件加载缓存时出错:%v", err)
		}
		return
	}
	defer file.Close()

	// 从 JSON 文件解码数据
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&c.Data); err != nil {
		slog.Errorf("缓存数据解码出错:%v", err)
	}

	// 删除过期记录
	c.cleanupExpired()
	c.SaveToFile()
}

// cleanupExpired 清理内存中过期的记录
func (c *Cache) cleanupExpired() {
	now := time.Now()
	for key, expiration := range c.Data {
		if now.After(expiration) {
			delete(c.Data, key)
		}
	}
}

// StartFileCleanupTask 定时清理本地文件中过期记录
func (c *Cache) StartFileCleanupTask() {
	ticker := time.NewTicker(Yaml.Sms.FileCleanupPeriod)
	defer ticker.Stop()

	for range ticker.C {
		c.Mutex.Lock()
		c.cleanupExpired()
		c.SaveToFile()
		c.Mutex.Unlock()
	}
}
