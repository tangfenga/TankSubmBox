package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	Charset  string `json:"charset"`
}

type Config struct {
	MySQL         MySQLConfig `json:"mysql"`
	MatterPath    string      `json:"matterPath"`
	DryRun        bool        `json:"dryRun"`
	MaxNameLength int         `json:"maxNameLength"`
}

type SubmissionInfo struct {
	MatterUuid    string
	MatterName    string
	MatterPath    string
	SpaceName     string
	RealName      string
	TrackName     string
	Title         string
	OldFolderPath string
	NewFolderName string
}

func main() {
	// 确定配置文件路径
	configFile := "config.json"
	if _, err := os.Stat("config_preview.json"); err == nil {
		configFile = "config_preview.json"
	}

	// 加载配置文件
	config, err := loadConfig(configFile)
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		return
	}

	// 如果指定了命令行参数，则使用参数作为 matter 路径
	if len(os.Args) > 1 {
		config.MatterPath = os.Args[1]
	}

	fmt.Printf("开始重命名提交文件夹...\n")
	fmt.Printf("Matter 路径: %s\n", config.MatterPath)
	if config.DryRun {
		fmt.Printf("🔍 DRY RUN 模式：仅预览，不执行实际操作\n")
	}

	// 连接数据库
	db, err := connectDB(config)
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		return
	}
	defer db.Close()

	// 获取所有提交信息
	submissions, err := getSubmissions(db, config)
	if err != nil {
		fmt.Printf("获取提交信息失败: %v\n", err)
		return
	}

	fmt.Printf("找到 %d 个提交\n", len(submissions))

	// 处理每个提交
	var successCount, failCount int
	for i, submission := range submissions {
		fmt.Printf("\n处理第 %d/%d 个提交:\n", i+1, len(submissions))
		fmt.Printf("  原文件夹: %s\n", submission.OldFolderPath)
		fmt.Printf("  新文件夹名: %s\n", submission.NewFolderName)
		fmt.Printf("  作者: %s, 赛道: %s, 作品: %s\n",
			submission.RealName, submission.TrackName, submission.Title)

		if config.DryRun {
			fmt.Printf("  🔍 DRY RUN: 跳过实际操作\n")
			successCount++
			continue
		}

		// 重命名文件夹
		err := renameFolder(config.MatterPath, submission)
		if err != nil {
			fmt.Printf("  ❌ 重命名失败: %v\n", err)
			failCount++
			continue
		}

		// 更新数据库中的路径
		err = updateMatterPath(db, submission)
		if err != nil {
			fmt.Printf("  ❌ 更新数据库失败: %v\n", err)
			failCount++
			continue
		}

		fmt.Printf("  ✅ 处理成功\n")
		successCount++
	}

	fmt.Printf("\n处理完成！成功: %d, 失败: %d\n", successCount, failCount)
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 设置默认值
	if config.MaxNameLength <= 0 {
		config.MaxNameLength = 50
	}
	if config.MySQL.Charset == "" {
		config.MySQL.Charset = "utf8"
	}

	return &config, nil
}

func connectDB(config *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.MySQL.Username,
		config.MySQL.Password,
		config.MySQL.Host,
		config.MySQL.Port,
		config.MySQL.Database,
		config.MySQL.Charset,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getSubmissions(db *sql.DB, config *Config) ([]SubmissionInfo, error) {
	query := `
		SELECT 
			m.uuid as matter_uuid,
			m.name as matter_name,
			m.path as matter_path,
			m.space_name as space_name,
			COALESCE(up.real_name, u.username) as real_name,
			COALESCE(t.name, '未知赛道') as track_name,
			s.title as submission_title
		FROM tank41_matter m
		INNER JOIN submission s ON m.uuid = s.matter_uuid
		LEFT JOIN tank41_track t ON s.track_id = t.id
		LEFT JOIN tank41_user u ON m.user_uuid = u.uuid
		LEFT JOIN user_profile up ON u.uuid = up.user_uuid
		WHERE m.dir = 1 
		AND m.deleted = 0
		AND m.puuid = 'root'
		ORDER BY m.create_time DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []SubmissionInfo
	for rows.Next() {
		var s SubmissionInfo
		err := rows.Scan(
			&s.MatterUuid,
			&s.MatterName,
			&s.MatterPath,
			&s.SpaceName,
			&s.RealName,
			&s.TrackName,
			&s.Title,
		)
		if err != nil {
			return nil, err
		}

		// 清理字符串，移除特殊字符
		s.RealName = sanitizeFileName(s.RealName, config.MaxNameLength)
		s.TrackName = sanitizeFileName(s.TrackName, config.MaxNameLength)
		s.Title = sanitizeFileName(s.Title, config.MaxNameLength)

		// 构建旧文件夹路径
		s.OldFolderPath = filepath.Join("matter", s.SpaceName, "root", s.MatterName)

		// 构建新文件夹名：名字-赛道-作品名
		s.NewFolderName = fmt.Sprintf("%s-%s-%s", s.RealName, s.TrackName, s.Title)

		submissions = append(submissions, s)
	}

	return submissions, nil
}

func sanitizeFileName(name string, maxLength int) string {
	// 移除文件名中不允许的字符
	replacer := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "-",
		"?", "-",
		"\"", "-",
		"<", "-",
		">", "-",
		"|", "-",
		"\n", "",
		"\r", "",
		"\t", "",
	)

	sanitized := replacer.Replace(name)

	// 移除首尾空格
	sanitized = strings.TrimSpace(sanitized)

	// 如果为空，返回默认值
	if sanitized == "" {
		sanitized = "未知"
	}

	// 限制长度，避免路径过长
	if len(sanitized) > maxLength {
		sanitized = sanitized[:maxLength]
	}

	return sanitized
}

func renameFolder(matterPath string, submission SubmissionInfo) error {
	oldPath := filepath.Join(matterPath, submission.SpaceName, "root", submission.MatterName)
	newPath := filepath.Join(matterPath, submission.SpaceName, "root", submission.NewFolderName)

	// 检查原文件夹是否存在
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return fmt.Errorf("原文件夹不存在: %s", oldPath)
	}

	// 检查新文件夹是否已存在
	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("目标文件夹已存在: %s", newPath)
	}

	// 重命名文件夹
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return fmt.Errorf("重命名失败: %v", err)
	}

	return nil
}

func updateMatterPath(db *sql.DB, submission SubmissionInfo) error {
	// 开启事务
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 更新主文件夹的名称
	query := `UPDATE tank41_matter SET name = ? WHERE uuid = ?`
	_, err = tx.Exec(query, submission.NewFolderName, submission.MatterUuid)
	if err != nil {
		return err
	}

	// 更新所有子文件和子文件夹的路径
	// 找到所有以原路径开头的文件
	oldPathPrefix := "/" + submission.MatterName
	newPathPrefix := "/" + submission.NewFolderName

	query = `UPDATE tank41_matter SET path = REPLACE(path, ?, ?) WHERE space_name = ? AND path LIKE ?`
	_, err = tx.Exec(query, oldPathPrefix, newPathPrefix, submission.SpaceName, oldPathPrefix+"%")
	if err != nil {
		return err
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
