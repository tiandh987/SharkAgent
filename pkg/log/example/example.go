package main

import (
	"context"
	"flag"
	"github.com/tiandh987/SharkAgent/pkg/log"
)

var (
	h      bool
	level  int
	format string
)

func main() {
	flag.BoolVar(&h, "h", false, "Print this help.")
	flag.IntVar(&level, "l", 0, "Log level.")
	flag.StringVar(&format, "f", "console", "log output format.")

	flag.Parse()

	if h {
		flag.Usage()

		return
	}

	// logger 配置
	opts := &log.Options{
		Level:            "debug",
		Format:           "console",
		DisableCaller:    true,
		EnableColor:      true,
		OutputPaths:      []string{"test.log", "stdout"},
		ErrorOutputPaths: []string{"error.log"},
	}

	// 初始化全局 logger
	log.Init(opts)
	defer log.Flush()

	// Debug、Info(with field)、Warnf、Errorw使用
	log.Debug("This is a debug message")
	log.Info("This is a info message", log.Int32("int_key", 10))
	log.Warnf("This is a formatted %s message", "warn")
	log.Errorw("Message printed with Errorw", "X-Request-ID", "fbf54504-64da-4088-9b86-67824a7fb508")

	// 2022-04-07 15:42:33.023	DEBUG	This is a debug message
	// 2022-04-07 15:42:33.028	INFO	This is a info message	{"int_key": 10}
	// 2022-04-07 15:42:33.028	WARN	This is a formatted warn message
	// 2022-04-07 15:42:33.028	ERROR	Message printed with Errorw	{"X-Request-ID": "fbf54504-64da-4088-9b86-67824a7fb508"}


	// WithValues使用
	lv := log.WithValues("X-Request-ID", "7a7b9f24-4cae-4b2a-9464-69088b45b904")
	lv.Infow("Info message printed with [WithValues] logger")
	lv.Debugw("Debug message printed with [WithValues] logger")

	// 2022-04-07 15:45:18.784	INFO	Info message printed with [WithValues] logger	{"X-Request-ID": "7a7b9f24-4cae-4b2a-9464-69088b45b904"}
	// 2022-04-07 15:45:18.784	DEBUG	Debug message printed with [WithValues] logger	{"X-Request-ID": "7a7b9f24-4cae-4b2a-9464-69088b45b904"}


	// Context使用
	ctx := lv.WithContext(context.Background())
	lc := log.FromContext(ctx)
	lc.Info("Message printed with [WithContext] logger")

	// 2022-04-07 15:48:44.060	INFO	Message printed with [WithContext] logger	{"X-Request-ID": "7a7b9f24-4cae-4b2a-9464-69088b45b904"}


	// WithName 使用
	ln := lv.WithName("test")
	ln.Info("Message printed with [WithName] logger")

	// 2022-04-07 15:51:03.996	INFO	test	Message printed with [WithName] logger	{"X-Request-ID": "7a7b9f24-4cae-4b2a-9464-69088b45b904"}


	// V level 使用
	log.V(5).Info("This is a V level message")
	log.V(5).Infow("This is a V level message with fields", "X-Request-ID", "7a7b9f24-4cae-4b2a-9464-69088b45b904")

	// 2022-04-07 15:58:56.950	INFO	This is a V level message
	// 2022-04-07 15:58:56.950	INFO	This is a V level message with fields	{"X-Request-ID": "7a7b9f24-4cae-4b2a-9464-69088b45b904"}
}
