

# 基于bluebell-plus重新开发

- [x] 加入积分系统 （Kafka） 

- [x] 重写评论系统

- [x] 添加redis预热

go 并发查

加入签到系统

- [x] 热门帖子redis查询

压测报告：

注册用户

```
#powershell
21..100 | ForEach-Object { 
    $user = "testuser_$_"
    $pass = "Password@$_"
    $email = $user+"@test.com"
    $ConfirmPassword = "Password@$_"
    $signup = Invoke-RestMethod -Uri "http://localhost:8081/api/v1/signup" -Method Post -Body (@{username=$user; password=$pass;email=$email;confirm_password=$ConfirmPassword} | ConvertTo-Json) -ContentType "application/json"
    $login = Invoke-RestMethod -Uri "http://localhost:8081/api/v1/login" -Method Post -Body (@{username=$user; password=$pass} | ConvertTo-Json) -ContentType "application/json"
    $login.access_token | Out-File -Append -FilePath "user_tokens.txt"
    Write-Host "用户 $_ 创建成功"
}
```

登录用户获取token

```
#他妈的，手动来好了

```









```shell
go-wrk -c 50 -d 30 -T 5000 http://localhost:8081/api/v1/posts
Running 30s test @ http://localhost:8081/api/v1/posts
  50 goroutine(s) running concurrently
59977 requests in 30.0754985s, 6.47MB read
Requests/sec:           1994.21
Transfer/sec:           220.26KB
Overall Requests/sec:   1989.03
Overall Transfer/sec:   219.69KB
Fastest Request:        0s
Avg Req Time:           25.072ms
Slowest Request:        855.327ms
Number of Errors:       0
10%:                    0s
50%:                    555µs
75%:                    614µs
99%:                    1.015ms
75%:                    614µs
99%:                    1.015ms
99.9%:                  1.051ms
99.9999%:               1.06ms
99.99999%:              1.06ms
stddev:                 66.508ms
```

```shell
go-wrk -c 100 -d 10 -T 5000 http://localhost:8081/api/v1/post/185445628082388993
Running 10s test @ http://localhost:8081/api/v1/post/1
  100 goroutine(s) running concurrently
21279 requests in 10.098893029s, 2.29MB read
Requests/sec:           2107.06
Transfer/sec:           232.57KB
Overall Requests/sec:   2083.20
Overall Transfer/sec:   229.94KB
Fastest Request:        0s
Avg Req Time:           47.459ms
Slowest Request:        355.535ms
Number of Errors:       0
10%:                    0s
50%:                    562µs
75%:                    658µs
99%:                    4.195ms
99.9%:                  4.909ms
99.9999%:               4.909ms
99.99999%:              4.909ms
stddev:                 80.803ms
```

```
go-wrk -c 50 -d 30 -T 5000 http://localhost:8081/api/v1/posts2
Running 30s test @ http://localhost:8081/api/v1/posts2
  50 goroutine(s) running concurrently
57655 requests in 30.105360312s, 6.22MB read
Requests/sec:           1915.11
Transfer/sec:           211.61KB
Overall Requests/sec:   1908.93
Overall Transfer/sec:   210.93KB
Fastest Request:        0s
Avg Req Time:           26.107ms
Slowest Request:        954.463ms
Number of Errors:       0
10%:                    0s
50%:                    567µs
75%:                    752µs
99%:                    4.369ms
99.9%:                  4.496ms
99.9999%:               4.505ms
99.99999%:              4.505ms
stddev:                 69.926ms

```







# bluebell-plus

1. 优化时间 2023年4月12日 至 2023年6月
2. 优化内容
- [ ] 对bluebell后端代码 进行全体重构优化
- [ ] 对bluebell前端代码 全部bug进行修复更新并进行全体升级优化
> 注：为了个人数据库安全，已删除全部配置文件。
> 4月份以后拉代码的新人后端会缺少配置，无法启动，需自行配置。
> 如需后端配置模版，关注公众号进微信交流群(群公告)获取。

