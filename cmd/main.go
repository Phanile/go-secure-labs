package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Phanile/go-secure-labs/seed"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	connectionString := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	if err := seed.RunSeed(ctx, pool); err != nil {
		panic(err)
	}

	<-ctx.Done()
	shutdown()
}

func shutdown() {
	fmt.Println("graceful shutdown...")
}
