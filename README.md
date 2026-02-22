# 🎂 Cake Shop — Online Cake Ordering Application

A full-stack, production-ready online cake ordering platform built with **Go** (backend) and **React + TypeScript** (frontend).

---

## Tech Stack

| Layer      | Technology                                              |
|------------|---------------------------------------------------------|
| Frontend   | React 18, Vite, TypeScript, TailwindCSS, shadcn/ui      |
| State      | TanStack Query (server), Zustand (client)               |
| Routing    | React Router v7                                         |
| Validation | Zod + react-hook-form                                   |
| Backend    | Go 1.23, chi router, Clean Architecture                 |
| Auth       | JWT (HTTP-only cookie + Bearer token)                   |
| Email      | SMTP / Mock (pluggable interface)                       |
| Database   | PostgreSQL 16                                           |
| ORM/Query  | sqlc (type-safe SQL code generation)                    |
| Migrations | golang-migrate                                          |
| Container  | Docker + Docker Compose                                 |

---

## Project Structure

```
online-shoping/
├── backend/
│   ├── cmd/api/main.go              # Application entry point
│   ├── internal/
│   │   ├── config/config.go         # Environment-based configuration
│   │   ├── domain/errors.go         # Sentinel errors
│   │   ├── handler/                 # HTTP handlers (auth, product, cart, order)
│   │   ├── middleware/              # Auth middleware, structured logger
│   │   ├── service/                 # Business logic layer
│   │   ├── repository/db/           # sqlc-generated repository layer
│   │   └── email/                   # Email sender interface + SMTP/Mock impls
│   ├── db/
│   │   ├── migrations/              # golang-migrate SQL files
│   │   ├── queries/                 # sqlc SQL query files
│   │   └── seed/seed.go             # Database seed script
│   ├── sqlc.yaml                    # sqlc configuration
│   ├── Dockerfile
│   └── .env.example
├── frontend/
│   └── src/
│       ├── components/
│       │   ├── ui/                  # shadcn/ui components (Button, Input, etc.)
│       │   ├── layout/              # Header, Footer
│       │   ├── home/                # Hero, FeaturedCakes, Testimonials
│       │   └── shared/              # ProductCard, CartDrawer, LoadingSpinner
│       ├── pages/                   # HomePage, RegisterPage, VerifyOTPPage, ProductsPage, CartPage
│       ├── services/                # API service functions (auth, products, cart, orders)
│       ├── store/                   # Zustand stores (auth, cart)
│       ├── types/                   # TypeScript type definitions
│       └── lib/                     # Utilities (cn, formatCurrency, etc.)
├── docker-compose.yml
├── Makefile
├── openapi.yaml                     # OpenAPI 3.0 specification
└── README.md
```

---

## Prerequisites