- QQ群：
  - 3群：805360166（活跃 巨佬云集）
  - 2群：579480724（满）
  - 1群：1007576722（满）
- 微信群：关注公众号回复：交流群
- 公众号：Gopher毛
- 微信：18836288306
- qq：2557523039

**感谢技术支持：**[☆往事随風☆](https://github.com/china-521)

* [个人博客](https://wk-blog.vip)
* [CSDN](https://blog.csdn.net/m0_47214030?spm=1000.2115.3001.5343)

## 技能清单
1. 雪花算法
2. gin框架
2. zap日志库
3. Viper配置管理
4. swagger生成文档
5. JWT认证
6. 令牌桶限流
7. Go语言操作MySQL **(sqlx)**
8. Go语言操作Redis **(go-redis)**
10. Gihub热榜
12. Docker部署
13. Vue框架
14. ElementUI
15. axios 
16. 畅言云评论系统

## 项目目录结构
### 后端结构树
```bash
.
├── Dockerfile
├── Makefile
├── README.md
├── bin
│   ├── bluebell-plus
│   └── bluebell-plus.conf
├── conf
│   └── config.yaml
├── controller
│   ├── code.go
│   ├── comment.go
│   ├── community.go
│   ├── doc_response_models.go
│   ├── post.go
│   ├── post_test.go
│   ├── request.go
│   ├── response.go
│   ├── user.go
│   ├── validator.go
│   └── vote.go
├── dao
│   ├── mysql
│   │   ├── comment.go
│   │   ├── community.go
│   │   ├── error_code.go
│   │   ├── mysql.go
│   │   ├── post.go
│   │   ├── post_test.go
│   │   └── user.go
│   └── redis
│       ├── error.go
│       ├── keys.go
│       ├── post.go
│       ├── redis.conf
│       ├── redis.go
│       └── vote.go
├── docker-compose.yml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── init.sql
├── log
│   └── bluebell-plus.log
├── logger
│   └── logger.go
├── logic
│   ├── community.go
│   ├── post.go
│   ├── truncate.go
│   ├── user.go
│   └── vote.go
├── main.go
├── middlewares
│   ├── auth.go
│   └── ratelimit.go
├── models
│   ├── comment.go
│   ├── community.go
│   ├── create_tables.sql
│   ├── params.go
│   ├── post.go
│   └── user.go
├── pkg
│   ├── jwt
│   │   └── jwt.go
│   └── snowflake
│       └── gen_id.go
├── routers
│   └── routers.go
├── settings
│   └── settings.go
├── static
│   ├── css
│   │   ├── app.5c39da08.css
│   │   └── chunk-vendors.5b539fe5.css
│   ├── favicon.ico
│   ├── fonts
│   │   ├── element-icons.535877f5.woff
│   │   ├── element-icons.732389de.ttf
│   │   ├── fontello.068ca2b3.ttf
│   │   ├── fontello.8d4a4e6f.woff2
│   │   ├── fontello.a782baa8.woff
│   │   └── fontello.e73a0647.eot
│   ├── img
│   │   ├── avatar.7b0a9835.png
│   │   ├── fontello.9354499c.svg
│   │   ├── iconfont.cdbe38a0.svg
│   │   ├── logo.938d1d61.png
│   │   └── search.8e85063d.png
│   └── js
│       ├── app.81e7c3d0.js
│       ├── app.81e7c3d0.js.map
│       ├── chunk-vendors.218b058e.js
│       └── chunk-vendors.218b058e.js.map
├── templates
│   └── index.html
├── version.go
└── wait-for.sh
```
### 前端结构树
```bash
├── bin
│   └── bluebell-plus
├── conf
│   └── config.yaml
├── static
│   ├── css
│   │   └── app.0afe9dae.css
│   ├── favicon.ico
│   ├── img
│   │   ├── avatar.7b0a9835.png
│   │   ├── iconfont.cdbe38a0.svg
│   │   ├── logo.da56125f.png
│   │   └── search.8e85063d.png
│   └── js
│       ├── app.9f3efa6d.js
│       ├── app.9f3efa6d.js.map
│       ├── chunk-vendors.57f9e9d6.js
│       └── chunk-vendors.57f9e9d6.js.map
└── templates
    └── index.html
```

## 项目预览图

[//]: # "[![bCORoR.png]&#40;https://s4.ax1x.com/2022/02/23/bCORoR.png&#41;]&#40;https://imgtu.com/i/bCORoR&#41;"

[![image.png](https://i.postimg.cc/brRyjhPL/image.png)](https://postimg.cc/zHVZn4VR)

## 项目全套笔记

- **视频教程地址**：[GoWeb进阶—两周开发一个基于vue+go+gin+mysql+redis的博客论坛web项目！！！从零到部署上线](https://www.bilibili.com/video/BV1Fb4y14747?spm_id_from=333.999.0.0)
- **GitHub仓库**：https://github.com/mao888/bluebell
- **GitEE仓库**：https://gitee.com/hu_maomao/bluebell
- 编程：用代码解决生活中的问题
- **技术与知识的区别**：
- - 知识：记住地球是圆的
  - 技术：自己学会游泳，自己学会开车 
- **基于雪花算法生成用户ID**
- - https://www.yuque.com/docs/share/e50bbca1-e019-45e2-b77b-a9ba01fbede3?# 《基于雪花算法生成用户ID》
- [gin框架中使用validator若干实用技巧](https://www.liwenzhou.com/posts/Go/validator_usages/)
- [《限制账号同一时间只能登录一个设备》](https://www.yuque.com/docs/share/584ddd0f-5158-4cea-8918-a4b6e1d41a07?# )
- [《基于Cookie、Session和基于Token的认证模式介绍》](https://www.yuque.com/docs/share/06a89a55-3e3c-452b-aeb1-acf4d2bac8a5?#)
- [在gin框架中使用JWT认证](https://www.liwenzhou.com/posts/Go/jwt_in_gin/)
- [为Go项目编写Makefile](https://www.liwenzhou.com/posts/Go/makefile/)
- [使用Air实现Go程序实时热重载](https://www.liwenzhou.com/posts/Go/live_reload_with_air/)
- [分页](https://zhidao.baidu.com/question/1573826651037645420.html)
- [JSON实战拾遗之数字精度](https://www.ituring.com.cn/article/506822)
- [你需要知道的那些go语言json技巧](https://www.liwenzhou.com/posts/Go/json_tricks_in_go)
- [帖子投票（点赞）功能设计与实现](https://www.yuque.com/docs/share/d09afe84-90d1-4e04-a73e-95848f073558?#)
- [《基于用户投票的排名算法》](https://www.yuque.com/docs/share/f40f5c41-f327-47d4-88bb-02bcf62515a8?# )
- [使用swagger生成接口文档](https://www.liwenzhou.com/posts/Go/gin_swagger/)
- [HTTP Server常用压测工具介绍](https://www.liwenzhou.com/posts/Go/benchmark_tool/)
- [漏桶和令牌桶限流策略介绍及使用](https://www.liwenzhou.com/posts/Go/ratelimit/)
- [option选项模式](https://www.liwenzhou.com/posts/Go/functional_options_pattern/)
- [Go pprof性能调优](https://www.liwenzhou.com/posts/Go/performance_optimisation/)
- [如何使用docker部署Go Web程序](https://www.liwenzhou.com/posts/Go/how_to_deploy_go_app_using_docker/)
- [部署Go语言程序的N种方法](https://www.liwenzhou.com/posts/Go/deploy_go_app/)
- [《企业代码发布流程及CICD介绍》](https://www.yuque.com/docs/share/e837e5bf-f6a9-4dc8-98e4-4b8ce24808ab?)