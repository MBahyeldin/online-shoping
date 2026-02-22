package service_test

import (
	"math/big"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/online-cake-shop/backend/internal/service"
)

func TestNumericToFloat(t *testing.T) {
	tests := []struct {
		name    string
		numeric pgtype.Numeric
		want    float64
	}{
		{
			name:    "zero",
			numeric: pgtype.Numeric{Int: big.NewInt(0), Exp: 0, Valid: true},
			want:    0,
		},
		{
			name:    "integer 10",
			numeric: pgtype.Numeric{Int: big.NewInt(10), Exp: 0, Valid: true},
			want:    10,
		},
		{
			name:    "decimal 9.99 (999 * 10^-2)",
			numeric: pgtype.Numeric{Int: big.NewInt(999), Exp: -2, Valid: true},
			want:    9.99,
		},
		{
			name:    "invalid numeric returns 0",
			numeric: pgtype.Numeric{Valid: false},
			want:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.NumericToFloat(tt.numeric)
			// Allow tiny floating-point delta
			diff := got - tt.want
			if diff < -0.0001 || diff > 0.0001 {
				t.Errorf("NumericToFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}
