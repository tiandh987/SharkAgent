package app

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/marmotedu/errors"
	cliflag "github.com/tiandh987/SharkAgent/pkg/cli/flag"
	"github.com/tiandh987/SharkAgent/pkg/cli/globalflag"
	"github.com/tiandh987/SharkAgent/pkg/log"
	"github.com/tiandh987/SharkAgent/pkg/term"
	"github.com/tiandh987/SharkAgent/pkg/version"
	"github.com/tiandh987/SharkAgent/pkg/version/verflag"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	progressMessage = color.GreenString("==>")

	usageTemplate = fmt.Sprintf(`%s{{if .Runnable}}
  %s{{end}}{{if .HasAvailableSubCommands}}
  %s{{end}}{{if gt (len .Aliases) 0}}

%s
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

%s
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

%s{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  %s {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

%s
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

%s
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

%s{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "%s --help" for more information about a command.{{end}}
`,
		color.CyanString("Usage:"),
		color.GreenString("{{.UseLine}}"),
		color.GreenString("{{.CommandPath}} [command]"),
		color.CyanString("Aliases:"),
		color.CyanString("Examples:"),
		color.CyanString("Available Commands:"),
		color.GreenString("{{rpad .Name .NamePadding }}"),
		color.CyanString("Flags:"),
		color.CyanString("Global Flags:"),
		color.CyanString("Additional help topics:"),
		color.GreenString("{{.CommandPath}} [command]"),
	)
)

// App 是 cli 应用的主要结构体
// 推荐使用 app.NewApp() 函数创建一个 app
// 创建一个应用
type App struct {
	basename    string                // 可执行文件名                 ==> 对应 cobra.Command.Use
	name        string                // 应用程序名                   ==> 对应 cobra.Command.Short
	description string                // 应用程序描述信息             ==> 对应 cobra.Command.Long
	options     CliOptions            // 从命令行读取的配置参数（设置FlagSet）
	runFunc     RunFunc               // 应用程序的启动回调函数        ==> 对应 cobra.Command.RunE
	silence     bool                  // 静默模式（不输出 启动/配置/版本 信息到控制台）
	noVersion   bool                  // 是否提供 version flag
	noConfig    bool                  // 是否提供 config flag
	commands    []*Command            // 应用程序的子命令
	args        cobra.PositionalArgs                             // ==> 对应 cobra.Command.Args
	cmd         *cobra.Command
}

// 1. func (a *App) buildCommand()
// 2. func (a *App) runCommand(cmd *cobra.Command, args []string) error
// 3. func (a *App) applyOptionRules() error
// 4. func (a *App) Run()
// 5. func (a *App) Command()

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

// ===============================

// NewApp 基于给定的 应用名称、二进制名称、其他选项 创建一个应用程序实例
// Options 个人理解：我有一个参数值（示例：description）,我需要一个App示例来进行赋值
func NewApp(name string, basename string, opts ...Option) *App {
	a := &App{
		name:     name,
		basename: basename,
	}

	for _, o := range opts {
		o(a)
	}

	a.buildCommand()

	return a
}

