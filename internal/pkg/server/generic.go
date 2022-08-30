package server

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/767829413/normal-frame/internal/apiserver/options"
	customerRouter "github.com/767829413/normal-frame/internal/apiserver/router"
	"github.com/767829413/normal-frame/internal/pkg/config"
	"github.com/767829413/normal-frame/pkg/middleware"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"golang.org/x/sync/errgroup"
)

type genericServer struct {
	healthz       bool
	bindAddress   string
	bindPort      int
	enabledGzip   bool
	gzipLevel     int
	enableMetrics bool
	enablePprof   bool

	enableHttps  bool
	httpsAddress string
	httpsPort    int
	certFile     string
	keyFile      string

	// ShutdownTimeout is the timeout used for server shutdown. This specifies the timeout before server
	// gracefully shutdown returns.
	ShutdownTimeout time.Duration
	*gin.Engine
	http, https *http.Server
}

func NewGenericServer(genericConfig *config.GenericConfig, extraConfig *config.ExtraConfig) (*genericServer, error) {
	// setMode before gin.New()
	gin.SetMode(genericConfig.Mode)

	s := &genericServer{
		healthz:       genericConfig.Healthz,
		bindAddress:   genericConfig.BindAddress,
		bindPort:      genericConfig.BindPort,
		enabledGzip:   genericConfig.EnabledGzip,
		gzipLevel:     genericConfig.GzipLevel,
		enableMetrics: genericConfig.EnableMetrics,
		enablePprof:   genericConfig.EnablePprof,
		enableHttps:   extraConfig.EnableHttps,
		httpsAddress:  extraConfig.HttpsAddress,
		httpsPort:     extraConfig.HttpsPort,
		certFile:      extraConfig.CertFile,
		keyFile:       extraConfig.KeyFile,
		Engine:        gin.New(),
	}
	s.initGenericAPIServer()
	return s, nil
}

func (s *genericServer) initGenericAPIServer() {
	s.Setup()
	s.InstallMiddlewares()
	s.InstallAPIs()
}

// Setup do some setup work for gin engine.
func (s *genericServer) Setup() {
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("%-6s %-s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}

func (s *genericServer) InstallMiddlewares() {
	// necessary middlewares
	if s.enabledGzip {
		s.Use(middleware.Gzip(s.gzipLevel))
	}
	s.Use(middleware.Dump())
	s.Use(middleware.Cors())
	s.Use(middleware.Recovery())
	// install custom middlewares
	for k, m := range middleware.Middlewares {
		log.Printf("install middleware: %s", k)
		s.Use(m)
	}
}

func (s *genericServer) InstallAPIs() {
	// install healthz handler
	if s.healthz {
		s.GET("/healthcheck", func(context *gin.Context) {
			context.String(http.StatusOK, "OK")
		})
	}

	// install metric handler
	if s.enableMetrics {
		prometheus := ginprometheus.NewPrometheus("gin")
		prometheus.Use(s.Engine)
	}

	// install pprof handler
	if s.enablePprof {
		pprof.Register(s.Engine)
	}

	// install customer API
	customerRouter.InitRouter(s.Engine)
}

func (s *genericServer) Run() error {
	var eg errgroup.Group
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below

	httpAddr := net.JoinHostPort(s.bindAddress, strconv.Itoa(s.bindPort))
	s.http = &http.Server{
		Addr:    httpAddr,
		Handler: s,
		// ReadTimeout:    10 * time.Second,
		// WriteTimeout:   10 * time.Second,
		// MaxHeaderBytes: 1 << 20,

	}
	eg.Go(func() error {
		log.Printf("Start to listening the incoming requests on http address: %s", httpAddr)

		if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Start to listening the incoming requests on http failed : %s", err.Error())
			return err
		}

		log.Printf("Server on %s stopped", httpAddr)

		return nil
	})

	if s.enableHttps && (s.certFile != "" && s.keyFile != "" && s.httpsPort != 0) {
		httpsAddr := net.JoinHostPort(s.httpsAddress, strconv.Itoa(s.httpsPort))
		s.https = &http.Server{
			Addr:    httpsAddr,
			Handler: s,
			// ReadTimeout:    10 * time.Second,
			// WriteTimeout:   10 * time.Second,
			// MaxHeaderBytes: 1 << 20,
		}

		eg.Go(func() error {
			log.Printf("Start to listening the incoming requests on https address: %s", httpsAddr)
			if err := s.https.ListenAndServeTLS(s.certFile, s.keyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Printf("Start to listening the incoming requests on https failed : %s", err.Error())
				return err
			}

			log.Printf("Server on %s stopped", httpsAddr)

			return nil
		})

	}

	if err := eg.Wait(); err != nil {
		log.Printf("eg.Wait() failed : %s", err.Error())
		return err
	}
	return nil
}

// Close graceful shutdown the api server.
func (s *genericServer) Close() {
	// The context is used to inform the server it has 10 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := s.http.Shutdown(ctx); err != nil {
		log.Printf("Shutdown http server failed: %s", err.Error())
	}
	if s.enableHttps && s.https != nil {
		if err := s.https.Shutdown(ctx); err != nil {
			log.Printf("Shutdown https server failed: %s", err.Error())
		}
	}

}

func buildGenericConfig(opts *options.Options) (genericConfig *config.GenericConfig, lastErr error) {
	genericConfig = config.NewGenericConfig()
	if lastErr = opts.GenericServerRunOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = opts.FeatureOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}
	return
}

func buildExtraConfig(opts *options.Options) (extraConfig *config.ExtraConfig, lastErr error) {
	extraConfig = config.NewExtraConfig()
	if lastErr = opts.GrpcOptions.ApplyTo(extraConfig); lastErr != nil {
		return
	}
	if lastErr = opts.SecureOptions.ApplyTo(extraConfig); lastErr != nil {
		return
	}
	if lastErr = opts.HttpsOptions.ApplyTo(extraConfig); lastErr != nil {
		return
	}
	return
}
