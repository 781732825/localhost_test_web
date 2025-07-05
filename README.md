# Web服务器监听工具

这是一个基于Golang的Web服务器监听工具，能够监听多个端口，支持HTTPS协议，根据预设规则返回内容。

## 功能特点

1. 监听web服务器，支持监控多个端口，支持HTTPS协议
2. 获取监听的请求域名和路径地址、请求参数等信息，支持GET/POST
3. 根据预设的规则进行返回请求的内容
4. 预设规则文件为 xxxxx.json（xxxxx是请求的域名），根据解析该json里面的路径信息和参数信息进行匹配请求项目，然后返回内容
5. 自动创建不存在的规则文件，并设置默认规则

## 安装

```bash
# 克隆仓库
git clone https://github.com/yourusername/web-server-listener.git

# 进入项目目录
cd web-server-listener

# 编译
go build -o host
```

## 使用方法

### 配置文件

配置文件为`config.json`，用于配置监听的端口列表：

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

### 规则文件

规则文件存放在`rules`目录下，文件名为域名加`.json`后缀，例如`example.com.json`：

```json
{
  "rules": [
    {
      "path": "/api/users",
      "method": "GET",
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

### 启动服务

```bash
# 使用默认配置启动
./host

# 指定配置文件
./host -config=/path/to/config.json

# 指定规则目录
./host -rules=/path/to/rules
```

## 工作原理

1. 服务启动后，会根据配置文件监听指定的端口
2. 当收到请求时，会解析请求的域名、路径和参数
3. 根据域名查找对应的规则文件
4. 如果规则文件不存在，会创建一个默认的规则文件
5. 根据请求的路径和参数匹配规则
6. 如果匹配成功，返回规则中定义的响应
7. 如果没有匹配的规则，返回默认响应

## 许可证

MIT