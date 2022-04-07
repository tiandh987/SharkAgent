package app

import (
	"fmt"
	"github.com/spf13/cobra"
)

// App 是 cli 应用程序的主要结构体
// 推荐使用 app.NewApp() 创建一个 app
type App struct {
	basename    string                 // 可执行文件名               ==> 对应 cobra.Command.Use
	name        string                 // 应用程序名                 ==> 对应 cobra.Command.Short
	description string                 // 应用程序描述信息           ==> 对应 cobra.Command.Long
	options     CliOptions             // 从命令行读取的配置参数（设置FlagSet）
	runFunc     RunFunc                // 应用程序的启动回调函数      ==> 对应 cobra.Command.RunE
	silence     bool                   // 静默模式（不输出 启动/配置/版本 信息到控制台）
	noVersion   bool                   // 是否提供 version flag
	noConfig    bool                   // 是否提供 config flag
	commands    []*Command             // 应用程序的子命令
	args        cobra.PositionalArgs                             // ==> 对应 cobra.Command.Args
	cmd         *cobra.Command
}

// RunFunc 定义应用程序的启动回调函数
type RunFunc func(basename string) error

// Option 定义用于初始化应用程序结构的可选参数。
type Option func(*App)

// 设计模式：选项模式
// 好处：在扩展新属性的时候，不用修改 NewApp()方法;
//       动态配置 APP
// 1. func WithDescription(desc string) Option
// 2. func WithOptions(opt CliOptions) Option
// 3. func WithSilence() Option
// 4. func WithNoVersion() Option
// 5. func WithNoConfig() Option
// 6. func WithValidArgs(args cobra.PositionalArgs) Option
// 7. func WithDefaultValidArgs() Option
// 8. func WithRunFunc(run RunFunc) Option

// WithDescription 用于设置应用程序的描述信息
func WithDescription(desc string) Option {
	return func(a *App) {
		a.description = desc
	}
}

// WithOptions 打开应用程序的函数，从命令行读取或从配置文件读取参数。
func WithOptions(opt CliOptions) Option {
	return func(a *App) {
		a.options = opt
	}
}

// WithSilence 设置应用程序为 silent 模式
// 项目的启动信息、配置信息、版本信息都不会打印到控制台
func WithSilence() Option {
	return func(a *App) {
		a.silence = true
	}
}

// WithNoVersion 设置应用程序不提供 version flag
func WithNoVersion() Option {
	return func(a *App) {
		a.noVersion = true
	}
}

// WithNoConfig 设置应用程序不提供 config flag
func WithNoConfig() Option {
	return func(a *App) {
		a.noConfig = true
	}
}

// WithValidArgs 设置验证函数，验证 non-flag 参数
func WithValidArgs(args cobra.PositionalArgs) Option {
	return func(a *App) {
		a.args = args
	}
}

// WithDefaultValidArgs 设置默认的验证函数去验证 non-flag 参数
// 校验命令行非选项参数的默认校验逻辑
func WithDefaultValidArgs() Option {
	return func(a *App) {
		a.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		}
	}
}

// WithRunFunc 用于设置应用程序的启动回调函数
func WithRunFunc(run RunFunc) Option {
	return func(a *App) {
		a.runFunc = run
	}
}