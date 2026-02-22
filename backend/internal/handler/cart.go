package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/online-cake-shop/backend/internal/middleware"
	"github.com/online-cake-shop/backend/internal/service"
)

type CartHandler struct {
	cartSvc *service.CartService
}

func NewCartHandler(cartSvc *service.CartService) *CartHandler {
	return &CartHandler{cartSvc: cartSvc}
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	cart, err := h.cartSvc.GetCart(r.Context(), userID)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeSuccess(w, http.StatusOK, cart)
}

type addCartItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	var req addCartItemRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, envelope{"success": false, "error": "invalid request body"})
		return
	}

	cart, err := h.cartSvc.AddItem(r.Context(), service.AddCartItemInput{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeSuccess(w, http.StatusOK, cart)
}

type updateCartItemRequest struct {
	Quantity int32 `json:"quantity"`
}

func (h *CartHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	itemID := chi.URLParam(r, "itemId")

	var req updateCartItemRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, envelope{"success": false, "error": "invalid request body"})
		return
	}

	cart, err := h.cartSvc.UpdateItem(r.Context(), service.UpdateCartItemInput{
		UserID:     userID,
		CartItemID: itemID,
		Quantity:   req.Quantity,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeSuccess(w, http.StatusOK, cart)
}

func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	itemID := chi.URLParam(r, "itemId")

	cart, err := h.cartSvc.RemoveItem(r.Context(), service.RemoveCartItemInput{
		UserID:     userID,
		CartItemID: itemID,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeSuccess(w, http.StatusOK, cart)
}

func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if err := h.cartSvc.ClearCart(r.Context(), userID); err != nil {
		writeError(w, r, err)
		return
	}
	writeSuccess(w, http.StatusOK, envelope{"message": "cart cleared"})
}
