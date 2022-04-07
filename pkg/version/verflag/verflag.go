package verflag

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/tiandh987/SharkAgent/pkg/version"
	"os"
	"strconv"
)

type versionValue int

const (
	VersionFalse versionValue = 0
	VersionTrue  versionValue = 1
	VersionRaw   versionValue = 2
)

const strRawVersion string = "raw"

func (v *versionValue) Set(s string) error {
	// raw
	if s == strRawVersion {
		*v = VersionRaw
		return nil
	}
	// ParseBool 返回字符串表示的布尔值。
	// 它接受 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False。
	// 任何其他值都会返回错误。
	boolVal, err := strconv.ParseBool(s)
	if boolVal {
		*v = VersionTrue
	} else {
		*v = VersionFalse
	}
	return err
}

func (v *versionValue) Get() interface{} {
	return v
}

func (v *versionValue) String() string {
	if *v == VersionRaw {
		return strRawVersion
	}
	return fmt.Sprintf("%v", bool(*v == VersionTrue))
}

func (v *versionValue) IsBoolFlag() bool {
	return true
}

// The type of the flag as required by the pflag.Value interface.
func (v *versionValue) Type() string {
	return "version"
}

// =======================================================

// VersionVar 使用指定的 name、usage 字符串定义一个 flag.
func VersionVar(p *versionValue, name string, value versionValue, usage string) {
	*p = value
	flag.Var(p, name, usage)
	// "--version" will be treated as "--version=true"
	flag.Lookup(name).NoOptDefVal = "true"
}

// Version wraps the VersionVar function.
func Version(name string, value versionValue, usage string) *versionValue {
	p := new(versionValue)
	VersionVar(p, name, value, usage)
	return p
}

// ======================================================================

const versionFlagName = "version"

var versionFlag = Version(versionFlagName, VersionFalse, "Print version information and quit.")

// AddFlags 在任意 FlagSets 上注册此包的 flag，以便它们指向与全局标志相同的值。
func AddFlags(fs *flag.FlagSet) {
	fs.AddFlag(flag.Lookup(versionFlagName))
}

// =====================================================================

// PrintAndExitIfRequested 将检查是否传递了 -version 标志，如果是，则打印版本并退出。
func PrintAndExitIfRequested() {
	if *versionFlag == VersionRaw {
		fmt.Printf("%#v\n", version.Get())
		os.Exit(0)
	} else if *versionFlag == VersionTrue {
		fmt.Printf("%s\n", version.Get())
		os.Exit(0)
	}
}