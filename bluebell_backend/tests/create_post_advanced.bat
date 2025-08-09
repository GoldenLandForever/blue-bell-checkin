@echo off
setlocal enabledelayedexpansion

REM 创建帖子压测高级脚本
REM 支持多种测试模式和参数配置

echo.
echo =================================
echo 🚀 Bluebell 创建帖子压测工具
echo =================================
echo.

REM 默认参数
set CONCURRENT=50
set DURATION=30
set BASE_URL=http://localhost:8081
set MODE=manual

REM 解析命令行参数
:parse_args
if "%~1"=="" goto :validate_args
if "%~1"=="-c" (
    set CONCURRENT=%~2
    shift
    shift
    goto :parse_args
)
if "%~1"=="-d" (
    set DURATION=%~2
    shift
    shift
    goto :parse_args
)
if "%~1"=="-u" (
    set BASE_URL=%~2
    shift
    shift
    goto :parse_args
)
if "%~1"=="-m" (
    set MODE=%~2
    shift
    shift
    goto :parse_args
)
if "%~1"=="-h" goto :show_help
shift
goto :parse_args

:validate_args
REM 验证参数
if %CONCURRENT% LSS 1 (
    echo ❌ 并发数必须大于0
    exit /b 1
)
if %DURATION% LSS 1 (
    echo ❌ 持续时间必须大于0秒
    exit /b 1
)

REM 显示配置
echo 📊 当前配置:
echo   并发数: %CONCURRENT%
echo   持续时间: %DURATION%秒
echo   目标URL: %BASE_URL%
echo   测试模式: %MODE%
echo.

REM 检查服务状态
echo 🔍 检查服务状态...
curl -s -o nul -w "%%{http_code}" %BASE_URL%/api/v1/post 2>nul > temp_status.txt
set /p STATUS=<temp_status.txt
del temp_status.txt

if "%STATUS%"=="405" (
    echo ✅ 服务正常运行
) else (
    echo ⚠️  服务可能未启动或端口错误
    echo 📝 请确保服务已启动: go run main.go
    pause
    exit /b 1
)

REM 根据模式执行测试
echo.
echo 🎯 开始压测...
echo.

if "%MODE%"=="manual" (
    goto :run_manual
) else if "%MODE%"=="go-wrk" (
    goto :run_go_wrk
) else if "%MODE%"=="compare" (
    goto :run_compare
) else (
    echo ❌ 无效的测试模式: %MODE%
    exit /b 1
)

:run_manual
echo 📈 使用Go程序进行压测...
if exist create_post_manual.go (
    echo 🏗️  编译压测程序...
    go build -o create_post_manual.exe create_post_manual.go
    if !errorlevel! neq 0 (
        echo ❌ 编译失败
        exit /b 1
    )
    
    echo 🚀 执行压测...
    create_post_manual.exe -c %CONCURRENT% -d %DURATION% -u %BASE_URL%
) else (
    echo ❌ 找不到 create_post_manual.go
    exit /b 1
)
goto :finish

:run_go_wrk
echo 📈 使用go-wrk进行压测...
echo 📝 准备测试数据...

REM 创建测试数据文件
set POST_DATA={"title":"测试帖子","content":"这是用于压测的测试内容","community_id":1}
set TEMP_FILE=temp_post_data.txt
echo %POST_DATA% > %TEMP_FILE%

echo 🚀 执行go-wrk压测...
go-wrk -c %CONCURRENT% -d %DURATION% -T 5000 -body "%POST_DATA%" -method POST -H "Content-Type: application/json" %BASE_URL%/api/v1/post

del %TEMP_FILE%
goto :finish

:run_compare
echo 📊 运行对比测试...
echo.
echo 🎯 测试1: Go程序压测
echo ==================
if exist create_post_manual.go (
    go build -o create_post_manual.exe create_post_manual.go
    create_post_manual.exe -c %CONCURRENT% -d 10 -u %BASE_URL%
)

echo.
echo 🎯 测试2: go-wrk压测
echo ==================
set POST_DATA={"title":"对比测试","content":"这是对比测试内容","community_id":1}
go-wrk -c %CONCURRENT% -d 10 -T 5000 -body "%POST_DATA%" -method POST -H "Content-Type: application/json" %BASE_URL%/api/v1/post

goto :finish

:show_help
echo.
echo 📖 使用帮助:
echo   %0 [选项]
echo.
echo 选项:
echo   -c ^<num^>     并发数 (默认: 50)
echo   -d ^<sec^>     持续时间秒数 (默认: 30)
echo   -u ^<url^>     目标URL (默认: http://localhost:8081)
echo   -m ^<mode^>    测试模式: manual/go-wrk/compare (默认: manual)
echo   -h            显示帮助

echo.
echo 示例:
echo   %0 -c 100 -d 60 -m manual
echo   %0 -c 50 -d 30 -m go-wrk
echo   %0 -c 20 -d 10 -m compare
exit /b 0

:finish
echo.
echo ✅ 压测完成!
echo.
echo 💡 性能分析建议:
echo   - 如果QPS低于100，检查数据库索引
if %CONCURRENT% GTR 50 (
echo   - 高并发下注意连接池配置
)
echo   - 查看数据库慢查询日志
echo   - 考虑使用Redis缓存

pause