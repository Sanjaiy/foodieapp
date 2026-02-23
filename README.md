# FoodieApp — Food Ordering API

## Quick Start (Docker)

The fastest way to get everything running:

```bash
# 1. Clone the repository
git clone https://github.com/Sanjaiy/foodieapp.git
cd foodieapp

# 2. Start all services (DB + Preprocessor + API)
docker compose up --build
```

This will:
- Start a **PostgreSQL 16** database on port `5433` (host) → `5432` (container)
- Run the **preprocessor** to generate `data/valid_codes.txt` from the coupon files
- Start the **API server** on `http://localhost:8080`

---

## The Coupon Preprocessor

### What It Does

The preprocessor reads 3 coupon source files (`data/coupon1.txt`, `data/coupon2.txt`, `data/coupon3.txt`) and produces a single sorted output file (`data/valid_codes.txt`) containing only the **valid** coupon codes.

### Validity Criteria

A coupon code is considered **valid** if it meets **both** conditions:
1. **Length**: Between 8 and 10 characters (inclusive)
2. **Frequency**: Appears in **at least 2 out of 3** source files

### How the Logic Works

The preprocessor uses a **bitmask** approach for memory-efficient tracking:

```
File 1 → bit 0 (0b001)
File 2 → bit 1 (0b010)
File 3 → bit 2 (0b100)
```

1. **Scan each file** — For each code with valid length (8–10 chars), set the corresponding bit in a map
2. **Filter** — After all files are processed, keep only codes that appeared in 2+ files
3. **Sort** — The valid codes are sorted alphabetically
4. **Write** — Output to a newline-delimited file (`valid_codes.txt`)

#### Example

| Code       | File 1 | File 2 | File 3 | Bitmask  | Count | Valid? |
|------------|--------|--------|--------|----------|-------|--------|
| `OVER9000` | ✅      | ✅      | ✅      | `0b111`  | 3     | ✅     |
| `GNULINUX` | ✅      | ✅      | ✅      | `0b111`  | 3     | ✅     |
| `SIXTYOFF` | ❌      | ✅      | ❌      | `0b010`  | 1     | ❌     |
| `FREEZAAA...` | ✅   | ❌      | ❌      | `0b001`  | 1     | ❌ (also too long) |

### How to Run the Preprocessor

```bash
# Run via docker compose (output lands in ./data/valid_codes.txt)
docker compose run --rm preprocess
```

### How the Server Uses `valid_codes.txt`

At startup, the API server **memory-maps** (`mmap`) the sorted `valid_codes.txt` file and builds an offset index. When a coupon code is submitted with an order, the server performs a **binary search** over the memory-mapped data — making lookups **O(log n)** with **zero heap allocation** for the file data.

---

## Why File-Based Instead of a Key-Value Store?

In a **production system**, an embedded key-value store (e.g. Redis) would be the ideal choice for storing and looking up valid coupon codes — it offers faster lookups, handles concurrent reads efficiently, and scales to millions of entries without issues.

However, for the purpose of this **assessment**, a file-based approach was chosen to demonstrate a low-level, dependency-free solution:

- **No external dependencies** — uses only Go's standard library (`syscall.Mmap`, `sort.Search`)
- **Simple to reason about** — the entire lookup mechanism is a sorted flat file + binary search

### File Size Scalability

The file-based approach works well for small to moderate datasets. Since the file is memory-mapped (`mmap`), the OS manages paging — only the accessed portions are loaded into RAM. However, practical limits exist:

| File Size | ~Number of Codes | Performance |
|-----------|-------------------|-------------|
| < 1 MB    | ~100K codes        | Excellent — fits in a single OS page batch |
| 1–50 MB   | ~5M codes          | Good — binary search is O(log n), mmap handles paging |

For datasets beyond ~50 MB, a key-value store would be the better choice as it avoids the overhead of scanning line offsets.

---

## Docker In-Depth

### Architecture

```
docker-compose.yml
├── db            — PostgreSQL 16 (Alpine)
├── preprocess    — One-shot container that generates valid_codes.txt
└── api           — The Go API server
```

### Individual Commands

```bash
# Build and start everything
docker compose up --build

# Start only the database
docker compose up db -d

# Run only the preprocessor (output lands in ./data/)
docker compose run --rm preprocess

# Build and start only the API
docker compose up api --build

# Stop all services
docker compose down

# Stop and remove all data (including DB volume)
docker compose down -v
```

### Dockerfile Breakdown

**`Dockerfile`** (API Server) — Multi-stage build:
1. **Builder stage**: Compiles the Go binary
2. **Runtime stage**: Contains only the compiled binary — no shell, no package manager

**`Dockerfile.preprocess`** (Preprocessor):
1. **Builder stage**: Compiles the preprocessor binary
2. **Runtime stage**: Copies the binary and coupon input files, runs to produce `valid_codes.txt`

### Volume Mounts

| Volume | Purpose |
|--------|---------|
| `pgdata` | Persists PostgreSQL data across restarts |
| `./data:/data` (preprocess) | Writes `valid_codes.txt` to the host's `data/` directory |
| `./data:/data:ro` (api) | Mounts the output as read-only for the API server |

---

## API Endpoints & curl Commands

### Health Check

```bash
curl http://localhost:8080/health
```

Response:
```json
{"status":"ok"}
```

### List Products

```bash
curl http://localhost:8080/api/product
```

Response:
```json
[
  {
    "id": "1",
    "name": "Waffle with Berries",
    "price": 6.5,
    "category": "Waffle",
    "image": {
      "thumbnail": "https://orderfoodonline.deno.dev/public/images/image-waffle-thumbnail.jpg",
      "mobile": "https://orderfoodonline.deno.dev/public/images/image-waffle-mobile.jpg",
      "tablet": "https://orderfoodonline.deno.dev/public/images/image-waffle-tablet.jpg",
      "desktop": "https://orderfoodonline.deno.dev/public/images/image-waffle-desktop.jpg"
    }
  }
]
```

### Get Product by ID

```bash
curl http://localhost:8080/api/product/1
```

### Place Order (Without Coupon)

```bash
curl -X POST http://localhost:8080/api/order \
  -H "Content-Type: application/json" \
  -H "api_key: apitest" \
  -d '{
    "items": [
      {"productId": "1", "quantity": 2},
      {"productId": "3", "quantity": 1}
    ]
  }'
```

### Place Order (With Valid Coupon)

```bash
curl -X POST http://localhost:8080/api/order \
  -H "Content-Type: application/json" \
  -H "api_key: apitest" \
  -d '{
    "items": [
      {"productId": "1", "quantity": 2}
    ],
    "couponCode": "OVER9000"
  }'
```

### Place Order (With Invalid Coupon — no discount applied)

```bash
curl -X POST http://localhost:8080/api/order \
  -H "Content-Type: application/json" \
  -H "api_key: apitest" \
  -d '{
    "items": [
      {"productId": "1", "quantity": 2}
    ],
    "couponCode": "INVALIDCODE"
  }'
```

### Unauthorized Request (Missing API Key)

```bash
curl -X POST http://localhost:8080/api/order \
  -H "Content-Type: application/json" \
  -d '{"items": [{"productId": "1", "quantity": 1}]}'
```

Response:
```json
{"code":"unauthorized","message":"api_key header is required"}
```

---

## Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run only handler tests
go test -v ./internal/handler/

# Run only preprocessor tests
go test -v ./cmd/preprocess/
```
