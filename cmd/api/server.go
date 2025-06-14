package api

import (
	"context"
	"ferry/database"
	"ferry/global/orm"
	"ferry/pkg/logger"
	"ferry/pkg/task"
	"ferry/router"
	"ferry/tools"
	config2 "ferry/tools/config"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config   string
	port     string
	mode     string
	StartCmd = &cobra.Command{
		Use:     "server",
		Short:   "Start API server",
		Example: "ferry server -c config/settings.yml",
		// PreRun 已被移除，所有逻辑都在 RunE 中
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&config, "config", "c", "config/settings.yml", "Start server with provided configuration file")
	StartCmd.PersistentFlags().StringVarP(&port, "port", "p", "", "Tcp port server listening on (overrides config file)")
	StartCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "", "server mode ; eg:dev,test,prod (overrides config file)")
}

func run() error {
	// --- 第一部分：健壮的初始化 ---
	fmt.Println("starting api server")

	// 1. 加载配置，并处理错误
	if err := config2.ConfigSetup(config); err != nil {
		// 在 logger 初始化前，只能用 fmt 打印致命错误
		fmt.Printf("FATAL: Failed to setup config: %v\n", err)
		return err
	}

	// 2. 初始化日志系统
	logger.Init()
	logger.Info("Configuration loaded and logger initialized successfully.")

	// 3. 初始化数据库链接
	database.Setup()
	logger.Info("Database connection setup successfully.")

	// 4. 启动异步任务队列
	go task.Start()
	logger.Info("Task worker started.")

	// --- 第二部分：服务器启动逻辑 ---

	// 应用命令行覆盖
	if mode != "" {
		viper.Set("settings.application.mode", mode)
	}
	if port != "" {
		viper.Set("settings.application.port", port)
	}

	if viper.GetString("settings.application.mode") == string(tools.ModeProd) {
		gin.SetMode(gin.ReleaseMode)
	}

	r := router.InitRouter()

	// 正确的 defer 块
	defer func() {
		// 根据编译器的指示，orm.Eloquent.DB() 只返回一个 *sql.DB 类型的值。
		sqlDB := orm.Eloquent.DB()
		if sqlDB == nil {
			logger.Error("Failed to get a valid *sql.DB instance for closing")
			return
		}
		// 关闭数据库连接池
		if err := sqlDB.Close(); err != nil {
			logger.Errorf("Failed to close database connection: %v", err)
		}
	}()

	srv := &http.Server{
		Addr:    viper.GetString("settings.application.host") + ":" + viper.GetString("settings.application.port"),
		Handler: r,
	}

	go func() {
		// 服务连接
		if viper.GetBool("settings.application.ishttps") {
			if err := srv.ListenAndServeTLS(viper.GetString("settings.ssl.pem"), viper.GetString("settings.ssl.key")); err != nil && err != http.ErrServerClosed {
				logger.Fatalf("listen: %s\n", err)
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Fatalf("listen: %s\n", err)
			}
		}
	}()

	fmt.Printf("%s Server Run http://%s:%s/ \r\n",
		tools.GetCurrntTimeStr(),
		viper.GetString("settings.application.host"),
		viper.GetString("settings.application.port"))
	fmt.Printf("%s Enter Control + C Shutdown Server \r\n", tools.GetCurrntTimeStr())

	// 优雅关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	fmt.Printf("%s Shutdown Server ... \r\n", tools.GetCurrntTimeStr())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown:", err)
	}
	logger.Info("Server exiting")
	return nil
}
