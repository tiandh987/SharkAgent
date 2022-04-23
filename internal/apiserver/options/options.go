package options

import (
	genericoptions "github.com/tiandh987/SharkAgent/internal/pkg/options"
	cliflag "github.com/tiandh987/SharkAgent/pkg/cli/flag"
	"github.com/tiandh987/SharkAgent/pkg/log"
)

// options 包 包含初始化 apiserver 的 flags 和 options

// Options 运行一个 apiserver 所需要的配置
// 在命令行 或 配置文件中进行配置
type Options struct {
	// 日志
	Log *log.Options `json:"log"      mapstructure:"log"`

	// apiServer
	GenericServerRunOptions *genericoptions.ServerRunOptions       `json:"server"   mapstructure:"server"`
	InsecureServing         *genericoptions.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	SecureServing           *genericoptions.SecureServingOptions   `json:"secure"   mapstructure:"secure"`
	FeatureOptions          *genericoptions.FeatureOptions         `json:"feature"  mapstructure:"feature"`

	// mysql
	MySQLOptions *genericoptions.MySQLOptions `json:"mysql"    mapstructure:"mysql"`
}

// NewOptions 使用默认参数创建一个 Options 对象
func NewOptions() *Options {
	o := Options{
		Log: log.NewOptions(),

		GenericServerRunOptions: genericoptions.NewServerRunOptions(),
		InsecureServing:         genericoptions.NewInsecureServingOptions(),
		SecureServing:           genericoptions.NewSecureServingOptions(),
		FeatureOptions:          genericoptions.NewFeatureOptions(),

		MySQLOptions: genericoptions.NewMySQLOptions(),
	}

	return &o
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))

	return fss
}

// Validate checks Options and return a slice of found errs.
func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.MySQLOptions.Validate()...)

	return errs
}
