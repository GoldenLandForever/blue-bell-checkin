@echo off
REM Bluebell论坛go-wrk压测脚本（Windows版）
REM 使用方法: go_wrk_windows.bat [场景]

set BASE_URL=http://localhost:8080
set CONCURRENT=100
set DURATION=30s
set TIMEOUT=5s

echo [INFO] 🚀 Bluebell论坛go-wrk压测开始

REM 检查go-wrk是否安装
go-wrk --help >nul 2>&1
if %errorlevel% neq 0 (
    echo [INFO] 安装go-wrk...
    go install github.com/tsliwowicz/go-wrk@latest
)

REM 压测场景函数
goto :main

:test_post_list
echo [TEST] 📋 帖子列表查询
go-wrk -c %CONCURRENT% -d %DURATION% -T %TIMEOUT% ^
    -H "Accept: application/json" ^
    "%BASE_URL%/api/v1/posts?page=1&size=10"
goto :eof

:test_post_detail
echo [TEST] 📖 帖子详情查看
go-wrk -c %CONCURRENT% -d %DURATION% -T %TIMEOUT% ^
    -H "Accept: application/json" ^
    "%BASE_URL%/api/v1/post/1"
goto :eof

:test_hot_posts
echo [TEST] 🔥 热门帖子排行
go-wrk -c %CONCURRENT% -d %DURATION% -T %TIMEOUT% ^
    -H "Accept: application/json" ^
    "%BASE_URL%/api/v1/posts/hot?days=7"
goto :eof

:test_community_posts
echo [TEST] 🏘️ 社区帖子列表
go-wrk -c %CONCURRENT% -d %DURATION% -T %TIMEOUT% ^
    -H "Accept: application/json" ^
    "%BASE_URL%/api/v1/community/1/posts?page=1&size=10"
goto :eof

:test_vote_post
echo [TEST] 👍 投票操作
REM 创建测试数据
curl -s -X POST "%BASE_URL%/api/v1/vote" ^
    -H "Content-Type: application/json" ^
    -d "{\"post_id\": 1, \"direction\": 1}" >nul 2>&1

go-wrk -c 50 -d %DURATION% -T %TIMEOUT% ^
    -M POST ^
    -H "Content-Type: application/json" ^
    -B "{\"post_id\": 1, \"direction\": 1}" ^
    "%BASE_URL%/api/v1/vote"
goto :eof

:test_create_post
echo [TEST] ✍️ 创建帖子
go-wrk -c 20 -d %DURATION% -T %TIMEOUT% ^
    -M POST ^
    -H "Content-Type: application/json" ^
    -B "{\"title\": \"压测帖子\", \"content\": \"这是压测内容\", \"community_id\": 1}" ^
    "%BASE_URL%/api/v1/post"
goto :eof

:run_full_test
echo [INFO] 🎯 开始完整压测流程

REM 预热缓存
echo [INFO] 1. 预热缓存...
curl -s "%BASE_URL%/api/v1/posts?page=1&size=10" >nul 2>&1

echo [INFO] 2. 帖子列表压测
call :test_post_list

echo [INFO] 3. 帖子详情压测
call :test_post_detail

echo [INFO] 4. 热门排行压测
call :test_hot_posts

echo [INFO] 5. 社区帖子压测
call :test_community_posts

echo [INFO] 6. 投票操作压测
call :test_vote_post

echo [INFO] 7. 发帖操作压测
call :test_create_post

echo [INFO] ✅ 完整压测完成
goto :eof

:run_progressive_test
echo [INFO] 📈 渐进式压测（逐步增加并发）

for %%c in (10 50 100 200 500 1000) do (
    echo [INFO] 并发数: %%c
    set CONCURRENT=%%c
    call :test_post_list
    timeout /t 2 >nul
)
goto :eof

:main
if "%1"=="list" goto test_post_list
if "%1"=="detail" goto test_post_detail
if "%1"=="hot" goto test_hot_posts
if "%1"=="community" goto test_community_posts
if "%1"=="vote" goto test_vote_post
if "%1"=="post" goto test_create_post
if "%1"=="progressive" goto run_progressive_test
if "%1"=="full" goto run_full_test

REM 默认执行完整测试
goto run_full_test

echo [INFO] 🎉 压测脚本执行完成！