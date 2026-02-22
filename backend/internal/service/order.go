package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/online-cake-shop/backend/internal/domain"
	"github.com/online-cake-shop/backend/internal/repository/db"
)

type OrderService struct {
	pool *pgxpool.Pool
	q    *db.Queries
}

func NewOrderService(pool *pgxpool.Pool, q *db.Queries) *OrderService {
	return &OrderService{pool: pool, q: q}
}

// ─── DTOs ────────────────────────────────────────────────────────────────────

type CreateOrderInput struct {
	UserID          uuid.UUID
	DeliveryAddress string
	DeliveryDate    time.Time
	Notes           string
	PaymentMethod   string
}

type OrderItemResponse struct {
	ID          string  `json:"id"`
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	ImageURL    *string `json:"image_url"`
	Quantity    int32   `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
}

type OrderResponse struct {
	ID              string              `json:"id"`
	DeliveryAddress string              `json:"delivery_address"`
	DeliveryDate    time.Time           `json:"delivery_date"`
	Notes           *string             `json:"notes"`
	PaymentMethod   string              `json:"payment_method"`
	Status          string              `json:"status"`
	TotalAmount     float64             `json:"total_amount"`
	Items           []OrderItemResponse `json:"items"`
	CreatedAt       time.Time           `json:"created_at"`
}

type ListOrdersOutput struct {
	Orders     []OrderResponse `json:"orders"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

// ─── Create Order (transactional) ────────────────────────────────────────────

func (s *OrderService) CreateOrder(ctx context.Context, in CreateOrderInput) (*OrderResponse, error) {
	if err := validateOrderInput(in); err != nil {
		return nil, err
	}

	// Load cart items outside the transaction first
	cart, err := s.q.GetCartByUserID(ctx, in.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEmptyCart
		}
		return nil, fmt.Errorf("get cart: %w", err)
	}

	cartItems, err := s.q.GetCartItems(ctx, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("get cart items: %w", err)
	}
	if len(cartItems) == 0 {
		return nil, domain.ErrEmptyCart
	}

	// Collect product IDs to lock for stock check
	productIDs := make([]uuid.UUID, 0, len(cartItems))
	for _, ci := range cartItems {
		productIDs = append(productIDs, ci.ProductID)
	}

	var order db.Order
	var orderItems []db.GetOrderItemsRow

	err = pgx.BeginTxFunc(ctx, s.pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		qtx := s.q.WithTx(tx)

		// Fetch and lock products
		products, err := qtx.GetProductsForOrder(ctx, productIDs)
		if err != nil {
			return fmt.Errorf("fetch products: %w", err)
		}

		productMap := make(map[uuid.UUID]db.GetProductsForOrderRow, len(products))
		for _, p := range products {
			productMap[p.ID] = p
		}

		// Validate stock
		for _, ci := range cartItems {
			p, ok := productMap[ci.ProductID]
			if !ok {
				return fmt.Errorf("product %s not found", ci.ProductID)
			}
			if p.StockQuantity < ci.Quantity {
				return &domain.AppError{
					Err:     domain.ErrInsufficientStock,
					Message: fmt.Sprintf("not enough stock for '%s'", p.Name),
				}
			}
		}

		// Compute total
		var totalAmount float64
		for _, ci := range cartItems {
			p := productMap[ci.ProductID]
			price := numericToFloat(p.Price)
			totalAmount += price * float64(ci.Quantity)
		}

		// Create order
		notes := pgtype.Text{}
		if in.Notes != "" {
			notes = pgtype.Text{String: in.Notes, Valid: true}
		}

		totalNumeric, err := floatToNumeric(totalAmount)
		if err != nil {
			return fmt.Errorf("convert total: %w", err)
		}

		paymentMethod := in.PaymentMethod
		if paymentMethod == "" {
			paymentMethod = "cash_on_delivery"
		}

		order, err = qtx.CreateOrder(ctx, db.CreateOrderParams{
			UserID:          in.UserID,
			DeliveryAddress: in.DeliveryAddress,
			DeliveryDate:    in.DeliveryDate,
			Notes:           notes,
			PaymentMethod:   paymentMethod,
			TotalAmount:     totalNumeric,
		})
		if err != nil {
			return fmt.Errorf("create order: %w", err)
		}

		// Create order items and deduct stock
		for _, ci := range cartItems {
			p := productMap[ci.ProductID]
			unitPrice := numericToFloat(p.Price)
			totalPrice := unitPrice * float64(ci.Quantity)

			unitPriceNumeric, err := floatToNumeric(unitPrice)
			if err != nil {
				return err
			}
			totalPriceNumeric, err := floatToNumeric(totalPrice)
			if err != nil {
				return err
			}

			if _, err := qtx.CreateOrderItem(ctx, db.CreateOrderItemParams{
				OrderID:    order.ID,
				ProductID:  ci.ProductID,
				Quantity:   ci.Quantity,
				UnitPrice:  unitPriceNumeric,
				TotalPrice: totalPriceNumeric,
			}); err != nil {
				return fmt.Errorf("create order item: %w", err)
			}

			// Deduct stock
			if err := qtx.DeductProductStock(ctx, db.DeductProductStockParams{
				ID:       ci.ProductID,
				Quantity: ci.Quantity,
			}); err != nil {
				return fmt.Errorf("deduct stock: %w", err)
			}
		}

		// Clear cart
		if err := qtx.ClearCart(ctx, cart.ID); err != nil {
			return fmt.Errorf("clear cart: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Load order items for response (outside transaction)
	orderItems, err = s.q.GetOrderItems(ctx, order.ID)
	if err != nil {
		return nil, fmt.Errorf("get order items: %w", err)
	}

	return mapOrderResponse(order, orderItems), nil
}

// ─── List Orders ──────────────────────────────────────────────────────────────

func (s *OrderService) ListOrders(ctx context.Context, userID uuid.UUID, page, limit int) (*ListOrdersOutput, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	offset := int32((page - 1) * limit)
	orders, err := s.q.ListOrdersByUserID(ctx, userID, int32(limit), offset)
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}

	total, err := s.q.CountOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("count orders: %w", err)
	}

	responses := make([]OrderResponse, 0, len(orders))
	for _, o := range orders {
		items, err := s.q.GetOrderItems(ctx, o.ID)
		if err != nil {
			return nil, fmt.Errorf("get order items: %w", err)
		}
		responses = append(responses, *mapOrderResponse(o, items))
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &ListOrdersOutput{
		Orders:     responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// ─── Get Order ────────────────────────────────────────────────────────────────

func (s *OrderService) GetOrder(ctx context.Context, userID uuid.UUID, orderID string) (*OrderResponse, error) {
	oid, err := uuid.Parse(orderID)
	if err != nil {
		return nil, &domain.AppError{Err: domain.ErrInvalidInput, Message: "invalid order id"}
	}

	order, err := s.q.GetOrderByID(ctx, oid, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get order: %w", err)
	}

	items, err := s.q.GetOrderItems(ctx, order.ID)
	if err != nil {
		return nil, fmt.Errorf("get order items: %w", err)
	}

	return mapOrderResponse(order, items), nil
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func validateOrderInput(in CreateOrderInput) error {
	if in.DeliveryAddress == "" {
		return &domain.AppError{Err: domain.ErrInvalidInput, Message: "delivery address is required"}
	}
	if in.DeliveryDate.IsZero() || in.DeliveryDate.Before(time.Now()) {
		return &domain.AppError{Err: domain.ErrInvalidInput, Message: "delivery date must be in the future"}
	}
	return nil
}

func floatToNumeric(f float64) (pgtype.Numeric, error) {
	bf := new(big.Float).SetFloat64(f)
	// Use 2 decimal places
	scaled := new(big.Float).Mul(bf, new(big.Float).SetFloat64(100))
	intVal, _ := new(big.Int).SetString(scaled.Text('f', 0), 10)
	return pgtype.Numeric{Int: intVal, Exp: -2, Valid: true}, nil
}

func mapOrderResponse(o db.Order, items []db.GetOrderItemsRow) *OrderResponse {
	resp := &OrderResponse{
		ID:              o.ID.String(),
		DeliveryAddress: o.DeliveryAddress,
		DeliveryDate:    o.DeliveryDate,
		PaymentMethod:   o.PaymentMethod,
		Status:          o.Status,
		TotalAmount:     numericToFloat(o.TotalAmount),
		Items:           make([]OrderItemResponse, 0, len(items)),
		CreatedAt:       o.CreatedAt,
	}
	if o.Notes.Valid {
		resp.Notes = &o.Notes.String
	}

	for _, item := range items {
		oi := OrderItemResponse{
			ID:          item.ID.String(),
			ProductID:   item.ProductID.String(),
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   numericToFloat(item.UnitPrice),
			TotalPrice:  numericToFloat(item.TotalPrice),
		}
		if item.ProductImageUrl.Valid {
			oi.ImageURL = &item.ProductImageUrl.String
		}
		resp.Items = append(resp.Items, oi)
	}
	return resp
}
