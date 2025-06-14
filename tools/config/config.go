// 文件: C:/go-project/ferry/tools/config/config.go

package config

import (
	// "ferry/pkg/logger" // 不再需要导入 logger
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// 全局变量定义保持不变
var cfgDatabase *viper.Viper
var cfgApplication *viper.Viper
var cfgJwt *viper.Viper
var cfgSsl *viper.Viper

// ConfigSetup 载入配置文件，成功返回 nil，失败返回 error
func ConfigSetup(path string) error { // <--- 修改 1：增加 "error" 返回值
	viper.SetConfigFile(path)
	content, err := os.ReadFile(path)
	if err != nil {
		// 修改 2：返回一个包装过的错误
		return fmt.Errorf("read config file fail: %w", err)
	}

	// Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		// 修改 3：返回错误
		return fmt.Errorf("parse config file fail: %w", err)
	}

	// 将所有的 panic(……) 调用替换为返回 error
	cfgDatabase = viper.Sub("settings.database")
	if cfgDatabase == nil {
		return fmt.Errorf("config not found: settings.database")
	}
	DatabaseConfig = InitDatabase(cfgDatabase)

	cfgApplication = viper.Sub("settings.application")
	if cfgApplication == nil {
		return fmt.Errorf("config not found: settings.application")
	}
	ApplicationConfig = InitApplication(cfgApplication)

	cfgJwt = viper.Sub("settings.jwt")
	if cfgJwt == nil {
		return fmt.Errorf("config not found: settings.jwt")
	}
	JwtConfig = InitJwt(cfgJwt)

	cfgSsl = viper.Sub("settings.ssl")
	if cfgSsl == nil {
		return fmt.Errorf("config not found: settings.ssl")
	}
	SslConfig = InitSsl(cfgSsl)

	// logger.Init() // <--- 修改 4：确保这一行被删除了

	return nil // <--- 修改 5：如果一切顺利，返回 nil
}

// SetConfig 函数保持不变
func SetConfig(configPath string, key string, value interface{}) {
	viper.AddConfigPath(configPath)
	viper.Set(key, value)
	_ = viper.WriteConfig()
}
