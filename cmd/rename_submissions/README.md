# 提交文件夹重命名脚本

这个脚本用于将已经提交的文件夹重命名为 "名字-赛道-作品名" 的格式，并同步更新数据库中的路径信息。

## 功能说明

脚本会：
1. 从数据库中查询所有提交的作品信息
2. 获取作者真实姓名、赛道名称、作品标题
3. 将文件夹重命名为 "姓名-赛道-作品名" 格式
4. 同步更新数据库中的文件路径信息
5. 处理文件名中的特殊字符，确保文件系统兼容性

## 使用方法

### 1. 进入脚本目录
```bash
cd cmd/rename_submissions
```

### 2. 配置数据库连接
编辑 `config.json` 文件，确保数据库连接信息正确：
```json
{
  "mysql": {
    "host": "127.0.0.1",
    "port": 3306,
    "username": "root",
    "password": "your_password",
    "database": "tank",
    "charset": "utf8"
  },
  "matterPath": "./matter",
  "dryRun": false,
  "maxNameLength": 50
}
```

### 3. 预览操作（推荐）
在执行实际重命名前，先运行预览模式：
```bash
# Windows
preview.bat

# Linux/Mac
# 手动设置 dryRun: true，然后运行：
go run main.go
```

### 4. 执行实际重命名
确认预览结果无误后，设置 `dryRun: false` 并运行：
```bash
# Windows
run.bat

# Linux/Mac
go mod tidy
go run main.go

# 或指定 matter 路径
go run main.go /path/to/matter
```

### 5. 编译版本
```bash
# 编译
go build -o rename_submissions main.go

# 运行
./rename_submissions

# 或指定路径
./rename_submissions /path/to/matter
```

## 配置说明

配置文件 `config.json` 参数说明：

```json
{
  "mysql": {
    "host": "127.0.0.1",           // 数据库主机地址
    "port": 3306,                  // 数据库端口
    "username": "root",            // 数据库用户名
    "password": "your_password",   // 数据库密码
    "database": "tank",            // 数据库名
    "charset": "utf8"              // 字符集
  },
  "matterPath": "./matter",        // matter 文件夹路径（相对或绝对路径）
  "dryRun": false,                 // 是否为预览模式（true: 仅预览，false: 执行操作）
  "maxNameLength": 50              // 文件名最大长度（避免路径过长）
}
```

### 重要参数说明：

- **dryRun**: 建议首次运行时设置为 `true`，预览操作结果
- **matterPath**: 可以是相对路径（如 `./matter`）或绝对路径（如 `D:\Project\TankSubmBox\matter`）
- **maxNameLength**: 控制生成的文件名长度，避免文件系统路径限制

## 安全注意事项

1. **备份数据**: 运行脚本前请务必备份数据库和文件夹
2. **测试环境**: 建议先在测试环境中验证脚本功能
3. **权限检查**: 确保脚本有足够的文件系统权限进行重命名操作
4. **数据库权限**: 确保数据库用户有 UPDATE 权限

## 处理逻辑

### 文件夹重命名规则
- 原格式：任意文件夹名
- 新格式：`姓名-赛道-作品名`
- 特殊字符处理：将 `/\:*?"<>|` 等字符替换为 `-`
- 长度限制：每个部分最长 50 个字符

### 数据库更新
1. 更新 `matter` 表中主文件夹的 `name` 字段
2. 更新所有子文件和子文件夹的 `path` 字段
3. 使用事务确保数据一致性

## 错误处理

脚本会处理以下情况：
- 原文件夹不存在
- 目标文件夹已存在
- 数据库连接失败
- 权限不足
- 路径过长

## 输出示例

```
开始重命名提交文件夹...
Matter 路径: ./matter
找到 3 个提交

处理第 1/3 个提交:
  原文件夹: matter/77/root/222-777-777
  新文件夹名: 张三-AI算法赛道-智能推荐系统
  作者: 张三, 赛道: AI算法赛道, 作品: 智能推荐系统
  ✅ 处理成功

处理第 2/3 个提交:
  原文件夹: matter/888/root/abc-def-ghi
  新文件夹名: 李四-Web开发赛道-在线学习平台
  作者: 李四, 赛道: Web开发赛道, 作品: 在线学习平台
  ✅ 处理成功

处理完成！成功: 2, 失败: 0
```

## 故障排除

1. **数据库连接失败**：检查数据库配置和网络连接
2. **权限错误**：确保脚本运行用户有文件系统写权限
3. **路径不存在**：检查 matter 路径是否正确
4. **中文乱码**：确保数据库连接字符集为 UTF-8