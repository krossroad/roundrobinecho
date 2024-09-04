package main

import (
	"context"
	"net/http"
	"os"
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
		os.Exit(1)
	}

	rrLB, err := lb.NewRoundRobin(logger, echoBackends)
	if err != nil {
		logger.Error("failed to create round-robin load-balancer", "error", err)
		os.Exit(1)
	}

	echoSvc := echo.NewService(logger, rrLB, echo.WithHealthCheckInterval(30*time.Second))
	go echoSvc.Monitor(ctx)

	app := NewApp(logger, echoSvc)
	mux := http.NewServeMux()
	mux.HandleFunc("/echo", app.fanoutHandler)

	logger.Info("starting round-robin server", "address", address)
	if err := http.ListenAndServe(address, mux); err != nil {
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
