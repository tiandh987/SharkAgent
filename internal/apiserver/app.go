package apiserver

import (
	"github.com/tiandh987/SharkAgent/internal/apiserver/config"
	"github.com/tiandh987/SharkAgent/internal/apiserver/options"
	"github.com/tiandh987/SharkAgent/pkg/app"
	"github.com/tiandh987/SharkAgent/pkg/log"
)

func NewApp(basename string) *app.App {
	opts := options.NewOptions()

	application := app.NewApp("IAM API Server",
		basename,
		app.WithOptions(opts),
		app.WithRunFunc(run(opts)),
	)

	return application
}

// 传入的 opts 是从命令行、配置文件获取的
func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		log.Init(opts.Log)
		defer log.Flush()

		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}

		return Run(cfg)
	}
}