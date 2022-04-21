package options

import cliflag "github.com/tiandh987/SharkAgent/pkg/cli/flag"

// options 包 包含初始化 apiserver 的 flags 和 options

// Options 运行一个 apiserver
type Options struct {

}

// NewOptions 使用默认参数创建一个 Options 对象
func NewOptions() *Options {
	o := Options{

	}

	return &o
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {

	return fss
}

// Validate checks Options and return a slice of found errs.
func (o *Options) Validate() []error {
	var errs []error

	return errs
}