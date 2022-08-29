package server

import (
	"fmt"

	gormPlugin "github.com/767829413/normal-frame/fork/SkyAPM/go2sky-plugins/gorm"

	v3 "github.com/767829413/normal-frame/fork/SkyAPM/go2sky-plugins/gin/v3"
	redisSkyHook "github.com/767829413/normal-frame/fork/SkyAPM/go2sky-plugins/redis-go2sky-hook"
	"github.com/767829413/normal-frame/internal/apiserver/options"
	"github.com/767829413/normal-frame/internal/pkg/logger"
	extDep "github.com/767829413/normal-frame/internal/pkg/options"
	"github.com/767829413/normal-frame/internal/pkg/store"
	"github.com/767829413/normal-frame/pkg/apm"
	"github.com/767829413/normal-frame/pkg/shutdown"
	"github.com/767829413/normal-frame/pkg/shutdown/shutdownmanagers/posixsignal"
)

type ApiServer struct {
	gs            *shutdown.GracefulShutdown
	genericServer *genericServer
	grpcServer    *grpcServer
	*extDep.MySQLOptions
	*extDep.RedisOptions
	*extDep.ApmOptions
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
	tracer := apm.GetApmTracer(s.ApmOptions)
	if s.ApmOptions.Http && tracer != nil {
		s.genericServer.Use(v3.Middleware(s.genericServer.Engine, tracer.Tracer))
	}

	//初始化外部依赖 数据库,redis等等
	st := store.GetMySQLIncOr(s.MySQLOptions)
	if st != nil {
		if s.ApmOptions.Mysql && tracer != nil {
			err := st.GetDb().Use(gormPlugin.New(tracer.Tracer,
				gormPlugin.WithPeerAddr(fmt.Sprintf("%s:%d", s.MySQLOptions.Host, s.MySQLOptions.Port)),
				gormPlugin.WithSqlDBType(gormPlugin.MYSQL),
				gormPlugin.WithParamReport(),
				gormPlugin.WithQueryReport(),
			))
			if err != nil {
				logger.LogErrorf(nil, logger.LogNameMysql, "mysql set apm plugin,error: %v", err)
			}
		}
		// TODO 如果有数据库的话要执行数据库迁移

	}

	r := store.GetRedisIncOr(s.RedisOptions)
	if r != nil {
		if s.ApmOptions.Redis && tracer != nil {
			r.Getclient().AddHook(redisSkyHook.NewSkyWalkingHook(tracer.Tracer))
		}
	}

	// 优雅关停
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		tr := apm.GetApmTracer(nil)
		if tr != nil {
			_ = tr.Close()
		}

		st := store.GetMySQLIncOr(nil)
		if st != nil {
			_ = st.Close()
		}

		r := store.GetRedisIncOr(nil)
		if r != nil {
			_ = r.Close()
		}

		if s.genericServer != nil {
			s.genericServer.Close()
		}

		if s.grpcServer != nil {
			s.grpcServer.Close()
		}

		return nil
	}))
	return s
}

func (s *ApiServer) Run() error {
	if s.grpcServer != nil && s.grpcServer.enable {
		go s.grpcServer.Run()
	}
	// start shutdown managers
	if err := s.gs.Start(); err != nil {
		logger.LogErrorf(nil, logger.LogNameNet, "start shutdown manager failed: %s", err.Error())
	}
	return s.genericServer.Run()
}
