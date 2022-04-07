package globalflag

import (
	"fmt"
	"github.com/spf13/pflag"
)

// AddGlobalFlags 显式注册 flags 库（log、verflag 等）针对“flag”中的全局标志集注册的标志。
// 我们这样做是为了防止不需要的 flags 泄漏到组件的 flagset 中。
func AddGlobalFlags(fs *pflag.FlagSet, name string) {
	fs.BoolP("help", "h", false, fmt.Sprintf("help for %s", name))
}