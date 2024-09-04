package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/krossroad/roundrobinecho/internal/echo"
	"github.com/krossroad/roundrobinecho/internal/env"
	lb "github.com/krossroad/roundrobinecho/internal/loadbalancer"
	"github.com/krossroad/roundrobinecho/internal/logger"
)

func main() {
	address := env.MustGet("HTTP_ADDRESS")
	serviceAddresses := strings.Split(env.MustGet("SERVICE_ADDRESSES"), ",")
	logger := logger.New("round-robin")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	echoBackends, err := bootEchoBackends(serviceAddresses)
	if err != nil {
		logger.Error("failed to boot echo backends", "error", err)
		return
	}

	rrLB, err := lb.NewRoundRobin(logger, echoBackends)
	if err != nil {
		logger.Error("failed to create round-robin load-balancer", "error", err)
		return
	}

	echoSvc := echo.NewService(logger, rrLB, echo.WithHealthCheckInterval(15*time.Second))
	go echoSvc.Monitor(ctx)

	app := NewApp(logger, echoSvc)
	mux := http.NewServeMux()
	mux.HandleFunc("/echo", app.fanoutHandler)

	logger.Info("starting round-robin server", "address", address)
	srv := &http.Server{
		Addr:        address,
		Handler:     mux,
		ReadTimeout: 5 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("server error", "error", err)
	}
}

func bootEchoBackends(serviceAddresses []string) ([]lb.Backend, error) {
	echoBackends := []lb.Backend{}
	for _, v := range serviceAddresses {
		b, err := echo.NewBackend(v)
		if err != nil {
			return nil, err
		}

		echoBackends = append(echoBackends, b)
	}

	return echoBackends, nil
}
