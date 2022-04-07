package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

// =====================================================

var (
	std = New(NewOptions())
	mu  sync.Mutex
)

// Init 使用给定的 Options 初始化 logger
func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()
	std = New(opts)
}

// New 通过命令参数自定义的 opts 创建一个 logger
func New(opts *Options) *zapLogger {
	// 如果传入的 opts 为 nil，则使用 NewOptions() 创建默认的 opts
	if opts == nil {
		opts = NewOptions()
	}

	var zapLevel zapcore.Level
	// UnmarshalText 将 text 反序列化为一个 level。
	// 与 MarshalText 一样，UnmarshalText 期望 level 的文本表示删除 -Level 后缀（参见示例）。
	// 特别是，这使得使用 YAML、TOML 或 JSON 文件配置日志记录级别变得容易。
	//
	// 可配置的日志级别值为：
	//	"debug", "DEBUG"
	//	"info", "INFO", ""
	//	"warn", "WARN"
	//	"error", "ERROR"
	//	"dpanic", "DPANIC"
	//	"panic", "PANIC"
	//	"fatal", "FATAL"
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	// CapitalLevelEncoder 将 Level 序列化为全大写字符串。
	// 例如， InfoLevel 被序列化为“INFO”。
	encodeLevel := zapcore.CapitalLevelEncoder

	// 输出格式为 console，并且启用颜色
	if opts.Format == consoleFormat && opts.EnableColor {
		// CapitalColorLevelEncoder 将 Level 序列化为全大写字符串并添加颜色。
		// 例如，InfoLevel 被序列化为“INFO”，颜色为蓝色。
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// EncoderConfig 允许用户配置 zapcore 提供的具体编码器。
	encoderConfig := zapcore.EncoderConfig{
		// 设置用于每个日志条目的键。 如果任何键为空，则条目的该部分将被省略。
		MessageKey:    "message",
		LevelKey:      "level",
		TimeKey:       "timestamp",
		NameKey:       "logger",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,

		// 配置常见复杂类型的原始表示。
		// 例如，一些用户可能希望将所有 time.Times 序列化为自纪元以来的浮点秒数，
		// 而其他用户可能更喜欢 ISO8601 字符串。
		EncodeLevel:    encodeLevel,
		EncodeTime:     timeEncoder,
		EncodeDuration: milliSecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		// ShortCallerEncoder 以 package/file:line 格式序列化caller，
		// 从完整路径中修剪除最终目录之外的所有目录。
	}

	// Config 提供了一种声明式的方式来构造一个记录器。
	// 它不会做任何用 New、Options 和各种 zapcore.WriteSyncer 和 zapcore.Core 包装器无法完成的事情，
	// 但它是切换常用选项的更简单方法。
	//
	// 请注意，Config 有意仅支持最常见的选项。
	// 更不寻常的日志设置（记录到网络连接或消息队列，在多个文件之间拆分输出等）是可能的， 但需要直接使用 zapcore 包。
	// 有关示例代码，请参阅包级 BasicConfiguration 和 AdvancedConfiguration 示例。
	//
	// 有关显示运行时日志级别更改的示例，请参阅 AtomicLevel 的文档。
	loggerConfig := &zap.Config{
		// Level 是最低启用的日志记录级别。
		// 请注意，这是一个动态级别，因此调用 Config.Level.SetLevel
		// 将自动更改从该配置继承的所有 logger 的日志级别。
		//
		// NewAtomicLevelAt 是一个方便的函数，它创建一个 AtomicLevel，
		// 然后使用给定的级别调用 SetLevel。
		Level: zap.NewAtomicLevelAt(zapLevel),

		// Development 将记录器置于开发模式，
		// 这会改变 DPanicLevel 的行为并更自由地获取堆栈跟踪。
		Development: opts.Development,

		// DisableCaller 停止使用调用函数的文件名和行号注释日志。
		// 默认情况下，所有日志都带有注释。
		DisableCaller: opts.DisableCaller,

		// DisableStacktrace 完全禁用自动堆栈跟踪捕获。
		// 默认情况下，会为 development 中的 WarnLevel 及以上日志
		// 和 production 中的 ErrorLevel 及以上日志捕获堆栈跟踪。
		DisableStacktrace: opts.DisableStacktrace,

		// Sampling 设置采样策略。 一个 nil SamplingConfig 将禁用采样。
		//
		// SamplingConfig：
		// 	SamplingConfig 为 logger 设置采样策略。
		// 	采样限制了 logger 对您的进程施加的全局 CPU 和 I/O 负载， 同时尝试保留日志的代表性子集。
		//
		// 	如果指定，采样器将在每次决策后调用 Hook。
		//
		// 	这里配置的值是每秒。
		//	有关详细信息，请参阅 zapcore.NewSamplerWithOptions。
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},

		// Encoding 设置 logger 的编码。
		// 有效值为 “json” 和 “console”，以及通过 RegisterEncoder 注册的任何第三方编码。
		Encoding:         opts.Format,

		// EncoderConfig 设置所选编码器的选项。
		// 有关详细信息，请参阅 zapcore.EncoderConfig。
		EncoderConfig:    encoderConfig,

		// OutputPaths 是要写入日志输出的 URL 或文件路径列表。
		// 详见 Open。
		OutputPaths:      opts.OutputPaths,

		// ErrorOutputPaths 是要写入 "内部logger错误" 的 URL 列表。
		// 默认为标准错误。
		//
		// 注意这个设置只影响内部错误；
		// 有关将 error 级别日志 发送到与 info 级别和 debug 级别日志 "不同的位置" 的示例代码，
		// 请参阅包级别的 AdvancedConfiguration 示例。
		ErrorOutputPaths: opts.ErrorOutputPaths,
	}

	var err error
	l, err := loggerConfig.Build(zap.AddStacktrace(zapcore.PanicLevel),
		zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	logger := &zapLogger{
		zapLogger: l.Named(opts.Name),
		infoLogger: infoLogger{
			log: 1,
			level: zap.InfoLevel,
		},
	}

	klog.InitLogger(l)
	zap.RedirectStdLog(l)

	return logger
}
