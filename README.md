# Pizza Backend

REST API for a pizza shop built with Go (Gin + PostgreSQL).

## Stack

- **Go 1.25** — language
- **Gin** — HTTP framework
- **PostgreSQL 17** — database
- **sqlx** — ergonomic wrapper over `database/sql`
- **golang-migrate** — migrations
- **Docker / Docker Compose** — environment

## Project Structure

```
.
├── main.go                  # Entry point, dependency wiring
├── router.go                # Routes
├── migrations/              # SQL migrations (up / down)
└── internal/
    ├── config/              # Config from env variables
    ├── database/            # DB connection and migration runner
    ├── model/               # Domain models
    ├── repository/          # Database layer
    ├── service/             # Business logic
    ├── handler/             # HTTP handlers
    ├── api/                 # Request / Response types
    ├── apperror/            # Error codes
    └── response/            # Response helpers
```

## Running

```bash
make run
```

Builds the image and starts two containers: `pizza-postgres` (port `5445`) and `fp_pizza` (port `3002`). Migrations are applied automatically on startup.

### Environment Variables

| Variable              | Description               | Example                                             |
|-----------------------|---------------------------|-----------------------------------------------------|
| `PIZZA_DATABASE_URL`  | PostgreSQL connection URL | `postgres://user:pass@host:5432/db?sslmode=disable` |
| `PIZZA_ADDR`          | Server address            | `:3002` (default)                                   |

Copy `.env` and fill in your values for local development.

## API

Base prefix: `/v1`

### Categories

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/categories` | List categories |
| `GET` | `/categories/:id/addons` | List addons for a category |

### Products

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/products` | List products |
| `GET` | `/products/:id` | Product details with variants |

Query parameters for `/products`: `categoryId`, `q` (search), `page`, `perPage`, `offset`, `limit`, `sortOrder` (`ASC` / `DESC`).  
Response includes `X-Total-Count` and `X-Total-Pages` headers.

### Cart

The cart is identified by a UUID returned on creation. The client is responsible for storing it.

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/cart` | Create a cart |
| `GET` | `/cart/:id` | Get cart with items |
| `DELETE` | `/cart/:id` | Clear cart |
| `POST` | `/cart/:id/items` | Add an item |
| `PATCH` | `/cart/:id/items/:itemId` | Update item quantity |
| `DELETE` | `/cart/:id/items/:itemId` | Remove an item |

**POST /cart/\:id/items**
```json
{
  "product_id": 1,
  "variant_id": 2,
  "quantity": 1,
  "addon_ids": [3, 4]
}
```

**PATCH /cart/\:id/items/\:itemId**
```json
{
  "quantity": 3
}
```

### Response Format

Success:
```json
{
  "success": true,
  "data": { ... },
  "message": "..."
}
```

Error:
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "cart not found"
  }
}
```

## Migrations

Migrations are applied automatically on server startup. Files are located in `migrations/` and follow the `{version}_{name}.up.sql` / `{version}_{name}.down.sql` naming convention.
