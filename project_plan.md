# Web服务器监听工具开发计划

## 项目概述

开发一个基于Golang的Web服务器监听工具，能够监听多个端口，支持HTTPS协议，根据预设规则返回内容。

## 功能需求

1. 监听web服务器，支持监控多个端口，支持HTTPS协议
2. 获取监听的请求域名和路径地址、请求参数等信息，支持GET/POST
3. 根据预设的规则进行返回请求的内容
4. 预设规则文件为 xxxxx.json（xxxxx是请求的域名），根据解析该json里面的路径信息和参数信息进行匹配请求项目，然后返回内容
5. 监听的端口配置文件为 config.json，该文件可配置监听的端口列表
6. 当收到请求时，判断请求的域名和参数信息，查找是否有对应的规则文件
7. 如果有规则文件，则从规则文件中读取规则进行匹配并返回
8. 如果没有规则文件，则创建该规则文件，并设置默认值
9. 如果有规则文件但没有匹配到规则，则返回默认规则

## 项目结构

```
/
├── main.go                 # 主程序入口
├── config.json             # 端口配置文件
├── config/                 # 配置相关
│   └── config.go           # 配置处理
├── server/                 # 服务器相关
│   ├── server.go           # 服务器实现
│   └── https.go            # HTTPS支持
├── handler/                # 请求处理
│   ├── handler.go          # 请求处理器
│   └── rule.go             # 规则处理
├── model/                  # 数据模型
│   ├── config.go           # 配置模型
│   └── rule.go             # 规则模型
└── utils/                  # 工具函数
    └── utils.go            # 通用工具函数
```

## 开发流程

### 1. 项目初始化 ✅
- 创建项目目录结构
- 初始化Go模块

### 2. 配置模块开发
- 实现配置文件读取功能
- 实现配置模型结构

### 3. 规则模块开发
- 实现规则文件读取功能
- 实现规则匹配功能
- 实现规则文件创建功能

### 4. 服务器模块开发
- 实现HTTP服务器
- 实现HTTPS支持
- 实现多端口监听

### 5. 请求处理模块开发
- 实现请求参数解析
- 实现域名和路径解析
- 实现规则匹配和响应

### 6. 主程序开发
- 实现主程序入口
- 整合各模块功能

### 7. 测试与优化
- 单元测试
- 集成测试
- 性能优化

## 技术栈

- 语言：Golang
- 依赖：
  - 标准库 net/http：HTTP服务器
  - 标准库 encoding/json：JSON处理
  - 标准库 crypto/tls：TLS/HTTPS支持

## 数据结构设计

### 配置文件结构 (config.json)

```json
{
  "ports": [
    {
      "port": 8080,
      "https": false
    },
    {
      "port": 8443,
      "https": true,
      "cert": "path/to/cert.pem",
      "key": "path/to/key.pem"
    }
  ],
  "defaultResponse": "没有匹配的规则"
}
```

### 规则文件结构 (example.com.json)

```json
{
  "rules": [
    {
      "path": "/api/users",
      "method": "GET",
      "params": {
        "id": "123"
      },
      "response": {
        "status": 200,
        "headers": {
          "Content-Type": "application/json"
        },
        "body": "{\"name\":\"John\",\"age\":30}"
      }
    },
    {
      "path": "/api/login",
      "method": "POST",
      "response": {
        "status": 200,
        "headers": {
          "Content-Type": "application/json"
        },
        "body": "{\"token\":\"abc123\"}"
      }
    }
  ],
  "default": {
    "status": 404,
    "headers": {
      "Content-Type": "text/plain"
    },
    "body": "没有匹配的规则"
  }
}
```