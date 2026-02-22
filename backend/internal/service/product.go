package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/online-cake-shop/backend/internal/domain"
	"github.com/online-cake-shop/backend/internal/repository/db"
)

type ProductService struct {
	q *db.Queries
}

func NewProductService(q *db.Queries) *ProductService {
	return &ProductService{q: q}
}

// ─── DTOs ────────────────────────────────────────────────────────────────────

type ProductResponse struct {
	ID            string   `json:"id"`
	CategoryID    *string  `json:"category_id"`
	CategoryName  *string  `json:"category_name"`
	CategorySlug  *string  `json:"category_slug"`
	Name          string   `json:"name"`
	Description   *string  `json:"description"`
	Price         float64  `json:"price"`
	ImageURL      *string  `json:"image_url"`
	StockQuantity int32    `json:"stock_quantity"`
	IsActive      bool     `json:"is_active"`
}

type ListProductsInput struct {
	CategoryID *string
	SortBy     string
	Page       int
	Limit      int
}

type ListProductsOutput struct {
	Products   []ProductResponse `json:"products"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// ─── List Products ────────────────────────────────────────────────────────────

func (s *ProductService) List(ctx context.Context, in ListProductsInput) (*ListProductsOutput, error) {
	if in.Page < 1 {
		in.Page = 1
	}
	if in.Limit < 1 || in.Limit > 100 {
		in.Limit = 20
	}

	var catID pgtype.UUID
	if in.CategoryID != nil && *in.CategoryID != "" {
		id, err := uuid.Parse(*in.CategoryID)
		if err != nil {
			return nil, &domain.AppError{Err: domain.ErrInvalidInput, Message: "invalid category_id"}
		}
		catID = pgtype.UUID{Bytes: id, Valid: true}
	}

	sortBy := pgtype.Text{}
	if in.SortBy == "price_asc" || in.SortBy == "price_desc" {
		sortBy = pgtype.Text{String: in.SortBy, Valid: true}
	}

	offset := int32((in.Page - 1) * in.Limit)

	rows, err := s.q.ListProducts(ctx, db.ListProductsParams{
		CategoryID: catID,
		SortBy:     sortBy,
		Limit:      int32(in.Limit),
		Offset:     offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	total, err := s.q.CountProducts(ctx, catID)
	if err != nil {
		return nil, fmt.Errorf("count products: %w", err)
	}

	products := make([]ProductResponse, 0, len(rows))
	for _, r := range rows {
		products = append(products, mapListProductRow(r))
	}

	totalPages := int(total) / in.Limit
	if int(total)%in.Limit > 0 {
		totalPages++
	}

	return &ListProductsOutput{
		Products:   products,
		Total:      total,
		Page:       in.Page,
		Limit:      in.Limit,
		TotalPages: totalPages,
	}, nil
}

// ─── Get Product ─────────────────────────────────────────────────────────────

func (s *ProductService) GetByID(ctx context.Context, id string) (*ProductResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, &domain.AppError{Err: domain.ErrInvalidInput, Message: "invalid product id"}
	}

	row, err := s.q.GetProductByID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get product: %w", err)
	}

	resp := mapGetProductRow(row)
	return &resp, nil
}

// ─── List Categories ─────────────────────────────────────────────────────────

func (s *ProductService) ListCategories(ctx context.Context) ([]db.Category, error) {
	cats, err := s.q.ListCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	return cats, nil
}

// ─── Mappers ─────────────────────────────────────────────────────────────────

func numericToFloat(n pgtype.Numeric) float64 {
	if !n.Valid || n.Int == nil {
		return 0
	}
	// Convert big.Int mantissa to float
	f := new(big.Float).SetInt(n.Int)
	if n.Exp != 0 {
		// 10^exp as a float (handles both positive and negative exponents)
		scale := new(big.Float).SetFloat64(1)
		ten := new(big.Float).SetFloat64(10)
		if n.Exp > 0 {
			for i := int32(0); i < n.Exp; i++ {
				scale.Mul(scale, ten)
			}
			f.Mul(f, scale)
		} else {
			for i := int32(0); i > n.Exp; i-- {
				scale.Mul(scale, ten)
			}
			f.Quo(f, scale)
		}
	}
	result, _ := f.Float64()
	return result
}

func mapListProductRow(r db.ListProductsRow) ProductResponse {
	p := ProductResponse{
		ID:            r.ID.String(),
		Name:          r.Name,
		Price:         numericToFloat(r.Price),
		StockQuantity: r.StockQuantity,
		IsActive:      r.IsActive,
	}
	if r.CategoryID.Valid {
		id := uuid.UUID(r.CategoryID.Bytes).String()
		p.CategoryID = &id
	}
	if r.CategoryName.Valid {
		p.CategoryName = &r.CategoryName.String
	}
	if r.CategorySlug.Valid {
		p.CategorySlug = &r.CategorySlug.String
	}
	if r.Description.Valid {
		p.Description = &r.Description.String
	}
	if r.ImageUrl.Valid {
		p.ImageURL = &r.ImageUrl.String
	}
	return p
}

func mapGetProductRow(r db.GetProductByIDRow) ProductResponse {
	p := ProductResponse{
		ID:            r.ID.String(),
		Name:          r.Name,
		Price:         numericToFloat(r.Price),
		StockQuantity: r.StockQuantity,
		IsActive:      r.IsActive,
	}
	if r.CategoryID.Valid {
		id := uuid.UUID(r.CategoryID.Bytes).String()
		p.CategoryID = &id
	}
	if r.CategoryName.Valid {
		p.CategoryName = &r.CategoryName.String
	}
	if r.CategorySlug.Valid {
		p.CategorySlug = &r.CategorySlug.String
	}
	if r.Description.Valid {
		p.Description = &r.Description.String
	}
	if r.ImageUrl.Valid {
		p.ImageURL = &r.ImageUrl.String
	}
	return p
}
