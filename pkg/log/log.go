package log

import (
	"context"
	"fmt"
	"github.com/tiandh987/SharkAgent/pkg/log/klog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"sync"
)

// 文件内容：
//	1、std
//		全局 zapLogger
//	2、Init()
//		初始化 std
// 	3、New()
//		通过给定的 Options 创建一个 zapLogger
//	4、实现 Logger 接口定义的方法
//
//	1、type InfoLogger interface
//	2、type infoLogger struct
//	3、type noopInfoLogger struct
//		disabledInfoLogger
//
//	1、type Logger interface
//	2、type zapLogger struct

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
		Encoding: opts.Format,

		// EncoderConfig 设置所选编码器的选项。
		// 有关详细信息，请参阅 zapcore.EncoderConfig。
		EncoderConfig: encoderConfig,

		// OutputPaths 是要写入日志输出的 URL 或文件路径列表。
		// 详见 Open。
		OutputPaths: opts.OutputPaths,

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
			log:   l,
			level: zap.InfoLevel,
		},
	}

	klog.InitLogger(l)
	zap.RedirectStdLog(l)

	return logger
}

// Debug
func Debug(msg string, fields ...Field) {
	std.zapLogger.Debug(msg, fields...)
}

func Debugf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Debugf(format, v...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Debugw(msg, keysAndValues...)
}

// Info
func Info(msg string, fields ...Field) {
	std.zapLogger.Info(msg, fields...)
}

