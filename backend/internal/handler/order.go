package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/online-cake-shop/backend/internal/middleware"
	"github.com/online-cake-shop/backend/internal/service"
)

type OrderHandler struct {
	orderSvc *service.OrderService
}

func NewOrderHandler(orderSvc *service.OrderService) *OrderHandler {
	return &OrderHandler{orderSvc: orderSvc}
}

type createOrderRequest struct {
	DeliveryAddress string `json:"delivery_address"`
	DeliveryDate    string `json:"delivery_date"` // RFC3339
	Notes           string `json:"notes"`
	PaymentMethod   string `json:"payment_method"`
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	var req createOrderRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, envelope{"success": false, "error": "invalid request body"})
		return
	}

	deliveryDate, err := time.Parse(time.RFC3339, req.DeliveryDate)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, envelope{
			"success": false,
			"error":   "invalid delivery_date format, use RFC3339 (e.g. 2024-12-25T10:00:00Z)",
		})
		return
	}

	order, err := h.orderSvc.CreateOrder(r.Context(), service.CreateOrderInput{
		UserID:          userID,
		DeliveryAddress: req.DeliveryAddress,
		DeliveryDate:    deliveryDate,
		Notes:           req.Notes,
		PaymentMethod:   req.PaymentMethod,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeSuccess(w, http.StatusCreated, order)
}

func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	page := queryInt(r.URL.Query().Get("page"), 1)
	limit := queryInt(r.URL.Query().Get("limit"), 10)

	out, err := h.orderSvc.ListOrders(r.Context(), userID, page, limit)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeSuccess(w, http.StatusOK, out)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	orderID := chi.URLParam(r, "id")

	order, err := h.orderSvc.GetOrder(r.Context(), userID, orderID)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeSuccess(w, http.StatusOK, order)
}
