package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/krossroad/roundrobinecho/internal/echo"
	"github.com/krossroad/roundrobinecho/internal/env"
	"github.com/krossroad/roundrobinecho/internal/logger"
	"github.com/krossroad/roundrobinecho/internal/middlewares"
)

func main() {
	address := env.MustGet("HTTP_ADDRESS")
	serviceAddresses := strings.Split(env.MustGet("SERVICE_ADDRESSES"), ",")
	logger := logger.New("round-robin")

	rrLB, err := echo.NewRoundRobin(logger, serviceAddresses)
	if err != nil {
		logger.Error("failed to create round-robin", "error", err)
		os.Exit(1)
	}

	echoSvc, err := echo.NewService(logger, rrLB)
	if err != nil {
		logger.Error("failed to create app", "error", err)
		os.Exit(1)
	}

	app := NewApp(logger, *echoSvc)
	mux := http.NewServeMux()
	mux.Handle("/echo", middlewares.JSONContentTypeValidator(app.fanoutHandler))

	logger.Info("starting round-robin server", "address", address)
	if err := http.ListenAndServe(address, mux); err != nil {
		logger.Error("server error", "error", err)
	}
}
