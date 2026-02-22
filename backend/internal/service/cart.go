package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/online-cake-shop/backend/internal/domain"
	"github.com/online-cake-shop/backend/internal/repository/db"
)

type CartService struct {
	q *db.Queries
}

func NewCartService(q *db.Queries) *CartService {
	return &CartService{q: q}
}

// ─── DTOs ────────────────────────────────────────────────────────────────────

type CartItemResponse struct {
	ID           string  `json:"id"`
	ProductID    string  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductImage *string `json:"product_image_url"`
	Price        float64 `json:"price"`
	Quantity     int32   `json:"quantity"`
	Subtotal     float64 `json:"subtotal"`
}

type CartResponse struct {
	ID    string             `json:"id"`
	Items []CartItemResponse `json:"items"`
	Total float64            `json:"total"`
}

type AddCartItemInput struct {
	UserID    uuid.UUID
	ProductID string
	Quantity  int32
}

type UpdateCartItemInput struct {
	UserID     uuid.UUID
	CartItemID string
	Quantity   int32
}

type RemoveCartItemInput struct {
	UserID     uuid.UUID
	CartItemID string
}

// ─── Get Cart ─────────────────────────────────────────────────────────────────

func (s *CartService) GetCart(ctx context.Context, userID uuid.UUID) (*CartResponse, error) {
	cart, err := s.q.GetOrCreateCart(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get or create cart: %w", err)
	}

	items, err := s.q.GetCartItems(ctx, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("get cart items: %w", err)
	}

	return buildCartResponse(cart, items), nil
}

// ─── Add Item ─────────────────────────────────────────────────────────────────

func (s *CartService) AddItem(ctx context.Context, in AddCartItemInput) (*CartResponse, error) {
	if in.Quantity < 1 {
		return nil, &domain.AppError{Err: domain.ErrInvalidInput, Message: "quantity must be at least 1"}
	}

	productID, err := uuid.Parse(in.ProductID)
	if err != nil {
		return nil, &domain.AppError{Err: domain.ErrInvalidInput, Message: "invalid product id"}
	}

	// Verify product exists and has stock
	product, err := s.q.GetProductByID(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get product: %w", err)
	}
	if product.StockQuantity < in.Quantity {
		return nil, domain.ErrInsufficientStock
	}

	cart, err := s.q.GetOrCreateCart(ctx, in.UserID)
	if err != nil {
		return nil, fmt.Errorf("get or create cart: %w", err)
	}

	if _, err := s.q.UpsertCartItem(ctx, db.UpsertCartItemParams{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  in.Quantity,
	}); err != nil {
		return nil, fmt.Errorf("upsert cart item: %w", err)
	}

	return s.GetCart(ctx, in.UserID)
}

// ─── Update Item ──────────────────────────────────────────────────────────────

func (s *CartService) UpdateItem(ctx context.Context, in UpdateCartItemInput) (*CartResponse, error) {
	if in.Quantity < 1 {
		return nil, &domain.AppError{Err: domain.ErrInvalidInput, Message: "quantity must be at least 1"}
	}

	itemID, err := uuid.Parse(in.CartItemID)
	if err != nil {
		return nil, &domain.AppError{Err: domain.ErrInvalidInput, Message: "invalid item id"}
	}

	cart, err := s.q.GetCartByUserID(ctx, in.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get cart: %w", err)
	}

	if _, err := s.q.UpdateCartItemQuantity(ctx, db.UpdateCartItemQuantityParams{
		ID:       itemID,
		Quantity: in.Quantity,
		CartID:   cart.ID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("update cart item: %w", err)
	}

	return s.GetCart(ctx, in.UserID)
}

// ─── Remove Item ──────────────────────────────────────────────────────────────

func (s *CartService) RemoveItem(ctx context.Context, in RemoveCartItemInput) (*CartResponse, error) {
	itemID, err := uuid.Parse(in.CartItemID)
	if err != nil {
		return nil, &domain.AppError{Err: domain.ErrInvalidInput, Message: "invalid item id"}
	}

	cart, err := s.q.GetCartByUserID(ctx, in.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get cart: %w", err)
	}

	if err := s.q.DeleteCartItem(ctx, db.DeleteCartItemParams{
		ID:     itemID,
		CartID: cart.ID,
	}); err != nil {
		return nil, fmt.Errorf("delete cart item: %w", err)
	}

	return s.GetCart(ctx, in.UserID)
}

// ─── Clear Cart ───────────────────────────────────────────────────────────────

func (s *CartService) ClearCart(ctx context.Context, userID uuid.UUID) error {
	cart, err := s.q.GetCartByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil // nothing to clear
		}
		return fmt.Errorf("get cart: %w", err)
	}
	return s.q.ClearCart(ctx, cart.ID)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func buildCartResponse(cart db.Cart, items []db.GetCartItemsRow) *CartResponse {
	resp := &CartResponse{
		ID:    cart.ID.String(),
		Items: make([]CartItemResponse, 0, len(items)),
	}
	var total float64
	for _, item := range items {
		price := numericToFloat(item.ProductPrice)
		subtotal := price * float64(item.Quantity)
		total += subtotal

		ci := CartItemResponse{
			ID:          item.ID.String(),
			ProductID:   item.ProductID.String(),
			ProductName: item.ProductName,
			Price:       price,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		}
		if item.ProductImageUrl.Valid {
			ci.ProductImage = &item.ProductImageUrl.String
		}
		resp.Items = append(resp.Items, ci)
	}
	resp.Total = total
	return resp
}
