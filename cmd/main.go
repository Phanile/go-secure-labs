package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	products_postgres_repository "github.com/Phanile/go-secure-labs/internal/features/products/repository/postgres"
	products_service "github.com/Phanile/go-secure-labs/internal/features/products/service"
	products_transport_http "github.com/Phanile/go-secure-labs/internal/features/products/transport/http"
	"github.com/Phanile/go-secure-labs/internal/logger"
	core_postgres_pool "github.com/Phanile/go-secure-labs/internal/repository/postgres/pool"
	http_middleware "github.com/Phanile/go-secure-labs/internal/transport/http/middleware"
	http_server "github.com/Phanile/go-secure-labs/internal/transport/http/server"
	"github.com/Phanile/go-secure-labs/seed"
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

	pool, err := core_postgres_pool.NewConnectionPool(ctx, core_postgres_pool.NewConfigMust())

	if err != nil {
		logger.Fatal("new postrgres connection pool: %w", zap.Error(err))
	}

	defer pool.Close()

	if err := seed.RunSeed(ctx, pool); err != nil {
		logger.Error("seed data: %w", zap.Error(err))
	}

	productsRepository := products_postgres_repository.NewProductsRepository(pool)
	productsService := products_service.NewProductsService(productsRepository)

	httpServer := http_server.NewHTTPServer(
		http_server.NewConfigMust(),
		logger,
		http_middleware.RequestID(),
		http_middleware.Logger(logger),
		http_middleware.Panic(),
		http_middleware.Trace(),
	)

	productsTransportHTTP := products_transport_http.NewProductsHTTPHandler(productsService)

	apiProductsRouter := http_server.NewAPIVersionRouter(http_server.APIVersion1)
	apiProductsRouter.RegisterRoutes(productsTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiProductsRouter)
	httpServer.RegisterStatic("/", "./web")

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run: ", zap.Error(err))
	}

	logger.Info("successfuly run http server...")

	<-ctx.Done()
}
