package apiserver

import (
	"github.com/767829413/normal-frame/internal/apiserver/options"
	"github.com/767829413/normal-frame/pkg/app"
)

const commandDesc = `Generic Golang server for validating and configuring data Validates and configures data for api objects, which include users, policies, secrets, and Others. API server serves REST operations for api object management.`

// NewApp creates an App object with default parameters.
func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	return app.NewApp("API Server",
		basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(GetRunFunc(opts)),
	)
}
