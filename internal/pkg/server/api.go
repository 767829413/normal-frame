package server

import (
	"github.com/767829413/normal-frame/internal/apiserver/options"
	"github.com/767829413/normal-frame/internal/pkg/logger"
	extDep "github.com/767829413/normal-frame/internal/pkg/options"
	"github.com/767829413/normal-frame/pkg/shutdown"
	"github.com/767829413/normal-frame/pkg/shutdown/shutdownmanagers/posixsignal"
)

type ApiServer struct {
	gs            *shutdown.GracefulShutdown
	genericServer *genericServer
	grpcServer    *grpcServer
	*extDep.MySQLOptions
	*extDep.RedisOptions
}

func CreateAPIServer(opts *options.Options) (*ApiServer, error) {
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	genericConfig, err := buildGenericConfig(opts)
	if err != nil {
		return nil, err
	}
	extraConfig, err := buildExtraConfig(opts)
	if err != nil {
		return nil, err
	}
	genericServer, err := NewGenericServer(genericConfig, extraConfig)
	if err != nil {
		return nil, err
	}
	server := &ApiServer{
		gs:            gs,
		genericServer: genericServer,
		MySQLOptions:  opts.MySQLOptions,
		RedisOptions:  opts.RedisOptions,
	}
	if extraConfig.EnableGRPC {
		extraServer, err := NewGrpcServer(extraConfig)
		if err != nil {
			return nil, err
		}
		server.setGrpcServer(extraServer)
	}

	return server, nil
}

func (s *ApiServer) setGrpcServer(grpcServer *grpcServer) {
	s.grpcServer = grpcServer
}

func (s *ApiServer) PrepareRun() *ApiServer {
	//TODO 初始化外部依赖 数据库,redis等等
	if s.MySQLOptions.Enabled {

	}

	if s.RedisOptions.Enabled {

	}

	// 优雅关停
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		s.genericServer.Close()
		if s.grpcServer.enable && s.grpcServer != nil {
			s.grpcServer.Close()
		}
		return nil
	}))
	return s
}

func (s *ApiServer) Run() error {
	if s.grpcServer.enable {
		go s.grpcServer.Run()
	}
	// start shutdown managers
	if err := s.gs.Start(); err != nil {
		logger.LogErrorf(nil, logger.LogNameNet, "start shutdown manager failed: %s", err.Error())
	}
	return s.genericServer.Run()
}
