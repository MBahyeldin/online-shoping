package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/online-cake-shop/backend/internal/service"
)

type ProductHandler struct {
	productSvc *service.ProductService
}

func NewProductHandler(productSvc *service.ProductService) *ProductHandler {
	return &ProductHandler{productSvc: productSvc}
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page := queryInt(q.Get("page"), 1)
	limit := queryInt(q.Get("limit"), 20)
	sortBy := q.Get("sort")
	categoryID := q.Get("category_id")

	var catIDPtr *string
	if categoryID != "" {
		catIDPtr = &categoryID
	}

	out, err := h.productSvc.List(r.Context(), service.ListProductsInput{
		CategoryID: catIDPtr,
		SortBy:     sortBy,
		Page:       page,
		Limit:      limit,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeSuccess(w, http.StatusOK, out)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	product, err := h.productSvc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeSuccess(w, http.StatusOK, product)
}

func (h *ProductHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.productSvc.ListCategories(r.Context())
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeSuccess(w, http.StatusOK, cats)
}

func queryInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	var n int
	if _, err := fmt.Sscanf(s, "%d", &n); err != nil || n < 1 {
		return defaultVal
	}
	return n
}
