package apiserver

import (
	"github.com/tiandh987/SharkAgent/internal/apiserver/config"
	"github.com/tiandh987/SharkAgent/internal/apiserver/options"
	"github.com/tiandh987/SharkAgent/pkg/app"
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

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		//log.Init(opts.Log)
		//defer log.Flush()

		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}

		return Run(cfg)
	}
}