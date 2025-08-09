# 🚀 go-wrk压测快速开始指南

## 📋 安装go-wrk

### Windows
```bash
go install github.com/tsliwowicz/go-wrk@latest
```

### Linux/Mac
```bash
go install github.com/tsliwowicz/go-wrk@latest
```

## 🔥 快速开始

### 1. 启动项目
```bash
cd bluebell_backend
go run main.go
```

### 2. 运行压测

#### 一键压测（推荐）
```bash
# Windows
.\tests\go_wrk_windows.bat

# Linux/Mac
./tests/go_wrk_test.sh
```

#### 单项测试
```bash
# 帖子列表查询
go-wrk -c 100 -d 30s -T 5s "http://localhost:8080/api/v1/posts?page=1&size=10"

# 帖子详情查看
go-wrk -c 100 -d 30s -T 5s "http://localhost:8080/api/v1/post/1"

# 热门排行
go-wrk -c 100 -d 30s -T 5s "http://localhost:8080/api/v1/posts/hot?days=7"
```

## 📊 压测结果解读

### 输出示例
```
Running 30s test @ http://localhost:8080/api/v1/posts
  100 goroutine(s) running concurrently
  3000 requests in 30.012s, 1.2MB read
Requests/sec:  99.96
Transfer/sec:  40.96KB
Avg Req Time:  1.000s
Fastest Request: 5.123ms
Slowest Request: 2.456s
Status Code 200: 3000 responses
```

### 关键指标
- **Requests/sec**: QPS（每秒请求数）
- **Avg Req Time**: 平均响应时间
- **Fastest/Slowest**: 最快/最慢响应时间
- **Status Code**: HTTP状态码分布

## 🎯 常用命令

### 渐进式压测
```bash
# 逐步增加并发
./tests/go_wrk_test.sh progressive
```

### 基准测试
```bash
# 性能基准测试
./tests/go_wrk_test.sh baseline
```

### 指定场景
```bash
# 只测试帖子列表
./tests/go_wrk_test.sh list

# 只测试投票
./tests/go_wrk_test.sh vote
```

## 🔧 参数调优

### 基础参数
| 参数 | 说明 | 示例 |
|------|------|------|
| `-c` | 并发数 | 100 |
| `-d` | 持续时间 | 30s |
| `-T` | 超时时间 | 5s |

### 高级参数
```bash
# 自定义Header
go-wrk -c 100 -d 30s -H "Authorization: Bearer token" "http://..."

# POST请求
go-wrk -c 50 -d 30s -M POST -B '{"key": "value"}' "http://..."

# 输出详细报告
./tests/go_wrk_test.sh full > tests/go_wrk_report.txt
```

## 📈 性能目标

### Bluebell论坛基准
| 场景 | 目标QPS | 响应时间 |
|------|---------|----------|
| 帖子列表 | >100 | <50ms |
| 帖子详情 | >200 | <30ms |
| 投票操作 | >50 | <100ms |
| 发帖操作 | >20 | <200ms |

## 🚨 常见问题

### 1. 连接超时
```bash
# 增加超时时间
go-wrk -c 100 -d 30s -T 10s "http://..."
```

### 2. 端口被占用
```bash
# 检查端口占用
netstat -ano | findstr :8080
```

### 3. Redis连接问题
```bash
# 检查Redis状态
redis-cli ping
```

## 🎉 一键运行

### Windows一键测试
```bash
cd tests
go_wrk_windows.bat
```

### Linux一键测试
```bash
cd tests
chmod +x go_wrk_test.sh
./go_wrk_test.sh
```

## 📱 监控建议

压测时同时监控：
- CPU使用率
- 内存占用
- Redis命中率
- MySQL连接数
- 网络带宽