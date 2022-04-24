package apiserver

import (
	"github.com/tiandh987/SharkAgent/internal/apiserver/config"
	genericapiserver "github.com/tiandh987/SharkAgent/internal/pkg/server"
	"github.com/tiandh987/SharkAgent/pkg/log"
	"github.com/tiandh987/SharkAgent/pkg/shutdown"
	"github.com/tiandh987/SharkAgent/pkg/shutdown/posixsignal"
)

// apiServer 包含 通用apiserver、优雅退出
type apiServer struct {
	// apiServer 优雅退出
	gs *shutdown.GracefulShutdown

	//redisOptions     *genericoptions.RedisOptions
	//gRPCAPIServer    *grpcAPIServer

	// apiserver 服务 - gin
	genericAPIServer *genericapiserver.GenericAPIServer
}

func (s *apiServer) PrepareRun() preparedAPIServer {
	//initRouter(s.genericAPIServer.Engine)
	//
	////s.initRedisStore()
	//
	//s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
	//	mysqlStore, _ := mysql.GetMySQLFactoryOr(nil)
	//	if mysqlStore != nil {
	//		_ = mysqlStore.Close()
	//	}
	//
	//	s.gRPCAPIServer.Close()
	//	s.genericAPIServer.Close()
	//
	//	return nil
	//}))

	return preparedAPIServer{s}
}

// ===================================================================

type preparedAPIServer struct {
	*apiServer
}

func (s preparedAPIServer) Run() error {
	//go s.gRPCAPIServer.Run()

	// start shutdown managers
	if err := s.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	return s.genericAPIServer.Run()
}

// ==================================================================

func createAPIServer(cfg *config.Config) (*apiServer, error) {
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}

	server := &apiServer{
		gs:               gs,
		genericAPIServer: genericServer,
	}

	return server, nil
}

// 基于 apiServer 的配置，生成 genericapiserver 的配置
func buildGenericConfig(cfg *config.Config) (genericConfig *genericapiserver.Config, lastErr error) {
	// 默认 genericapiserver 配置
	genericConfig = genericapiserver.NewConfig()

	// 将 apiServer 配置中的 通用服务器运行Options 合并到 genericConfig
	if lastErr = cfg.GenericServerRunOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	// HTTP 配置
	if lastErr = cfg.InsecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	// TLS 配置
	if lastErr = cfg.SecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	// Feature 配置
	if lastErr = cfg.FeatureOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	return
}
