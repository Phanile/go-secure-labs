package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	products_transport_http "github.com/Phanile/go-secure-labs/internal/features/products/transport/http"
	"github.com/Phanile/go-secure-labs/internal/logger"
	http_middleware "github.com/Phanile/go-secure-labs/internal/transport/http/middleware"
	http_server "github.com/Phanile/go-secure-labs/internal/transport/http/server"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger, err := logger.NewLogger(logger.NewConfigMust())

	if err != nil {
		fmt.Println("failed to init app logger: ", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Info("successfully seed data...")

	productsTransportHTTP := products_transport_http.NewProductsHTTPHandler(nil)

	productsRoutes := productsTransportHTTP.Routes()

	apiProductsRouter := http_server.NewAPIVersionRouter(http_server.APIVersion1)
	apiProductsRouter.RegisterRoutes(productsRoutes...)

	httpServer := http_server.NewHTTPServer(
		http_server.NewConfigMust(),
		logger,
		http_middleware.RequestID(),
		http_middleware.Logger(logger),
		http_middleware.Panic(),
		http_middleware.Trace(),
	)

	httpServer.RegisterAPIRouters(apiProductsRouter)
	httpServer.RegisterStatic("/", "./web")

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run: ", zap.Error(err))
	}

	logger.Info("successfuly run http server...")

	<-ctx.Done()
}
