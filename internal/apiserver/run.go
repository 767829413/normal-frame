package apiserver

import (
	"github.com/767829413/normal-frame/internal/apiserver/options"
	"github.com/767829413/normal-frame/internal/pkg/logger"
	apiSver "github.com/767829413/normal-frame/internal/pkg/server"
	"github.com/767829413/normal-frame/pkg/app"
)

func GetRunFunc(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		logger.Init(opts.LogsOptions)
		return Run(opts)
	}
}

// Run runs the specified APIServer. This should never exit.
func Run(opts *options.Options) error {
	server, err := apiSver.CreateAPIServer(opts)
	if err != nil {
		return err
	}
	return server.PrepareRun().Run()
}
