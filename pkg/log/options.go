package log

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

const (
	flagLevel             = "log.level"
	flagDisableCaller     = "log.disable-caller"
	flagDisableStacktrace = "log.disable-stacktrace"
	flagFormat            = "log.format"
	flagEnableColor       = "log.enable-color"
	flagOutputPaths       = "log.output-paths"
	flagErrorOutputPaths  = "log.error-output-paths"
	flagDevelopment       = "log.development"
	flagName              = "log.name"

	consoleFormat = "console"
	jsonFormat    = "json"
)

// Options 包含与 log 包相关的配置项。
type Options struct {
	// OutputPaths，可以设置日志输出， 支持同时输出到多个输出。
	OutputPaths       []string `json:"output-paths"       mapstructure:"output-paths"`

	ErrorOutputPaths  []string `json:"error-output-paths" mapstructure:"error-output-paths"`

	Level             string   `json:"level"              mapstructure:"level"`

	// Format 支持 console 和 json 2 种格式
	Format            string   `json:"format"             mapstructure:"format"`

	DisableCaller     bool     `json:"disable-caller"     mapstructure:"disable-caller"`

	DisableStacktrace bool     `json:"disable-stacktrace" mapstructure:"disable-stacktrace"`

	// EnableColor 为 true 开启颜色输出，为 false 关闭颜色输出。
	EnableColor       bool     `json:"enable-color"       mapstructure:"enable-color"`

	Development       bool     `json:"development"        mapstructure:"development"`

	Name              string   `json:"name"               mapstructure:"name"`
}

// NewOptions 使用默认的参数，创建一个 Options 对象
func NewOptions() *Options {
	return &Options{
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		Level:             zapcore.InfoLevel.String(),
		Format:            consoleFormat,
		DisableCaller:     false,
		DisableStacktrace: false,
		EnableColor:       false,
		Development:       false,
	}
}

// AddFlags 添加 log 包的 flags 到指定的 FlagSet 对象。
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(&o.OutputPaths, flagOutputPaths, o.OutputPaths,
		"Output paths of log.")
	fs.StringSliceVar(&o.ErrorOutputPaths, flagErrorOutputPaths, o.ErrorOutputPaths,
		"Error output paths of log.")

	fs.StringVar(&o.Level, flagLevel, o.Level,
		"Minimum log output `LEVEL`.")
	fs.StringVar(&o.Format, flagFormat, o.Format,
		"Log output `FORMAT`, support plain or json format.")

	fs.BoolVar(&o.DisableCaller, flagDisableCaller, o.DisableCaller,
		"Disable output of caller information in the log.")
	fs.BoolVar(&o.DisableStacktrace, flagDisableStacktrace, o.DisableStacktrace,
		"Disable the log to record a stack trace for all messages at or above panic level.")
	fs.BoolVar(&o.EnableColor, flagEnableColor, o.EnableColor,
		"Enable output ansi colors in plain format logs.")
	fs.BoolVar(&o.Development, flagDevelopment, o.Development,
		"Development puts the logger in development mode, which changes "+
			"the behavior of DPanicLevel and takes stacktraces more liberally.",
	)

	fs.StringVar(&o.Name, flagName, o.Name, "The name of the logger.")
}

// Validate 验证 Options 字段
func (o *Options) Validate() []error {
	var errs []error

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, err)
	}

	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != jsonFormat {
		errs = append(errs, fmt.Errorf("not a valid log format: %q", o.Format))
	}

	return errs
}

// String
func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

// Build 从 Config 和 Options 构造一个全局 zap logger。
func (o *Options) Build() error {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	encodeLevel := zapcore.CapitalLevelEncoder
	if o.Format == consoleFormat && o.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zc := &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Development:       o.Development,
		DisableCaller:     o.DisableCaller,
		DisableStacktrace: o.DisableStacktrace,
		Sampling:          &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:          o.Format,
		EncoderConfig:     zapcore.EncoderConfig{
			MessageKey:          "message",
			LevelKey:            "level",
			TimeKey:             "timestamp",
			NameKey:             "logger",
			CallerKey:           "caller",
			StacktraceKey:       "stacktrace",
			LineEnding:          zapcore.DefaultLineEnding,
			EncodeLevel:         encodeLevel,
			EncodeTime:          timeEncoder,
			EncodeDuration:      milliSecondsDurationEncoder,
			EncodeCaller:        zapcore.ShortCallerEncoder,
			EncodeName:          zapcore.FullNameEncoder,
		},
		OutputPaths:       o.OutputPaths,
		ErrorOutputPaths:  o.ErrorOutputPaths,
	}

	logger, err := zc.Build(zap.AddStacktrace(zapcore.PanicLevel))
	if err != nil {
		return err
	}
	zap.RedirectStdLog(logger.Named(o.Name))
	zap.ReplaceGlobals(logger)

	return nil
}