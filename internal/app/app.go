package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flytrap/gin_template/internal/app/config"
	"github.com/flytrap/gin_template/pkg/redis"
	"github.com/jinzhu/copier"
	logger "github.com/sirupsen/logrus"
)

type options struct {
	ConfigFile string
	InitFile   string
	WWWDir     string
	Version    string
}

type Option func(*options)

func SetConfigFile(s string) Option {
	return func(o *options) {
		o.ConfigFile = s
	}
}

func SetInitFile(s string) Option {
	return func(o *options) {
		o.InitFile = s
	}
}

func SetVersion(s string) Option {
	return func(o *options) {
		o.Version = s
	}
}

func Import(ctx context.Context, opts ...Option) error {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	config.MustLoad(o.ConfigFile)
	config.PrintWithJSON()

	injector, injectorCleanFunc, err := BuildInjector()
	if len(o.InitFile) > 0 {
		err = injector.ImportService.Import(o.InitFile)
	}
	defer injectorCleanFunc()
	return err
}

func Init(ctx context.Context, opts ...Option) (func(), error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}

	config.MustLoad(o.ConfigFile)
	if v := o.WWWDir; v != "" {
		config.C.WWW = v
	}
	config.PrintWithJSON()

	logger.Info(fmt.Sprintf("Start server,#run_mode %s,#version %s,#pid %d", config.C.RunMode, o.Version, os.Getpid()))

	injector, injectorCleanFunc, err := BuildInjector()
	if err != nil {
		return nil, err
	}

	httpServerCleanFunc := InitHTTPServer(ctx, injector.Engine)
	return func() {
		httpServerCleanFunc()
		injectorCleanFunc()
	}, nil
}

func InitHTTPServer(ctx context.Context, handler http.Handler) func() {
	cfg := config.C.HTTP
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		logger.Info("HTTP server is running at %s.", addr)

		var err error
		if cfg.CertFile != "" && cfg.KeyFile != "" {
			srv.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			err = srv.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}

	}()

	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(cfg.ShutdownTimeout))
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error(err)
		}
	}
}

func InitStore() (*redis.Store, func(), error) {
	cfg := config.C.Redis
	c := redis.Config{}
	copier.Copy(&c, cfg)
	store := redis.NewStore(&c)
	return store, func() { store.Close() }, nil
}

func Run(ctx context.Context, opts ...Option) error {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cleanFunc, err := Init(ctx, opts...)
	if err != nil {
		return err
	}

EXIT:
	for {
		sig := <-sc
		logger.Info("Receive signal[%s]", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	cleanFunc()
	logger.Info("Server exit")
	time.Sleep(time.Second)
	os.Exit(state)
	return nil
}
