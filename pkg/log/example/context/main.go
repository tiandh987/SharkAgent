package main

import (
	"context"
	"flag"
	"github.com/tiandh987/SharkAgent/pkg/log"
)

var (
	h bool

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

	// logger配置
	opts := &log.Options{
		Level:            "debug",
		Format:           "console",
		EnableColor:      true,
		DisableCaller:    true,
		OutputPaths:      []string{"test.log", "stdout"},
		ErrorOutputPaths: []string{"error.log"},
	}
	// 初始化全局logger
	log.Init(opts)
	defer log.Flush()

	// WithValues使用
	lv := log.WithValues("X-Request-ID", "7a7b9f24-4cae-4b2a-9464-69088b45b904")

	// Context使用
	lv.Infof("Start to call pirntString function")

	// 2022-04-07 16:11:19.605	INFO	Start to call pirntString function	{"X-Request-ID": "7a7b9f24-4cae-4b2a-9464-69088b45b904"}


	ctx := lv.WithContext(context.Background())
	pirntString(ctx, "World")

	//2022-04-07 16:11:19.625	INFO	Hello World	{"X-Request-ID": "7a7b9f24-4cae-4b2a-9464-69088b45b904"}
}

func pirntString(ctx context.Context, str string) {
	lc := log.FromContext(ctx)
	lc.Infof("Hello %s", str)
}