// 创建 Cobra Command 类型的命令,命令的功能通过指定 Cobra Command 类型的各个
// 字段来实现。
func (a *App) buildCommand() {
	cmd := cobra.Command{
		Use:   FormatBaseName(a.basename),
		Short: a.name,
		Long:  a.description,
		// stop printing usage when the command errors
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          a.args,
	}

	// 设置标准输出
	cmd.SetOut(os.Stdout)
	// 设置标准错误
	cmd.SetErr(os.Stderr)
	// Flags() 返回适用于此命令的完整 FlagSet（在此处和所有父级声明的 local 和 persistent）。
	// SortFlags 用于指示用户是否希望在 帮助/使用消息中具有 排序 标志。
	cmd.Flags().SortFlags = true

	// InitFlags()
	// 1. SetNormalizeFunc 允许您添加一个可以 "翻译Flag名称" 的函数。
	// 添加到 FlagSet 的标志将被翻译，然后当尝试查找标志时也将被翻译的。
	// 因此，可以创建一个名为 "getURL" 的标志，并将其翻译为 "geturl"。
	// 然后，用户可以传递 "--getUrl"，它也可以翻译为 "geturl"，一切都会正常进行。
	// 1.1 WordSepNormalizeFunc() 更改所有包含“_”分隔符的标志。
	//     将 “_” 替换为 “-”
	//
	// 2. AddGoFlagSet(newSet *goflag.FlagSet) 添加给定的 flag.FlagSet 到 pflag.FlagSet
	// 2.1 flag.CommandLine （给定的 flag.FlagSet）
	//
	cliflag.InitFlags(cmd.Flags())

	// 附加命令
	if len(a.commands) > 0 {
		for _, command := range a.commands {
			cmd.AddCommand(command.cobraCommand())
		}
		cmd.SetHelpCommand(helpCommand(FormatBaseName(a.basename)))
	}

	// 设置应用程序运行回调函数
	if a.runFunc != nil {
		cmd.RunE = a.runCommand
	}

	// 从命令行读取的参数
	// NamedFlagSets 按调用 FlagSet 的顺序存储命名标志集。
	var namedFlagSets cliflag.NamedFlagSets

	// 将应用程序（eg：api-server）提供的 FlagSet 添加到 cobra.command.flags
	if a.options != nil {
		namedFlagSets = a.options.Flags()
		// cmd.Flags() 返回 cobra command 的 FlagSet
		fs := cmd.Flags()
		for _, f := range namedFlagSets.FlagSets {
			// AddFlagSet 将一个标记集添加到另一个标记集。如果 f 中已经存在一个标志，那么来自newSet的标志将被忽略。
			fs.AddFlagSet(f)
		}
	}

	// --version 默认值为 true
	if !a.noVersion {
		verflag.AddFlags(namedFlagSets.FlagSet("global"))
	}

	// --config
	if !a.noConfig {
		addConfigFlag(a.basename, namedFlagSets.FlagSet("global"))
	}

	// 添加全局 FlagSet global
	// cmd.Name() 返回应command的名称（use line （Use 字段） 的第一个单词）
	// --help 添加到 global flagSet
	globalflag.AddGlobalFlags(namedFlagSets.FlagSet("global"), cmd.Name())

	// add new global flagset to cmd FlagSet
	// 将 global FlagSet 添加到 cobra.Command.flags
	cmd.Flags().AddFlagSet(namedFlagSets.FlagSet("global"))

	// 自定义 Usage、Help
	addCmdTemplate(&cmd, namedFlagSets)
	a.cmd = &cmd
}

func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	// 打印工作目录
	printWorkingDir()
	// 打印 Options
	cliflag.PrintFlags(cmd.Flags())

	if !a.noVersion {
		// display application version information
		verflag.PrintAndExitIfRequested()
	}

	// 将配置文件中的配置项和命令行参数绑定
	if !a.noConfig {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}

		if err := viper.Unmarshal(a.options); err != nil {
			return err
		}
	}

	if !a.silence {
		log.Infof("%v Starting %s ...", progressMessage, a.name)
		if !a.noVersion {
			log.Infof("%v Version: `%s`", progressMessage, version.Get().ToJSON())
		}
		if !a.noConfig {
			log.Infof("%v Config file used: `%s`", progressMessage, viper.ConfigFileUsed())
		}
	}

	if a.options != nil {
		// 完整填充、检查、日志记录 Optioins
		if err := a.applyOptionRules(); err != nil {
			return err
		}
	}

	// 启动应用程序
	if a.runFunc != nil {
		return a.runFunc(a.basename)
	}

	return nil
}

// 1. 完整填充 Options
// 2. 检查 Options
// 3. 记录 Options 到日志中
func (a *App) applyOptionRules() error {
	if completeableOptions, ok := a.options.(CompleteableOptions); ok {
		// 完善 Options，完整填写所有未设置的字段（这些字段必须具有有效数据）。
		if err := completeableOptions.Complete(); err != nil {
			return err
		}
	}

	// 检查各类options，并返回一个errors slice
	if errs := a.options.Validate(); len(errs) != 0 {
		return errors.NewAggregate(errs)
	}

	// 记录 Options 到日志中
	if printableOptions, ok := a.options.(PrintableOptions); ok && !a.silence {
		log.Infof("%v Config: `%s`", progressMessage, printableOptions.String())
	}

	return nil
}

// Run 用于运行应用程序
func (a *App) Run() {
	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}

// Command 返回应用程序内的 cobra 命令实例。
func (a *App) Command() *cobra.Command {
	return a.cmd
}

// 打印工作目录
func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Infof("%v WorkingDir: %s", progressMessage, wd)
}

// 自定义 Usage、Help
func addCmdTemplate(cmd *cobra.Command, namedFlagSets cliflag.NamedFlagSets) {
	usageFmt := "Usage:\n  %s\n"

	// 获取终端宽度、高度
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())

	// 自定义 Usage
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)

		return nil
	})

	// 自定义 Help
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})
}
