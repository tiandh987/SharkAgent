package app

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"runtime"
	"strings"
)

// Command 是 cli 应用的子命令结构体。
// 推荐使用 app.NewCommand() 创建一个命令
type Command struct {
	usage    string
	desc     string
	options  CliOptions
	commands []*Command
	runFunc  RunCommandFunc
}

// 1. func (c *Command) cobraCommand() *cobra.Command
// 2.

// RunCommandFunc 定义应用程序的命令的启动回调函数
type RunCommandFunc func(args []string) error

//
func (c *Command) cobraCommand() *cobra.Command {
	cmd := &cobra.Command {
		Use: c.usage,
		Short: c.desc,
	}

	cmd.SetOut(os.Stdout)
	cmd.Flags().SortFlags = false

	if len(c.commands) > 0 {
		for _, command := range c.commands {
			cmd.AddCommand(command.cobraCommand())
		}
	}

	if c.runFunc != nil {
		cmd.Run = c.runCommand
	}

	if c.options != nil {
		for _, f := range c.options.Flags().FlagSets {
			cmd.Flags().AddFlagSet(f)
		}
	}

	addHelpCommandFlag(c.usage, cmd.Flags())

	return cmd
}

func (c *Command) runCommand(cmd *cobra.Command, args []string) {
	if c.runFunc != nil {
		if err := c.runFunc(args); err != nil {
			fmt.Printf("%v %v\n", color.RedString("Error:"), err)
			os.Exit(1)
		}
	}
}

// FormatBaseName 基于给定的 basenema，在不同的操作系统下格式化为不同的可执行文件名
func FormatBaseName(basename string) string {
	if runtime.GOOS == "windows" {
		basename = strings.ToLower(basename)
		basename = strings.TrimSuffix(basename, ".exe")
	}

	return basename
}