func Infof(format string, v ...interface{}) {
	std.zapLogger.Sugar().Infof(format, v...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Infow(msg, keysAndValues...)
}

// Warn
func Warn(msg string, fields ...Field) {
	std.zapLogger.Warn(msg, fields...)
}

func Warnf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Warnf(format, v...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Warnw(msg, keysAndValues...)
}

// Error
func Error(msg string, fields ...Field) {
	std.zapLogger.Error(msg, fields...)
}

func Errorf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Errorf(format, v...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Errorw(msg, keysAndValues...)
}

// Panic
func Panic(msg string, fields ...Field) {
	std.zapLogger.Panic(msg, fields...)
}

func Panicf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Panicf(format, v...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Panicw(msg, keysAndValues...)
}

// Fatal
func Fatal(msg string, fields ...Field) {
	std.zapLogger.Fatal(msg, fields...)
}

func Fatalf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Fatalf(format, v...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Fatalw(msg, keysAndValues...)
}

// V return a leveled InfoLogger.
func V(level int) InfoLogger {
	return std.V(level)
}

// L method output with specified context value.
func L(ctx context.Context) *zapLogger {
	return std.L(ctx)
}

func Write(p []byte) (n int, err error) {
	return std.Write(p)
}

// WithValues creates a child logger and adds adds Zap fields to it.
func WithValues(keysAndValues ...interface{}) Logger {
	return std.WithValues(keysAndValues...)
}

// WithName adds a new path segment to the logger's name.
// Segments are joined by periods.
// By default, Loggers are unnamed.
func WithName(s string) Logger {
	return std.WithName(s)
}

// WithContext returns a copy of context in which the log value is set.
func WithContext(ctx context.Context) context.Context {
	return std.WithContext(ctx)
}

// Flush calls the underlying Core's Sync method, flushing any buffered log entries.
// Applications should take care to call Sync before exiting.
func Flush() {
	std.Flush()
}

// ===================================================================

// InfoLogger 表示以特定详细程度记录 non-error 消息的能力。
type InfoLogger interface {
	// Info 以给定的 key/value对 作为上下文, 记录一条 non-error 消息。
	//
	// msg 参数应该用于向日志行添加一些常量描述。
	// 可以使用 key/value对 添加其他变量信息。
	// key/value对 应该交替字符串键和任意值。
	Info(msg string, fields ...Field)
	Infof(format string, v ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	// Info 使用指定的 key/value 记录日志。
	// Infof 格式化记录日志。
	// Infow 也是使用指定的 key/value 记录日志。
	//
	// Infow 跟 Info 的区别是：
	//	使用 Info 需要指定值的类型，通过指定值的日志类型，日志库底层不需要进行反射操作，
	// 	所以使用 Info 记录日志性能最高。

	// Enabled 测试此 InfoLogger 是否已启用。
	// 例如，命令行标志可能用于设置日志记录的详细程度, 并禁用某些信息日志。
	Enabled() bool
}

//====================================================================

// infoLogger 是一个 logr.InfoLogger，它使用 Zap 在特定 level 进行日志记录。
// 该 level 已经转换为 Zap level，即 `logrLevel = -1*zapLevel`。
type infoLogger struct {
	level zapcore.Level
	log   *zap.Logger
}

func (l *infoLogger) Enabled() bool {
	return true
}

func (l *infoLogger) Info(msg string, fields ...Field) {
	if checkedEntry := l.log.Check(l.level, msg); checkedEntry != nil {
		checkedEntry.Write(fields...)
	}
}

func (l *infoLogger) Infof(format string, args ...interface{}) {
	if checkedEntry := l.log.Check(l.level, fmt.Sprintf(format, args...)); checkedEntry != nil {
		checkedEntry.Write()
	}
}

func (l *infoLogger) Infow(msg string, keysAndValues ...interface{}) {
	if checkedEntry := l.log.Check(l.level, msg); checkedEntry != nil {
		checkedEntry.Write(handleFields(l.log, keysAndValues)...)
	}
}

// ===================================================================

// noopInfoLogger 是一个 logr.InfoLogger，它总是被禁用，并且什么都不做。
type noopInfoLogger struct{}

func (l *noopInfoLogger) Enabled() bool {
	return false
}

func (l *noopInfoLogger) Info(_ string, _ ...Field) {}

func (l *noopInfoLogger) Infof(_ string, _ ...interface{}) {}

func (l *noopInfoLogger) Infow(_ string, _ ...interface{}) {}

var disabledInfoLogger = &noopInfoLogger{}

// ===================================================================

// Logger 表示记录消息的能力，包括错误和非错误。
type Logger interface {
	// InfoLogger
	// 所有 Logger 都实现 InfoLogger。
	// 直接在 Logger 值上调用 InfoLogger 方法, 等效于在 V(0) InfoLogger 上调用它们。
	// 例如，logger.Info() 产生与 logger.V(0).Info 相同的结果。
	InfoLogger

	// log 包为每种级别的日志都提供了 3 种日志记录方式

	Debug(msg string, fields ...Field)
	Debugf(format string, v ...interface{})
	Debugw(msg string, keysAndValues ...interface{})

	Warn(msg string, fields ...Field)
	Warnf(format string, v ...interface{})
	Warnw(msg string, keysAndValues ...interface{})

	Error(msg string, fields ...Field)
	Errorf(format string, v ...interface{})
	Errorw(msg string, keysAndValues ...interface{})

	Panic(msg string, fields ...Field)
	Panicf(format string, v ...interface{})
	Panicw(msg string, keysAndValues ...interface{})

	Fatal(msg string, fields ...Field)
	Fatalf(format string, v ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})

	// V 返回特定详细级别的 InfoLogger 值。
	// 更高的详细级别意味着日志消息不太重要。
	// 传递小于零的日志级别是非法的。
	//
	// log 包支持 V Level，可以通过整型数值来灵活指定日志级别，数值越大，优先级越低
	// 0: Fatal  1: Panic  2: DPanic  3: Error  4: Warn  5: Info  6: Debug
	V(level int) InfoLogger

	// L 可以很方便地从 Context 中提取出指定的 key-value，
	// 作为上下文添加到日志输出中
	L(ctx context.Context) *zapLogger

	Write(p []byte) (n int, err error)

	// WithValues 可以返回一个携带指定 key-value 的 Logger，供后面使用。
	WithValues(keysAndValues ...interface{}) Logger

	// WithName 为 logger 的名称添加一个新元素。
	// 连续调用 WithName 继续为 logger 的名称附加后缀。
	// 强烈建议名称段仅包含字母、数字和连字符
	//（有关更多信息，请参阅 package 文档）。
	WithName(name string) Logger

	// WithContext
	// log 包提供 WithContext 和 FromContext 用来将指定的 Logger 添加到某个 Context 中，
	// 以及从某个 Context 中获取 Logger
	WithContext(ctx context.Context) context.Context

	// Flush 调用底层核心的 Sync 方法，刷新所有缓冲的日志条目。
	// 应用程序应注意在退出之前调用 Sync。
	Flush()
}

// ===================================================================

var _ Logger = &zapLogger{}

// zapLogger is a logr.Logger that uses Zap to log.
type zapLogger struct {
	// 注意：这看起来与 zap.SugaredLogger 非常相似，但处理的是我们想要拥有多个详细级别的愿望。
	zapLogger *zap.Logger
	infoLogger
}

// Debug
func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.zapLogger.Debug(msg, fields...)
}

func (l *zapLogger) Debugf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Debugf(format, v...)
}

func (l *zapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Debugw(msg, keysAndValues...)
}

// Info
func (l *zapLogger) Info(msg string, fields ...Field) {
	l.zapLogger.Info(msg, fields...)
}

func (l *zapLogger) Infof(format string, v ...interface{}) {
	l.zapLogger.Sugar().Infof(format, v...)
}

func (l *zapLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Infow(msg, keysAndValues...)
}

// Warn
func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.zapLogger.Warn(msg, fields...)
}

func (l *zapLogger) Warnf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Warnf(format, v...)
}

func (l *zapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Warnw(msg, keysAndValues...)
}

// Error
func (l *zapLogger) Error(msg string, fields ...Field) {
	l.zapLogger.Error(msg, fields...)
}

func (l *zapLogger) Errorf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Errorf(format, v...)
}

func (l *zapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Errorw(msg, keysAndValues...)
}

// Panic
func (l *zapLogger) Panic(msg string, fields ...Field) {
	l.zapLogger.Panic(msg, fields...)
}

func (l *zapLogger) Panicf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Panicf(format, v...)
}

func (l *zapLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Panicw(msg, keysAndValues...)
}

// Fatal
func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.zapLogger.Fatal(msg, fields...)
}

func (l *zapLogger) Fatalf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Fatalf(format, v...)
}

