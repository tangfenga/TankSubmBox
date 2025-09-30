#!/bin/bash

echo "提交文件夹重命名脚本"
echo "========================"

# 检查是否在正确的目录
if [ ! -f "main.go" ]; then
    echo "错误：请在 cmd/rename_submissions 目录下运行此脚本"
    exit 1
fi

echo "正在安装依赖..."
go mod tidy
if [ $? -ne 0 ]; then
    echo "依赖安装失败"
    exit 1
fi

echo ""
echo "开始运行重命名脚本..."
echo ""

# 如果有参数传入，使用参数作为 matter 路径
if [ -z "$1" ]; then
    go run main.go
else
    go run main.go "$1"
fi

if [ $? -ne 0 ]; then
    echo ""
    echo "脚本执行失败"
    exit 1
else
    echo ""
    echo "脚本执行完成"
fi