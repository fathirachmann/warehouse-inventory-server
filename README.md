# Warehouse Inventory Server

## Description

This is a backend service for a Warehouse Inventory System designed to manage stock, transactions, and product data. Built with **Go (Golang)** using the **Fiber** framework and **PostgreSQL**, it provides a robust RESTful API for inventory operations.

The system handles:

- **Barang (Items)**: Management of product master data with auto-generated codes.
- **Stok (Stock)**: Real-time tracking of inventory levels.
- **Pembelian (Purchases)**: Recording incoming stock from suppliers.
- **Penjualan (Sales)**: Recording outgoing stock to customers.
- **History Stok**: Audit trail for all stock movements.
- **Authentication**: Role-based access control (Admin and Staff) using JWT.

## Project Structure

```
warehouse-inventory-server/
├── config/         # Database configuration
├── docs/           # Swagger documentation files
├── handlers/       # HTTP request handlers (Controllers)
├── middleware/     # Auth, Error Handling, Rate Limiting
├── models/         # Database models (GORM)
├── repositories/   # Data access layer
├── tools/          # Utility tools (Hashing, etc.)
├── utils/          # Helper functions (Validators, etc.)
├── main.go         # Entry point
└── ...
```

## Setup Instruction

### Prerequisites

- **Go (Golang)**: Version 1.23 or higher.
- **PostgreSQL**: Ensure you have PostgreSQL installed and running.
- **Docker & Docker Compose**: (Optional) Recommended for easy setup and deployment.

### Installation & Running in Docker

1. **Clone the repository**

   ```bash
   git clone https://github.com/fathirachmann/warehouse-inventory-server
   cd warehouse-inventory-server
   ```

2. **Environment Configuration**
   The application relies on environment variables. A `.env` file is required.
   Create a `.env` file in the root directory with the following content:

   ```env
   DB_HOST=db
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=warehouse_db
   DB_PORT=5432
   JWT_SECRET=your_secret_key_here
   PORT=8080
   ```

3. **Run with Docker**
   Build and start the application and database containers:

   ```bash
   docker-compose up --build
   ```

   The API will be accessible at `http://localhost:8080`.

4. **Stopping Containers**
   To stop and remove containers:

   ```bash
   docker-compose down
   ```

### Running without Docker

1. **Run in Local Machine**
   Ensure you have Go and PostgreSQL installed. Update `.env` to point to your local database (e.g., `DB_HOST=localhost`).

   ```bash
   go mod tidy
   go run main.go
   ```

2. **Database Setup**
   Create a PostgreSQL database named `warehouse_db` (or as per your `.env` configuration).

3. **Migrations**
   Run the `database-dump.sql` file to set up the database schema if needed, or let GORM auto-migrate (configured in `main.go`).

## API Documentation

### Swagger UI

The API documentation is auto-generated using Swagger. Once the server is running, you can access the interactive documentation at:

**URL**: `http://localhost:8080/swagger/index.html`

### Postman Collection

A comprehensive Postman Collection is included in this repository for testing.

- **File**: `warehouse-auth-seeded.postman_collection.json`
- **Usage**: Import into Postman and use the "Login" requests to set environment variables automatically.

### Registered Users

The system comes with pre-seeded users for testing:

- **Admin**: `admin@warehouse.com` / `Admin123!`
- **Staff 1**: `staff1@warehouse.com` / `Staff1GDA!`
- **Staff 2**: `staff2@warehouse.com` / `Staff2GDB!`

## Error Handling

The API uses a standardized error response format for all endpoints.

**Standard Error:**

```json
{
  "error": "Bad Request",
  "message": "Invalid input data"
}
```

**Validation Error:**

```json
{
  "error": "validation error",
  "message": {
    "email": "Email is required",
    "password": "Password must be at least 6 characters"
  }
}
```

## API Reference (Summary)

For full details, request bodies, and responses, please refer to the **Swagger UI**.

### Authentication

- `POST /api/auth/register` - Register new user (Admin only)
- `POST /api/auth/login` - Login and get JWT

### Barang

- `GET /api/barang` - List all items
- `POST /api/barang` - Create new item
- `GET /api/barang/:id` - Get item details
- `PUT /api/barang/:id` - Update item
- `DELETE /api/barang/:id` - Delete item

### Stok (Stock)

- `GET /api/stok` - List stock for all items
- `GET /api/stok/:barang_id` - Get stock for specific item

### History Stok

- `GET /api/history-stok` - View stock movement history
- `GET /api/history-stok/:barang_id` - View history for specific item

### Transaksi Pembelian

- `GET /api/pembelian` - List purchase transactions
- `POST /api/pembelian` - Create new purchase
- `GET /api/pembelian/:id` - Get purchase details

### Transaksi Penjualan

- `GET /api/penjualan` - List sales transactions
- `POST /api/penjualan` - Create new sale
- `GET /api/penjualan/:id` - Get sale details