func (l *zapLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Fatalw(msg, keysAndValues...)
}

func (l zapLogger) V(level int) InfoLogger {
	lvl := zapcore.Level(5 - 1*level)
	if l.zapLogger.Core().Enabled(lvl) {
		return &infoLogger{
			level: lvl,
			log:   l.zapLogger,
		}
	}
	return disabledInfoLogger
}

func (l *zapLogger) L(ctx context.Context) *zapLogger {
	lg := l.clone()

	// L() 方法会从传入的 Context 中提取出 requestID 和 username ，追加到 Logger 中，并返回 Logger。
	// 这时候调用该 Logger 的 Info、Infof、Infow 等方法记录日志， 输出的日志中均包含 requestID 和 username 字段

	if requestID := ctx.Value(KeyRequestID); requestID != nil {
		lg.zapLogger = lg.zapLogger.With(zap.Any(KeyRequestID, requestID))
	}
	if username := ctx.Value(KeyUsername); username != nil {
		lg.zapLogger = lg.zapLogger.With(zap.Any(KeyUsername, username))
	}
	if watcherName := ctx.Value(KeyWatcherName); watcherName != nil {
		lg.zapLogger = lg.zapLogger.With(zap.Any(KeyWatcherName, watcherName))
	}

	return lg
}

func (l zapLogger) Write(p []byte) (n int, err error) {
	l.zapLogger.Info(string(p))

	return len(p), nil
}

func (l *zapLogger) WithValues(keysAndValues ...interface{}) Logger {
	newLogger := l.zapLogger.With(handleFields(l.zapLogger, keysAndValues)...)

	return NewLogger(newLogger)
}

func (l *zapLogger) WithName(name string) Logger {
	newLogger := l.zapLogger.Named(name)

	return NewLogger(newLogger)
}

func (l *zapLogger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logContextKey, l)
}

func (l *zapLogger) Flush() {
	_ = l.zapLogger.Sync()
}

//nolint:predeclared
func (l *zapLogger) clone() *zapLogger {
	copy := *l

	return &copy
}

//===================================================================

// handleFields 将一堆任意 key-value对 转换为 Zap fields。
// 它需要额外的预转换 Zap fields，用于自动附加的字段，如`error`。
func handleFields(l *zap.Logger, args []interface{}, additional ...zap.Field) []zap.Field {
	// zap.SugaredLogger.sweetenFields 的略微修改版本
	if len(args) == 0 {
		// 如果我们没有 suggarded 字段，则快速返回。
		return additional
	}

	// 与 Zap 不同，我们可以很确定用户没有传递结构化字段（因为 logr 没有这个概念），
	// 所以猜测我们需要更少的空间。
	fields := make([]zap.Field, 0, len(args)/2+len(additional))
	for i := 0; i < len(args); {
		if _, ok := args[i].(zap.Field); ok {
			l.DPanic("strongly-typed Zap Field passed to logr", zap.Any("zap field", args[i]))

			break
		}

		// 确保这不是一个不匹配的键
		if i == len(args)-1 {
			l.DPanic("odd number of arguments passed as key-value pairs for logging", zap.Any("ignored key", args[i]))

			break
		}

		key, val := args[i], args[i+1]
		keyStr, isString := key.(string)
		if !isString {
			l.DPanic(
				"non-string key argument passed to logging, ignoring all later arguments",
				zap.Any("invalid key", key),
			)

			break
		}

		fields = append(fields, zap.Any(keyStr, val))
		i += 2
	}

	return append(fields, additional...)
}

// ==============================================================

// NewLogger creates a new logr.Logger using the given Zap Logger to log.
func NewLogger(l *zap.Logger) Logger {
	return &zapLogger{
		zapLogger: l,
		infoLogger: infoLogger{
			log:   l,
			level: zap.InfoLevel,
		},
	}
}

// SugaredLogger returns global sugared logger.
func SugaredLogger() *zap.SugaredLogger {
	return std.zapLogger.Sugar()
}

// ZapLogger used for other log wrapper such as klog.
func ZapLogger() *zap.Logger {
	return std.zapLogger
}

// StdErrLogger returns logger of standard library which writes to supplied zap
// logger at error level.
func StdErrLogger() *log.Logger {
	if std == nil {
		return nil
	}
	if l, err := zap.NewStdLogAt(std.zapLogger, zapcore.ErrorLevel); err == nil {
		return l
	}

	return nil
}

// StdInfoLogger returns logger of standard library which writes to supplied zap
// logger at info level.
func StdInfoLogger() *log.Logger {
	if std == nil {
		return nil
	}
	if l, err := zap.NewStdLogAt(std.zapLogger, zapcore.InfoLevel); err == nil {
		return l
	}

	return nil
}

// CheckIntLevel used for other log wrapper such as klog which return if logging a
// message at the specified level is enabled.
func CheckIntLevel(level int32) bool {
	var lvl zapcore.Level
	if level < 5 {
		lvl = zapcore.InfoLevel
	} else {
		lvl = zapcore.DebugLevel
	}
	checkEntry := std.zapLogger.Check(lvl, "")

	return checkEntry != nil
}
