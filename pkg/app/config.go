package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tiandh987/SharkAgent/pkg/util/homedir"
	"os"
	"path/filepath"
	"strings"
)

const configFlagName = "config"

var cfgFile string

func init() {
	pflag.StringVarP(&cfgFile, "config", "c", cfgFile, "Read configuration from specified `FILE`, "+
		"support JSON, TOML, YAML, HCL, or Java properties formats.")
}

// addConfigFlag 将特定服务器的 flags 添加到指定的 FlagSet 对象。
func addConfigFlag(basename string, fs *pflag.FlagSet) {
	// 通过 init() 函数将 config 这个 flag 添加到 CommandLine 这个pfalg全局FlagSet
	fs.AddFlag(pflag.Lookup(configFlagName))

	// 支持环境变量
	viper.AutomaticEnv()
	// 设置环境变量前缀
	viper.SetEnvPrefix(strings.Replace(strings.ToUpper(basename), "-", "_", -1))
	// 重写 Env 键
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// 在每个 command 的执行方法(Run、RunE)被调用前 要执行的函数
	// 在命令执行前.读取配置文件
	cobra.OnInitialize(func() {
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		} else {
			viper.AddConfigPath(".")

			if names := strings.Split(basename, "-"); len(names) > 1 {
				viper.AddConfigPath(filepath.Join(homedir.HomeDir(), "." + names[0]))
				viper.AddConfigPath(filepath.Join("/etc", names[0]))
			}

			viper.SetConfigName(basename)
		}

		if err := viper.ReadInConfig(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: failed to read configuration file(%s): %v\n", cfgFile, err)
			os.Exit(1)
		}
	})

}
