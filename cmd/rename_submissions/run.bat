@echo off
echo 提交文件夹重命名脚本
echo ========================

REM 检查是否在正确的目录
if not exist "main.go" (
    echo 错误：请在 cmd/rename_submissions 目录下运行此脚本
    pause
    exit /b 1
)

echo 正在安装依赖...
go mod tidy
if %errorlevel% neq 0 (
    echo 依赖安装失败
    pause
    exit /b 1
)

echo.
echo 开始运行重命名脚本...
echo.

REM 如果有参数传入，使用参数作为 matter 路径
if "%1"=="" (
    go run main.go
) else (
    go run main.go "%1"
)

if %errorlevel% neq 0 (
    echo.
    echo 脚本执行失败
) else (
    echo.
    echo 脚本执行完成
)

echo.
pause