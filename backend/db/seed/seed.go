// Seed populates the database with sample data for development/demo purposes.
// Run: go run ./db/seed/seed.go
package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/online-cake-shop/backend/internal/config"
	"github.com/online-cake-shop/backend/internal/repository/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.Database.URL())
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("db ping: %v", err)
	}

	q := db.New(pool)
	ctx := context.Background()

	fmt.Println("Seeding categories...")
	cats := []db.CreateCategoryParams{
		{Name: "Birthday Cakes", Slug: "birthday-cakes"},
		{Name: "Wedding Cakes", Slug: "wedding-cakes"},
		{Name: "Custom Cakes", Slug: "custom-cakes"},
		{Name: "Cheesecakes", Slug: "cheesecakes"},
		{Name: "Cupcakes", Slug: "cupcakes"},
	}

	categoryMap := map[string]db.Category{}
	for _, c := range cats {
		cat, err := q.CreateCategory(ctx, c)
		if err != nil {
			log.Printf("  skip category %s (already exists?): %v", c.Name, err)
			cat, _ = q.GetCategoryBySlug(ctx, c.Slug)
		} else {
			fmt.Printf("  created category: %s\n", cat.Name)
		}
		categoryMap[c.Slug] = cat
	}

	fmt.Println("Seeding products...")
	products := []struct {
		CategorySlug  string
		Name          string
		Description   string
		Price         float64
		ImageURL      string
		StockQuantity int32
	}{
		{
			CategorySlug:  "birthday-cakes",
			Name:          "Classic Chocolate Birthday Cake",
			Description:   "Rich triple-layer chocolate cake with smooth ganache frosting and colorful sprinkles. Perfect for any birthday celebration.",
			Price:         45.00,
			ImageURL:      "https://images.unsplash.com/photo-1578985545062-69928b1d9587?w=800",
			StockQuantity: 20,
		},
		{
			CategorySlug:  "birthday-cakes",
			Name:          "Strawberry Dream Birthday Cake",
			Description:   "Light vanilla sponge layered with fresh strawberries and whipped cream frosting. A summer favorite.",
			Price:         42.00,
			ImageURL:      "https://images.unsplash.com/photo-1565958011703-44f9829ba187?w=800",
			StockQuantity: 15,
		},
		{
			CategorySlug:  "wedding-cakes",
			Name:          "Elegant White Wedding Cake",
			Description:   "Five-tier fondant-covered wedding cake with delicate floral decorations and pearl details. Customizable flavors.",
			Price:         350.00,
			ImageURL:      "https://images.unsplash.com/photo-1535254973040-607b474cb50d?w=800",
			StockQuantity: 5,
		},
		{
			CategorySlug:  "wedding-cakes",
			Name:          "Rustic Naked Wedding Cake",
			Description:   "Three-tier semi-naked cake with fresh flowers and berries. A boho-chic option for modern weddings.",
			Price:         280.00,
			ImageURL:      "https://images.unsplash.com/photo-1519671282429-b44b4b0d9b13?w=800",
			StockQuantity: 8,
		},
		{
			CategorySlug:  "custom-cakes",
			Name:          "Princess Castle Cake",
			Description:   "A magical multi-tier castle cake with edible towers, a drawbridge, and sparkle dust. Makes dreams come true.",
			Price:         95.00,
			ImageURL:      "https://images.unsplash.com/photo-1558636508-e0969431f628?w=800",
			StockQuantity: 10,
		},
		{
			CategorySlug:  "custom-cakes",
			Name:          "Galaxy Space Cake",
			Description:   "Stunning galaxy-themed cake with deep blue, purple, and silver swirls. Stars and edible glitter included.",
			Price:         75.00,
			ImageURL:      "https://images.unsplash.com/photo-1563729784474-d77dbb933a9e?w=800",
			StockQuantity: 12,
		},
		{
			CategorySlug:  "cheesecakes",
			Name:          "New York Style Cheesecake",
			Description:   "Classic dense and creamy New York cheesecake with a buttery graham cracker crust. A timeless dessert.",
			Price:         38.00,
			ImageURL:      "https://images.unsplash.com/photo-1533134242443-d4fd215305ad?w=800",
			StockQuantity: 25,
		},
		{
			CategorySlug:  "cheesecakes",
			Name:          "Blueberry Swirl Cheesecake",
			Description:   "Creamy cheesecake with a gorgeous blueberry compote swirl and fresh blueberry topping.",
			Price:         42.00,
			ImageURL:      "https://images.unsplash.com/photo-1571115177098-24ec42ed204d?w=800",
			StockQuantity: 18,
		},
		{
			CategorySlug:  "cupcakes",
			Name:          "Dozen Assorted Cupcakes",
			Description:   "A dozen delicious cupcakes in assorted flavors: chocolate, vanilla, red velvet, and lemon. Perfect for parties.",
			Price:         28.00,
			ImageURL:      "https://images.unsplash.com/photo-1486427944299-d1955d23e34d?w=800",
			StockQuantity: 30,
		},
		{
			CategorySlug:  "cupcakes",
			Name:          "Mini Wedding Cupcake Tower",
			Description:   "An elegant tower of 36 beautifully decorated mini cupcakes, an alternative to a traditional wedding cake.",
			Price:         120.00,
			ImageURL:      "https://images.unsplash.com/photo-1576618148400-f54bed99fcfd?w=800",
			StockQuantity: 10,
		},
	}

	for _, p := range products {
		cat := categoryMap[p.CategorySlug]

		// Build pgtype values
		catID := pgtype.UUID{Bytes: cat.ID, Valid: cat.ID != [16]byte{}}
		desc := pgtype.Text{String: p.Description, Valid: true}
		imgURL := pgtype.Text{String: p.ImageURL, Valid: true}

		// Price as NUMERIC
		scaledPrice := int64(p.Price * 100)
		priceNumeric := pgtype.Numeric{Int: big.NewInt(scaledPrice), Exp: -2, Valid: true}

		product, err := q.CreateProduct(ctx, db.CreateProductParams{
			CategoryID:    catID,
			Name:          p.Name,
			Description:   desc,
			Price:         priceNumeric,
			ImageUrl:      imgURL,
			StockQuantity: p.StockQuantity,
		})
		if err != nil {
			log.Printf("  skip product %s: %v", p.Name, err)
		} else {
			fmt.Printf("  created product: %s ($%.2f)\n", product.Name, p.Price)
		}
	}

	fmt.Println("\nSeed complete!")
	os.Exit(0)
}
