package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/trae/host/model"
)

// RuleManager 管理规则文件的加载和匹配
type RuleManager struct {
	RuleDir         string                      // 规则文件存储目录
	DefaultResponse string                      // 默认响应内容
	ruleCache       map[string]*model.RuleFile  // 规则缓存，避免频繁读取文件
	lastUpdate      time.Time                   // 上次更新缓存的时间
	cacheMutex      sync.RWMutex                // 缓存读写锁
	updateInterval  time.Duration               // 缓存更新间隔
}

// NewRuleManager 创建规则管理器
func NewRuleManager(ruleDir string, defaultResponse string) *RuleManager {
	// 确保规则目录存在
	if _, err := os.Stat(ruleDir); os.IsNotExist(err) {
		os.MkdirAll(ruleDir, 0755)
		log.Printf("创建规则目录: %s", ruleDir)
	}

	// 创建规则管理器
	rm := &RuleManager{
		RuleDir:         ruleDir,
		DefaultResponse: defaultResponse,
		ruleCache:       make(map[string]*model.RuleFile),
		lastUpdate:      time.Now(),
		updateInterval:  10 * time.Minute, // 设置缓存更新间隔为10分钟
	}

	// 初始加载规则文件到缓存
	rm.updateCache()

	log.Printf("规则管理器初始化完成，规则目录: %s，已加载规则文件到缓存，缓存更新间隔: %v", ruleDir, rm.updateInterval)
	return rm
}

// GetRuleFile 获取指定域名的规则文件
func (rm *RuleManager) GetRuleFile(domain string) (*model.RuleFile, error) {
	// 检查是否需要更新缓存
	rm.checkAndUpdateCache()

	// 先从缓存中查找
	rm.cacheMutex.RLock()
	ruleFile, exists := rm.ruleCache[domain]
	rm.cacheMutex.RUnlock()

	if exists {
		log.Printf("从缓存中获取规则文件: %s，包含 %d 条规则", domain, len(ruleFile.Rules))
		return ruleFile, nil
	}

	// 缓存中不存在，加载规则文件
	filePath := filepath.Join(rm.RuleDir, domain+".json")
	log.Printf("缓存中不存在规则文件，加载规则文件: %s", filePath)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("规则文件不存在，创建默认规则文件: %s", filePath)
		// 如果规则文件不存在，创建默认规则文件
		return rm.createDefaultRuleFile(domain)
	}

	// 读取规则文件
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取规则文件失败: %v", err)
	}

	// 解析JSON
	var newRuleFile model.RuleFile
	if err := json.Unmarshal(data, &newRuleFile); err != nil {
		return nil, fmt.Errorf("解析规则文件失败: %v", err)
	}

	// 编译规则中的正则表达式
	for i := range newRuleFile.Rules {
		if err := newRuleFile.Rules[i].CompileRegex(); err != nil {
			return nil, fmt.Errorf("编译正则表达式失败: %v", err)
		}
	}

	// 更新缓存
	rm.cacheMutex.Lock()
	rm.ruleCache[domain] = &newRuleFile
	rm.cacheMutex.Unlock()

	log.Printf("成功加载规则文件并更新缓存: %s，包含 %d 条规则", filePath, len(newRuleFile.Rules))
	return &newRuleFile, nil
}

// createDefaultRuleFile 创建默认规则文件
func (rm *RuleManager) createDefaultRuleFile(domain string) (*model.RuleFile, error) {
	// 创建默认规则文件
	ruleFile := &model.RuleFile{
		Rules: []model.Rule{},
		Default: model.ResponseRule{
			Status:  404,
			Headers: map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			Body:    rm.DefaultResponse,
		},
	}

	// 将规则写入文件
	filePath := filepath.Join(rm.RuleDir, domain+".json")
	data, err := json.MarshalIndent(ruleFile, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("序列化规则失败: %v", err)
	}

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return nil, fmt.Errorf("写入规则文件失败: %v", err)
	}

	// 更新缓存
	rm.cacheMutex.Lock()
	rm.ruleCache[domain] = ruleFile
	rm.cacheMutex.Unlock()

	log.Printf("创建默认规则文件并更新缓存: %s", filePath)
	return ruleFile, nil
}

// MatchRule 匹配请求与规则
func (rm *RuleManager) MatchRule(req *http.Request) (*model.ResponseRule, error) {
	// 获取域名
	domain := req.Host
	// 移除端口号（如果有）
	if strings.Contains(domain, ":") {
		domain = strings.Split(domain, ":")[0]
	}

	// 获取规则文件（会自动检查缓存是否需要更新）
	ruleFile, err := rm.GetRuleFile(domain)
	if err != nil {
		return nil, err
	}

	// 匹配规则
	for i, rule := range ruleFile.Rules {
		if rule.MatchRequest(req) {
			log.Printf("匹配成功: 规则 #%d [%s %s]", i+1, rule.Method, rule.Path)
			return &rule.Response, nil
		}
	}

	// 如果没有匹配的规则，返回默认响应
	log.Printf("未匹配到规则，使用默认响应")
	return &ruleFile.Default, nil
}

// checkAndUpdateCache 检查是否需要更新缓存，如果需要则更新
func (rm *RuleManager) checkAndUpdateCache() {
	rm.cacheMutex.RLock()
	timeElapsed := time.Since(rm.lastUpdate)
	rm.cacheMutex.RUnlock()

	// 如果距离上次更新时间超过更新间隔，则更新缓存
	if timeElapsed >= rm.updateInterval {
		rm.updateCache()
	}
}

// updateCache 更新规则文件缓存
func (rm *RuleManager) updateCache() {
	rm.cacheMutex.Lock()
	defer rm.cacheMutex.Unlock()

	// 更新上次更新时间
	rm.lastUpdate = time.Now()

	// 清空缓存
	rm.ruleCache = make(map[string]*model.RuleFile)

	// 读取规则目录下的所有规则文件
	files, err := ioutil.ReadDir(rm.RuleDir)
	if err != nil {
		log.Printf("读取规则目录失败: %v", err)
		return
	}

	// 加载所有规则文件到缓存
	loadedCount := 0
	for _, file := range files {
		// 跳过目录和非json文件
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		// 获取域名（文件名去掉.json后缀）
		domain := strings.TrimSuffix(file.Name(), ".json")

		// 读取规则文件
		filePath := filepath.Join(rm.RuleDir, file.Name())
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Printf("读取规则文件失败: %s, %v", filePath, err)
			continue
		}

		// 解析JSON
		var ruleFile model.RuleFile
		if err := json.Unmarshal(data, &ruleFile); err != nil {
			log.Printf("解析规则文件失败: %s, %v", filePath, err)
			continue
		}

		// 编译规则中的正则表达式
		for i := range ruleFile.Rules {
			if err := ruleFile.Rules[i].CompileRegex(); err != nil {
				log.Printf("编译正则表达式失败: %s, %v", filePath, err)
				continue
			}
		}

		// 更新缓存
		rm.ruleCache[domain] = &ruleFile
		loadedCount++
	}

	log.Printf("缓存更新完成，已加载 %d 个规则文件，下次更新时间: %v", loadedCount, rm.lastUpdate.Add(rm.updateInterval))
}