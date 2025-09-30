@echo off
echo 预览提交文件夹重命名操作
echo ============================

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
echo 创建临时配置文件（DRY RUN 模式）...

REM 创建临时的预览配置文件
echo {> config_preview.json
echo   "mysql": {>> config_preview.json
echo     "host": "127.0.0.1",>> config_preview.json
echo     "port": 3306,>> config_preview.json
echo     "username": "root",>> config_preview.json
echo     "password": "q1472580369.",>> config_preview.json
echo     "database": "tank",>> config_preview.json
echo     "charset": "utf8">> config_preview.json
echo   },>> config_preview.json
echo   "matterPath": "./matter",>> config_preview.json
echo   "dryRun": true,>> config_preview.json
echo   "maxNameLength": 50>> config_preview.json
echo }>> config_preview.json

echo.
echo 开始预览重命名操作...
echo.

go run main.go

if %errorlevel% neq 0 (
    echo.
    echo 预览执行失败
) else (
    echo.
    echo 预览完成！如果结果正确，请将 config.json 中的 dryRun 设置为 false 后执行实际操作
)

REM 清理临时文件
del config_preview.json

echo.
pause