- [Go 1.23+](https://go.dev/dl/)
- [Node.js 20+](https://nodejs.org/) and npm
- [Docker + Docker Compose](https://docs.docker.com/compose/)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) *(for running migrations manually)*
- [sqlc](https://sqlc.dev/) *(optional — only needed if modifying queries)*

---

## Quick Start (Recommended)

### 1. Clone and configure

```bash
git clone <repo-url>
cd online-shoping

# Copy backend environment config
cp backend/.env.example backend/.env
# Edit backend/.env to set secrets if needed (defaults work for local dev)
```

### 2. Start PostgreSQL with Docker

```bash
make docker-db
# Wait a few seconds for Postgres to be ready
```

### 3. Run database migrations

```bash
make migrate-up
```

### 4. Seed the database

```bash
make seed
```

### 5. Start the backend

```bash
make dev
# Server starts on http://localhost:8080
```

### 6. Install frontend dependencies and start the dev server

```bash
make deps-frontend
make frontend-dev
# Frontend starts on http://localhost:5173
```

Open [http://localhost:5173](http://localhost:5173) in your browser.

---

## Docker Compose (Full Stack)

Start the entire stack (PostgreSQL + API) with a single command:

```bash
docker compose up -d --build

# Check logs
make docker-logs

# Stop everything
make docker-down
```

> The frontend is not containerized in this setup — run it locally with `make frontend-dev`.

---

## Environment Variables (Backend)

Copy `backend/.env.example` to `backend/.env` and configure:

| Variable              | Default                                | Description                         |
|-----------------------|----------------------------------------|-------------------------------------|
| `SERVER_PORT`         | `8080`                                 | HTTP server port                    |
| `ENV`                 | `development`                          | Environment name                    |
| `ALLOWED_ORIGINS`     | `http://localhost:5173`               | CORS allowed origins (comma-sep)    |
| `DB_HOST`             | `localhost`                            | PostgreSQL host                     |
| `DB_PORT`             | `5432`                                 | PostgreSQL port                     |
| `DB_NAME`             | `cake_shop`                            | Database name                       |
| `DB_USER`             | `postgres`                             | Database user                       |
| `DB_PASSWORD`         | `postgres`                             | Database password                   |
| `JWT_SECRET`          | *(change this!)*                       | HS256 signing secret (min 32 chars) |
| `JWT_ACCESS_TOKEN_TTL`| `24h`                                  | Token expiry duration               |
| `EMAIL_PROVIDER`      | `mock`                                 | `mock` or `smtp`                    |
| `EMAIL_FROM`          | `noreply@cakeshop.com`                | Sender email address                |
| `SMTP_HOST`           | *(empty)*                              | SMTP server host                    |
| `SMTP_PORT`           | `587`                                  | SMTP server port                    |
| `SMTP_USER`           | *(empty)*                              | SMTP username                       |
| `SMTP_PASS`           | *(empty)*                              | SMTP password / app password        |

> When `EMAIL_PROVIDER=mock`, OTPs are printed to the server console — perfect for development.

---

## API Documentation

See [openapi.yaml](./openapi.yaml) for the full OpenAPI 3.0 specification.

You can view it in Swagger UI:
```bash
npx @redocly/cli preview-docs openapi.yaml
```

### Key Endpoints

| Method | Path                    | Auth | Description                          |
|--------|-------------------------|------|--------------------------------------|
| POST   | `/api/v1/auth/register` | —    | Register and trigger OTP email       |
| POST   | `/api/v1/auth/verify-otp` | —  | Verify OTP and receive JWT           |
| POST   | `/api/v1/auth/resend-otp` | —  | Resend OTP (rate-limited: 3/hour)    |
| GET    | `/api/v1/products`      | —    | List products (filter, sort, paginate)|
| GET    | `/api/v1/products/:id`  | —    | Get single product                   |
| GET    | `/api/v1/categories`    | —    | List categories                      |
| GET    | `/api/v1/cart`          | ✓    | Get cart                             |
| POST   | `/api/v1/cart/items`    | ✓    | Add item to cart                     |
| PUT    | `/api/v1/cart/items/:id`| ✓    | Update item quantity                 |
| DELETE | `/api/v1/cart/items/:id`| ✓    | Remove item                          |
| DELETE | `/api/v1/cart`          | ✓    | Clear cart                           |
| POST   | `/api/v1/orders`        | ✓    | Create order (transactional)         |
| GET    | `/api/v1/orders`        | ✓    | List user orders                     |
| GET    | `/api/v1/orders/:id`    | ✓    | Get specific order                   |

---

## Running Tests

```bash
make test
# Or with coverage
make test-coverage
```

---

## Database Migrations

```bash
# Apply all pending migrations
make migrate-up

# Roll back last migration
make migrate-down

# Show current version
make migrate-status
```

---

## Regenerating sqlc Code

If you modify SQL queries in `backend/db/queries/`, regenerate the Go code:

```bash
make sqlc-generate
```

---

## Makefile Reference

```bash
make help           # Show all available targets
make setup          # Full first-time setup
make dev            # Run backend locally
make frontend-dev   # Run frontend dev server
make build          # Build backend binary
make test           # Run tests
make migrate-up     # Apply migrations
make seed           # Seed database
make docker-up      # Start all Docker services
make docker-down    # Stop Docker services
make docker-logs    # Tail API logs
```

---

## Security Notes

- OTPs are hashed with bcrypt before storage
- OTPs expire in 5 minutes
- Max 3 OTP requests per hour per user
- Max 5 OTP verification attempts before lockout
- JWT signed with HS256; stored in HTTP-only cookie + `Authorization` header
- All inputs validated on both frontend (Zod) and backend
- SQL injection prevented by parameterized queries (sqlc/pgx)
- CORS configured to only allow specified origins

---

## License

MIT
# online-shoping
