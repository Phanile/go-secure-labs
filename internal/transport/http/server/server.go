package http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Phanile/go-secure-labs/internal/logger"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux    *http.ServeMux
	logger *logger.Logger
	config Config
}

func NewHTTPServer(config Config, log *logger.Logger) *HTTPServer {
	return &HTTPServer{
		mux:    http.NewServeMux(),
		config: config,
		logger: log,
	}
}

func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		s.mux.Handle(prefix+"/", http.StripPrefix(prefix, router))
	}
}

func (s *HTTPServer) RegisterStatic(path, dir string) {
	fileServer := http.FileServer(http.Dir(dir))

	s.mux.Handle(path, http.StripPrefix(path, fileServer))
}

func (s *HTTPServer) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: s.mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.logger.Info("start HTTP Server:", zap.String("Addr", s.config.Addr))

		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				ch <- err
			}
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve http: %w", err)
		}
	case <-ctx.Done():
		s.logger.Warn("shutdown http server...")

		shutdownContext, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownContext); err != nil {
			server.Close()

			return fmt.Errorf("shutdown http server: %w", err)
		}
	}

	s.logger.Warn("http server stopped")

	return nil
}
