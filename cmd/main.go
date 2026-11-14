package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	auth_postgres_repository "github.com/Phanile/go-secure-labs/internal/features/auth/repository/postgres"
	auth_service "github.com/Phanile/go-secure-labs/internal/features/auth/service"
	auth_transport_http "github.com/Phanile/go-secure-labs/internal/features/auth/transport"
	orders_postgres_repository "github.com/Phanile/go-secure-labs/internal/features/orders/repository/postgres"
	orders_service "github.com/Phanile/go-secure-labs/internal/features/orders/service"
	orders_transport_http "github.com/Phanile/go-secure-labs/internal/features/orders/transport/http"
	products_postgres_repository "github.com/Phanile/go-secure-labs/internal/features/products/repository/postgres"
	products_postgres_vulnerable_repository "github.com/Phanile/go-secure-labs/internal/features/products/repository/postgres/vulnerable"
	products_service "github.com/Phanile/go-secure-labs/internal/features/products/service"
	products_transport_http "github.com/Phanile/go-secure-labs/internal/features/products/transport/http"
	users_postgres_repository "github.com/Phanile/go-secure-labs/internal/features/users/repository/postgres"
	users_service "github.com/Phanile/go-secure-labs/internal/features/users/service"
	"github.com/Phanile/go-secure-labs/internal/hasher"
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
	productsVulnerableRepository := products_postgres_vulnerable_repository.NewProductsVulnerableRepository(pool)

	userRepository := users_postgres_repository.NewUserRepository(pool)

	authRepository := auth_postgres_repository.NewAuthRepository(pool)

	productsService := products_service.NewProductsService(productsRepository)
	productsVulnerableService := products_service.NewProductsService(productsVulnerableRepository)
	ordersService := orders_service.New(orders_postgres_repository.New(pool))

	userService := users_service.NewUserService(userRepository, hasher.SHA256Hasher{})

	authService := auth_service.NewAuthService(userService, authRepository, hasher.SHA256Hasher{}, hasher.SHA256Hasher{})

	httpServer := http_server.NewHTTPServer(
		http_server.NewConfigMust(),
		logger,
		http_middleware.RequestID(),
		http_middleware.Logger(logger),
		http_middleware.Panic(),
		http_middleware.Trace(),
		http_middleware.Auth(authService),
	)

	productsTransportHTTP := products_transport_http.NewProductsHTTPHandler(productsService)
	productsVulnerableTransportHTTP := products_transport_http.NewProductsHTTPHandler(productsVulnerableService)

	authTransportHTTP := auth_transport_http.NewAuthHTTPHandler(authService)
	ordersTransportHTTP := orders_transport_http.New(ordersService)

	apiProductsRouter := http_server.NewAPIVersionRouter(http_server.APIVersion1)
	apiProductsRouter.RegisterRoutes(productsTransportHTTP.Routes()...)
	apiProductsRouter.RegisterRoutes(productsVulnerableTransportHTTP.VulnerableRoutes()...)
	apiProductsRouter.RegisterRoutes(authTransportHTTP.Routes()...)
	apiProductsRouter.RegisterRoutes(ordersTransportHTTP.Routes()...)
	apiProductsRouter.RegisterRoutes(ordersTransportHTTP.VulnerableRoutes()...)

	httpServer.RegisterAPIRouters(apiProductsRouter)
	httpServer.RegisterStatic("/", "./web")

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run: ", zap.Error(err))
	}

	logger.Info("successfuly run http server...")

	<-ctx.Done()
}
