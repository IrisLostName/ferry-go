package migrate

import (
	"ferry/database"
	"ferry/global/orm"
	"ferry/models/gorm"
	"ferry/pkg/logger"
	"ferry/pkg/task"
	config2 "ferry/tools/config"
	"fmt"
	"os" // 需要导入 os 包

	"github.com/spf13/cobra"
)

var (
	config   string
	mode     string
	StartCmd = &cobra.Command{
		Use:   "init",
		Short: "initialize the database",
		Run: func(cmd *cobra.Command, args []string) {
			if err := run(); err != nil {
				// run 函数出错时，确保程序以非零状态码退出
				os.Exit(1)
			}
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&config, "config", "c", "config/settings.yml", "Start server with provided configuration file")
	StartCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "dev", "server mode ; eg:dev,test,prod")
}

// 你的 run() 函数必须是这个结构！

func run() error {
	// --- 第一部分：健壮的初始化 ---
	fmt.Println("starting api server")

	// 1. 加载配置，并处理错误
	if err := config2.ConfigSetup(config); err != nil {
		fmt.Printf("FATAL: Failed to setup config: %v\n", err)
		return err
	}

	// 2. 初始化日志系统 (这是解决你当前问题的关键！)
	logger.Init()
	logger.Info("Configuration loaded and logger initialized successfully.")

	// 3. 初始化数据库链接
	database.Setup()
	logger.Info("Database connection setup successfully.")

	// 4. 启动异步任务队列
	go task.Start()
	logger.Info("Task worker started.")

	// --- 第二部分：服务器启动逻辑 ---
	// …… (省略)

	return nil
}

func migrateModel() error {
	if config2.DatabaseConfig.Dbtype == "mysql" {
		orm.Eloquent = orm.Eloquent.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
	}
	return gorm.AutoMigrate(orm.Eloquent)
}
