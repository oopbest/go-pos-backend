# ☕ Brew & Bite POS - Backend API (Go + Fiber + PostgreSQL)

A high-performance, lightweight, and robust RESTful API server built for **Restaurant & Cafe Point of Sale (POS)** and **Kitchen Display System (KDS)**.

Part of the **Zero-to-Hero Go POS** series, designed to serve the React frontend [go-pos-frontend](https://github.com/oopbest/go-pos-frontend).

---

## ✨ Features

- 🪑 **Dining Table Lifecycle**:
  - Track dining table states (`available`, `occupied`, `reserved`).
  - Zone management (`Indoor`, `Outdoor`) with seat capacities.
- 📖 **Menu & Station Routing**:
  - Categorized menu items (Main dishes, Drinks & Coffee, Desserts).
  - Station-specific routing: automatically tags items for **Kitchen (ครัว)** or **Bar (บาร์น้ำ)**.
- 📝 **Order Lifecycle & Transactions**:
  - Open table and place initial orders with ACID database transactions (`tx.Begin()`, `tx.Commit()`, `tx.Rollback()`).
  - Add extra items to open orders seamlessly.
  - Automatic table state switching: available ➔ occupied upon order placement.
- 🍳 **Kitchen Display System (KDS) API**:
  - Real-time ticket tracking for kitchen & bar staff.
  - Granular dish status transitions: `pending` ➔ `cooking` ➔ `ready` ➔ `served`.
- 💳 **Cashier Checkout & Billing**:
  - Itemized bill review with discount calculation.
  - Multi-payment support: `cash` (with change calculation), `promptpay` QR, and `credit_card`.
  - Automatically completes order and frees table back to `available`.
- 📑 **Interactive Swagger / OpenAPI UI**:
  - Self-documenting API accessible at `/swagger/` with instant "Try it out" capabilities.
- 💾 **Auto-Migration & Seed Data**:
  - Database schema automatically managed by GORM.
  - Idempotent seed data automatically populates default tables and sample menus on first boot.

---

## 🛠️ Tech Stack

- **Language**: [Go (Golang)](https://go.dev/)
- **Web Framework**: [Fiber v2](https://gofiber.io/) (High-performance Express-inspired framework)
- **Database**: [PostgreSQL 16](https://www.postgresql.org/) (Containerized via Docker Compose)
- **ORM**: [GORM](https://gorm.io/) (Postgres driver, auto-migration, relational preloading)
- **API Documentation**: [Swagger / Swag CLI](https://github.com/swaggo/swag)
- **Containerization**: [Docker & Docker Compose](https://www.docker.com/)

---

## 📂 Project Architecture

```text
backend/
├── cmd/
│   └── api/
│       └── main.go           # Entry point: Server initialization & route grouping
├── internal/
│   ├── config/
│   │   └── config.go         # Environment variables & configuration loader
│   ├── database/
│   │   └── db.go             # PostgreSQL connection pool, auto-migration & seeder
│   ├── models/
│   │   ├── table.go          # Dining table model & status enums
│   │   ├── category.go       # Category model & default station enum
│   │   ├── product.go        # Product model with station tagging
│   │   └── order.go          # Order & OrderItem models with kitchen status
│   └── handlers/
│       ├── table_handler.go   # Table HTTP controllers & swagger annotations
│       ├── product_handler.go # Menu & category HTTP controllers
│       ├── order_handler.go   # Order creation, add items & checkout controllers
│       └── kitchen_handler.go # KDS items listing & status update controllers
├── docs/                     # Swagger OpenAPI generated docs (docs.go, swagger.json, swagger.yaml)
├── docker-compose.yml        # PostgreSQL 16 container definition
├── go.mod
└── go.sum
```

---

## 🌐 API Endpoints Summary

| Method | Endpoint | Description | Tag |
| :--- | :--- | :--- | :--- |
| `GET` | `/health` | Server health check | System |
| `GET` | `/swagger/*` | Swagger UI Interactive Documentation | Docs |
| `GET` | `/api/tables` | Get all dining tables | Tables |
| `GET` | `/api/tables/:id` | Get single table by ID | Tables |
| `PUT` | `/api/tables/:id/status` | Update table status (`available`, `occupied`, etc.) | Tables |
| `GET` | `/api/categories` | List all menu categories | Menu |
| `GET` | `/api/products` | List menu products (supports `?category_id=` & `?station=`) | Menu |
| `POST` | `/api/orders` | Open table & place initial order | Orders |
| `GET` | `/api/orders/table/:table_id` | Get active open order for table | Orders |
| `POST` | `/api/orders/:id/items` | Add extra items to open order | Orders |
| `POST` | `/api/orders/:id/checkout` | Complete checkout, close order & free table | Orders |
| `GET` | `/api/kitchen/items` | List pending kitchen/bar dishes (supports `?station=`) | Kitchen |
| `PUT` | `/api/kitchen/items/:id/status` | Update dish status (`cooking`, `ready`, `served`) | Kitchen |

---

## 🚀 Getting Started

### Prerequisites
- [Go](https://go.dev/dl/) (v1.22 or higher)
- [Docker](https://www.docker.com/) & Docker Compose
- (Optional) [Swag CLI](https://github.com/swaggo/swag): `go install github.com/swaggo/swag/cmd/swag@latest`

### 1. Clone the repository
```bash
git clone https://github.com/oopbest/go-pos-backend.git
cd go-pos-backend
```

### 2. Start PostgreSQL Database
```bash
docker compose up -d
```
*(PostgreSQL will boot up on port `5432` with user `posuser` and database `pos_db`)*

### 3. Run the Go Server
```bash
go run cmd/api/main.go
```

The server will automatically:
1. Connect to PostgreSQL
2. Run Auto-Migration for all models
3. Seed initial dining tables and sample menu items
4. Listen on `http://localhost:8080`

### 4. Interactive Swagger Documentation
Open your browser and visit:
👉 **`http://localhost:8080/swagger/`**

---

## 🔄 Regenerating Swagger Docs

Whenever you modify handler annotations or data models, regenerate Swagger definitions by running:

```bash
swag init -g cmd/api/main.go
```

---

## 🔗 Related Repositories

- **React Frontend**: [oopbest/go-pos-frontend](https://github.com/oopbest/go-pos-frontend) (React 19 + TypeScript + Vite + TailwindCSS + Lucide)
