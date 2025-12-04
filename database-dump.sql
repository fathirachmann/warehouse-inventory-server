------------------------------ SCHEMA -----------------------------

-- Use this file to create the database schema for the Warehouse Inventory System

-- Table User
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    full_name VARCHAR(200) NOT NULL,
    role VARCHAR(50) DEFAULT 'staff',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table Master Barang
CREATE TABLE IF NOT EXISTS master_barang (
    id SERIAL PRIMARY KEY,
    kode_barang VARCHAR(50) UNIQUE NOT NULL,
    nama_barang VARCHAR(200) NOT NULL,
    deskripsi TEXT,
    satuan VARCHAR(50) NOT NULL,
    harga_beli DECIMAL(15,2) DEFAULT 0,
    harga_jual DECIMAL(15,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Table Stok
CREATE TABLE IF NOT EXISTS mstok (
    id SERIAL PRIMARY KEY,
    barang_id INTEGER REFERENCES master_barang(id),
    stok_akhir INTEGER DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table History Stok
CREATE TABLE IF NOT EXISTS history_stok (
    id SERIAL PRIMARY KEY,
    barang_id INTEGER REFERENCES master_barang(id),
    user_id INTEGER REFERENCES users(id),
    jenis_transaksi VARCHAR(50) NOT NULL, -- 'masuk', 'keluar', 'adjustment'
    jumlah INTEGER NOT NULL,
    stok_sebelum INTEGER NOT NULL,
    stok_sesudah INTEGER NOT NULL,
    keterangan TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table Pembelian Header
CREATE TABLE IF NOT EXISTS beli_header (
    id SERIAL PRIMARY KEY,
    no_faktur VARCHAR(100) UNIQUE NOT NULL,
    supplier VARCHAR(200) NOT NULL,
    total DECIMAL(15,2) DEFAULT 0,
    user_id INTEGER REFERENCES users(id),
    status VARCHAR(50) DEFAULT 'selesai',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table Pembelian Detail
CREATE TABLE IF NOT EXISTS beli_detail (
    id SERIAL PRIMARY KEY,
    beli_header_id INTEGER REFERENCES beli_header(id),
    barang_id INTEGER REFERENCES master_barang(id),
    qty INTEGER NOT NULL,
    harga DECIMAL(15,2) NOT NULL,
    subtotal DECIMAL(15,2) NOT NULL
);

-- Table Penjualan Header
CREATE TABLE IF NOT EXISTS jual_header (
    id SERIAL PRIMARY KEY,
    no_faktur VARCHAR(100) UNIQUE NOT NULL,
    customer VARCHAR(200) NOT NULL,
    total DECIMAL(15,2) DEFAULT 0,
    user_id INTEGER REFERENCES users(id),
    status VARCHAR(50) DEFAULT 'selesai',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table Penjualan Detail
CREATE TABLE IF NOT EXISTS jual_detail (
    id SERIAL PRIMARY KEY,
    jual_header_id INTEGER REFERENCES jual_header(id),
    barang_id INTEGER REFERENCES master_barang(id),
    qty INTEGER NOT NULL,
    harga DECIMAL(15,2) NOT NULL,
    subtotal DECIMAL(15,2) NOT NULL
);

----------------------------- DATA DUMMY -----------------------------

-- Insert Users: Passwords are bcrypt hashed
INSERT INTO users (username, password, email, full_name, role) VALUES
('admin', '$2a$10$SLvSsMu6kbS5CsZcswlOlOpDDMYVPhOT3hlq15XZGFQe15IoTZOr6', 'admin@warehouse.com', 'Administrator System', 'admin'), -- Password: Admin123!
('staff1', '$2a$10$z7BgTYBk3jonuRV76Gn8jO7OKBkengAZelCHZQj0CzpGJof3srR7G', 'staff1@warehouse.com', 'Staff Gudang A', 'staff'), -- Password: Staff1GDA!
('staff2', '$2a$10$t.57bYH7QMj7i9cKGVBvUOP33pNzDt69knzcYxbYaLDc4qt2eFm56', 'staff2@warehouse.com', 'Staff Gudang B', 'staff'); -- Password: Staff2GDB!

-- Insert Master Barang
INSERT INTO master_barang (kode_barang, nama_barang, deskripsi, satuan, harga_beli, harga_jual) VALUES
('BRG001', 'Laptop Dell XPS 13', 'Laptop Business Grade', 'unit', 15000000, 17500000),
('BRG002', 'Mouse Wireless Logitech', 'Mouse Wireless 2.4GHz', 'pcs', 250000, 350000),
('BRG003', 'Keyboard Mechanical', 'Keyboard Mechanical RGB', 'pcs', 800000, 1200000),
('BRG004', 'Monitor 24 inch', 'Monitor LED 24 inch Full HD', 'unit', 2000000, 2800000),
('BRG005', 'Webcam HD 1080p', 'Webcam High Definition', 'pcs', 450000, 650000);

-- Insert Initial Stock
INSERT INTO mstok (barang_id, stok_akhir) VALUES
(1, 10), (2, 50), (3, 30), (4, 15), (5, 25);

-- Insert Pembelian Data
INSERT INTO beli_header (no_faktur, supplier, total, user_id, status) VALUES
('BLI001', 'PT Supplier Elektronik', 32500000, 2, 'selesai'),
('BLI002', 'CV Komputer Jaya', 12500000, 3, 'selesai');

INSERT INTO beli_detail (beli_header_id, barang_id, qty, harga, subtotal) VALUES
(1, 1, 2, 15000000, 30000000),
(1, 2, 10, 250000, 2500000),
(2, 3, 5, 800000, 4000000),
(2, 4, 3, 2000000, 6000000),
(2, 5, 4, 450000, 1800000);

-- Insert Penjualan Data
INSERT INTO jual_header (no_faktur, customer, total, user_id, status) VALUES
('JUAL001', 'PT Customer Indonesia', 18700000, 2, 'selesai'),
('JUAL002', 'CV Tech Solution', 4150000, 3, 'selesai');

INSERT INTO jual_detail (jual_header_id, barang_id, qty, harga, subtotal) VALUES
(1, 1, 1, 17500000, 17500000),
(1, 2, 2, 350000, 700000),
(1, 3, 1, 1200000, 1200000),
(2, 2, 5, 350000, 1750000),
(2, 4, 1, 2800000, 2800000);

-- Insert History Stok (automatically triggered by transactions)
INSERT INTO history_stok (barang_id, user_id, jenis_transaksi, jumlah, stok_sebelum, stok_sesudah, keterangan) VALUES
(1, 2, 'masuk', 2, 0, 2, 'Pembelian BLI001'),
(2, 2, 'masuk', 10, 0, 10, 'Pembelian BLI001'),
(3, 3, 'masuk', 5, 0, 5, 'Pembelian BLI002'),
(4, 3, 'masuk', 3, 0, 3, 'Pembelian BLI002'),
(5, 3, 'masuk', 4, 0, 4, 'Pembelian BLI002'),
(1, 2, 'keluar', 1, 2, 1, 'Penjualan JUAL001'),
(2, 2, 'keluar', 2, 10, 8, 'Penjualan JUAL001'),
(3, 2, 'keluar', 1, 5, 4, 'Penjualan JUAL001'),
(2, 3, 'keluar', 5, 8, 3, 'Penjualan JUAL002'),
(4, 3, 'keluar', 1, 3, 2, 'Penjualan JUAL002');
