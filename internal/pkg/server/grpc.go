package server

import (
	"net"

	"github.com/767829413/normal-frame/internal/pkg/config"
	"github.com/767829413/normal-frame/internal/pkg/logger"
	"google.golang.org/grpc"
)

type grpcServer struct {
	enable bool
	*grpc.Server
	address string
}

func NewGrpcServer(extraConfig *config.ExtraConfig) (*grpcServer, error) {
	// creds, err := credentials.NewServerTLSFromFile(extraConfig.CertFile, extraConfig.KeyFile)
	// if err != nil {
	// 	logger.LogErrorf(nil,logger.LogNameGRpc,"Failed to generate credentials %s", err.Error())
	// }
	// opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(extraConfig.MaxMsgSize), grpc.Creds(creds)}
	// storeIns, _ := mysql.GetMySQLFactoryOr(c.mysqlOptions)
	// storeIns, _ := etcd.GetEtcdFactoryOr(c.etcdOptions, nil)
	// store.SetClient(storeIns)

	return nil, nil
}

func (s *grpcServer) Run() {
	listen, err := net.Listen("tcp", s.address)
	if err != nil {
		logger.LogErrorf(nil, logger.LogNameGRpc, "failed to listen: %s", err.Error())
	}

	go func() {
		if err := s.Serve(listen); err != nil {
			logger.LogErrorf(nil, logger.LogNameGRpc, "failed to start grpc server: %s", err.Error())
		}
	}()
	logger.LogInfof(nil, logger.LogNameGRpc, "start grpc server at %s", s.address)
}

func (s *grpcServer) Close() {
	s.GracefulStop()
	logger.LogInfof(nil, logger.LogNameGRpc, "GRPC server on %s stopped", s.address)
}
