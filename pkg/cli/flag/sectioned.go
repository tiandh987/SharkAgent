package flag

import (
	"bytes"
	"fmt"
	"github.com/spf13/pflag"
	"io"
	"strings"
)

// FlagSet 是一些预先定义好的 Flag 的集合，几乎所有的 Pflag 操作，
// 都需要借助 FlagSet 提供的方法来完成。
// 在实际开发中，我们可以使用两种方法来获取并使用 FlagSet：
//   方法一，调用 NewFlagSet 创建一个 FlagSet。
//   方法二，使用 Pflag 包定义的全局 FlagSet：CommandLine。实际上 CommandLine 也是由 NewFlagSet 函数创建的。

// NamedFlagSets 以调用 FlagSet() 的顺序存储 named flag sets
type NamedFlagSets struct {
	// Order 是 flag set names 的顺序列表
	Order []string
	// FlagSets 通过 name 存储 flag sets
	FlagSets map[string]*pflag.FlagSet
}

// FlagSet 根据给定的 name 返回 FlagSet（如果不存在则创建，并将 name 添加到 Order 中）
func (nfs *NamedFlagSets) FlagSet(name string) *pflag.FlagSet {
	// map 必须先初始化才能使用
	if nfs.FlagSets == nil {
		nfs.FlagSets = map[string]*pflag.FlagSet{}
	}

	if _, ok := nfs.FlagSets[name]; !ok {
		// 调用 NewFlagSet 创建一个 FlagSet。
		// NewFlagSet返回一个新的 空标志集，其指定名称、错误处理属性、SortFlags设置为true。
		nfs.FlagSets[name] = pflag.NewFlagSet(name, pflag.ExitOnError)
		nfs.Order = append(nfs.Order, name)
	}

	return nfs.FlagSets[name]
}

// PrintSections 在 sections 中打印 给定名称 的标志集，最大列数为给定列数。如果cols为零，则不包裹线条。
func PrintSections(w io.Writer, fss NamedFlagSets, cols int) {
	for _, name := range fss.Order {
		fs := fss.FlagSets[name]
		if !fs.HasFlags() {
			continue
		}

		wideFS := pflag.NewFlagSet("", pflag.ExitOnError)
		wideFS.AddFlagSet(fs)

		var zzz string
		if cols > 24 {
			zzz = strings.Repeat("z", cols-24)
			wideFS.Int(zzz, 0, strings.Repeat("z", cols-24))
		}

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "\n%s flags:\n\n%s", strings.ToUpper(name[:1])+name[1:], wideFS.FlagUsagesWrapped(cols))

		if cols > 24 {
			i := strings.Index(buf.String(), zzz)
			lines := strings.Split(buf.String()[:i], "\n")
			fmt.Fprint(w, strings.Join(lines[:len(lines)-1], "\n"))
			fmt.Fprintln(w)
		} else {
			fmt.Fprint(w, buf.String())
		}
	}
}