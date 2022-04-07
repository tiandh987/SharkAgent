package flag

import (
	goflag "flag"
	"github.com/spf13/pflag"
	"github.com/tiandh987/SharkAgent/pkg/log"
	"strings"
)

// WordSepNormalizeFunc 改变所有包含“_” 分隔符的 flag（将 “_” 替换为 “-”）
func WordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.ReplaceAll(name, "_", "-"))
	}
	return pflag.NormalizedName(name)
}

// InitFlags 规范化，解析，然后记录命令行 flags。
// 1. SetNormalizeFunc 允许您添加一个可以 "翻译Flag名称" 的函数。
// 添加到 FlagSet 的 flag 将被翻译，然后当尝试查找 flag 时也将被翻译的。
// 因此，可以创建一个名为 "getURL" 的标志，并将其翻译为 "geturl"。
// 然后，用户可以传递 "--getUrl"，它也可以翻译为 "geturl"，一切都会正常进行。
//
// 2. AddGoFlagSet(newSet *goflag.FlagSet) 添加给定的 flag.FlagSet 到 pflag.FlagSet
// 2.1 flag.CommandLine （给定的 flag.FlagSet）
func InitFlags(flags *pflag.FlagSet) {
	flags.SetNormalizeFunc(WordSepNormalizeFunc)
	flags.AddGoFlagSet(goflag.CommandLine)
}

// PrintFlags 在日志（Debug）中记录 FlagSet
func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		log.Debugf("FLAG: --%s=%q", flag.Name, flag.Value)
	})
}