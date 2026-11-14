package orders_transport_http

import (
	"context"
	"encoding/json"
	core_errors "github.com/Phanile/go-secure-labs/internal/errors"
	orders_service "github.com/Phanile/go-secure-labs/internal/features/orders/service"
	"github.com/Phanile/go-secure-labs/internal/logger"
	mw "github.com/Phanile/go-secure-labs/internal/transport/http/middleware"
	response "github.com/Phanile/go-secure-labs/internal/transport/http/response"
	http_server "github.com/Phanile/go-secure-labs/internal/transport/http/server"
	"net/http"
	"strconv"
)

type Service interface {
	Create(context.Context, int64, orders_service.CreateOrderInput) (orders_service.Order, error)
	List(context.Context, int64, bool) ([]orders_service.Order, error)
	Get(context.Context, int64, int64, bool) (orders_service.Order, error)
}
type Handler struct{ service Service }

func New(service Service) *Handler { return &Handler{service: service} }
func (h *Handler) Routes() []http_server.Route {
	return []http_server.Route{{Method: http.MethodPost, Path: "/orders", Handler: h.create}, {Method: http.MethodGet, Path: "/orders", Handler: h.list}, {Method: http.MethodGet, Path: "/orders/{id}", Handler: h.get}}
}
func (h *Handler) VulnerableRoutes() []http_server.Route {
	return []http_server.Route{{Method: http.MethodGet, Path: "/vulnerable/orders", Handler: h.vulnerableList}, {Method: http.MethodGet, Path: "/vulnerable/orders/{id}", Handler: h.vulnerableGet}}
}
func userID(r *http.Request) (int64, bool) { return mw.UserIdFromContext(r.Context()) }
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())
	rh := response.NewHTTPResponseHandler(log, w)
	uid, ok := userID(r)
	if !ok {
		http.Error(w, `{"message":"authentication required"}`, http.StatusUnauthorized)
		return
	}
	var input orders_service.CreateOrderInput
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&input); err != nil {
		http.Error(w, `{"message":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	o, err := h.service.Create(r.Context(), uid, input)
	if err != nil {
		rh.ErrorResponse(err, "failed to create order")
		return
	}
	rh.JSONResponse(o, http.StatusCreated)
}
func (h *Handler) list(w http.ResponseWriter, r *http.Request)           { h.writeList(w, r, false) }
func (h *Handler) vulnerableList(w http.ResponseWriter, r *http.Request) { h.writeList(w, r, true) }
func (h *Handler) writeList(w http.ResponseWriter, r *http.Request, v bool) {
	log := logger.FromContext(r.Context())
	rh := response.NewHTTPResponseHandler(log, w)
	uid, ok := userID(r)
	if !v && !ok {
		rh.ErrorResponse(core_errors.ErrUnauthorized, "authentication required")
		return
	}
	orders, err := h.service.List(r.Context(), uid, v)
	if err != nil {
		rh.ErrorResponse(err, "failed to list orders")
		return
	}
	rh.JSONResponse(orders, http.StatusOK)
}
func (h *Handler) get(w http.ResponseWriter, r *http.Request)           { h.writeGet(w, r, false) }
func (h *Handler) vulnerableGet(w http.ResponseWriter, r *http.Request) { h.writeGet(w, r, true) }
func (h *Handler) writeGet(w http.ResponseWriter, r *http.Request, v bool) {
	log := logger.FromContext(r.Context())
	rh := response.NewHTTPResponseHandler(log, w)
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, `{"message":"invalid order id"}`, http.StatusBadRequest)
		return
	}
	uid, ok := userID(r)
	if !v && !ok {
		rh.ErrorResponse(core_errors.ErrUnauthorized, "authentication required")
		return
	}
	o, err := h.service.Get(r.Context(), uid, id, v)
	if err != nil {
		rh.ErrorResponse(err, "failed to get order")
		return
	}
	rh.JSONResponse(o, http.StatusOK)
}
