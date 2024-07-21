package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"

	"github.com/ardanlabs/conf/v3"
	_ "github.com/ardanlabs/conf/v3"
	"go.uber.org/zap"

	"wheatherAPI/app/weather/api"
	v1 "wheatherAPI/app/weather/api/v1"
	"wheatherAPI/domain/weather/get"
)

type config struct {
	ServerAddr         string        `conf:"env:SERVER_ADDR,default:0.0.0.0:3000"`
	ServerReadTimeout  time.Duration `conf:"env:SERVER_READ_TIMEOUT,default:30s"`
	ServerWriteTimeout time.Duration `conf:"env:SERVER_WRITE_TIMEOUT,default:30s"`

	PprofServerAddr         string        `conf:"env:PPROF_SERVER_ADDR,default:0.0.0.0:3100"`
	PprofServerReadTimeout  time.Duration `conf:"env:PPROF_SERVER_READ_TIMEOUT,default:30s"`
	PprofServerWriteTimeout time.Duration `conf:"env:PPROF_SERVER_WRITE_TIMEOUT,default:30s"`
}

func main() {
	var cfg config

	_, err := conf.Parse("", &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			return
		}

	}

	if err = StartApp(cfg); err != nil {
		fmt.Println("running Reporter Doc 4111 API", zap.Error(err))
	}
}

func StartApp(cfg config) error {
	getWeatherUC := get.NewUseCase()

	// Handlers V1 and their dependencies
	apiV1 := v1.API{
		GetWeatherHandler: v1.GenerateReportHandler(getWeatherUC),
	}

	router := api.NewServer()

	// Routing
	apiV1.Routes(router)

	// Server
	server := http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      router,
		ReadTimeout:  cfg.ServerReadTimeout,
		WriteTimeout: cfg.ServerWriteTimeout,
	}

	// pprof
	pprofServer := http.Server{
		Addr:         cfg.PprofServerAddr,
		Handler:      getDebugHTTPMux(),
		ReadTimeout:  cfg.PprofServerReadTimeout,
		WriteTimeout: cfg.PprofServerWriteTimeout,
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		fmt.Println("server started",
			zap.String("address", server.Addr),
			zap.Duration("read timeout", server.ReadTimeout),
			zap.Duration("write timeout", server.WriteTimeout),
		)

		if serverErr := server.ListenAndServe(); serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
			fmt.Println("failed to listen and serve server", zap.Error(serverErr))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		fmt.Println("pprof server started", zap.String("address", cfg.PprofServerAddr))

		if err := pprofServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("failed to listen and serve pprof server", zap.Error(err))
		}
	}()

	shutdownCtx, shutdownCancel := context.WithCancel(context.Background())
	defer shutdownCancel()

	<-shutdownCtx.Done()
	wg.Wait()

	return nil
}

func getDebugHTTPMux() *http.ServeMux {
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/debug/pprof/", pprof.Index)
	httpMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	httpMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	httpMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	httpMux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return httpMux
}
