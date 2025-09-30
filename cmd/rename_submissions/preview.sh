#!/bin/bash

echo "预览提交文件夹重命名操作"
echo "============================"

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
echo "创建临时配置文件（DRY RUN 模式）..."

# 创建临时的预览配置文件
cat > config_preview.json << EOF
{
  "mysql": {
    "host": "127.0.0.1",
    "port": 3306,
    "username": "root",
    "password": "q1472580369.",
    "database": "tank",
    "charset": "utf8"
  },
  "matterPath": "./matter",
  "dryRun": true,
  "maxNameLength": 50
}
EOF

echo ""
echo "开始预览重命名操作..."
echo ""

go run main.go

if [ $? -ne 0 ]; then
    echo ""
    echo "预览执行失败"
    rm -f config_preview.json
    exit 1
else
    echo ""
    echo "预览完成！如果结果正确，请将 config.json 中的 dryRun 设置为 false 后执行实际操作"
fi

# 清理临时文件
rm -f config_preview.json