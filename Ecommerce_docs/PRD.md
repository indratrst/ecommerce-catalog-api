## PRD
1. Scope Aplikasi & Target Pengguna
Apakah aplikasi ini nantinya hanya berupa katalog (pengguna melihat produk, filter, dan jika ingin beli diarahkannya ke WhatsApp/Marketplace lain)?

# Kita akan membuat ecommercenya langsung dengan menggunakan golang dan juga nextJs 

Atau mencakup checkout lengkap (keranjang belanja, perhitungan ongkir, dan payment gateway)?

# iya betul mencakup checkout lengkap (keranjang belanja, perhitungan ongkir, dan payment gateway)

2. Pengelolaan Produk & Stok
Selain varian Warna dan Ukuran, apakah ada atribut khusus fashion lain yang perlu dicatat? (Contoh: Bahan/Material, Chart Size, Brand/Merek).

# Cukup size,warna(optional) dan stock yang lainnya akan fleksibel mengikuti kebutuhan

Apakah satu produk bisa memiliki banyak gambar/galeri per varian warna, atau gambar produk disatukan di tingkat produk utama saja?

# iya satu produk bisa memiliki beberapa gambar 

3. Fitur Utama di Frontend (Next.js)
Pencarian & Filter: Filter apa saja yang paling krusial untuk calon pembeli? (Misal: filter harga, rentang ukuran, warna, atau kategori bertingkat/sub-kategori).

# filter rentang harga,produk,dan kategori yang saat ini dibutuhkan mungkin yang lainnya bisa improve di kemudian hari  


Apakah ada fitur Wishlist, Ulasan/Rating Produk, atau Promo/Kupon Diskon di rilis pertama ini (V1)?

# iya saya ingin ada ketiga fiturnya

4. Admin Panel & Manajemen (Go Fiber Backend)
Bagaimana pengelola toko akan memasukkan/mengelola data produk? Apakah perlu fitur Admin Dashboard lengkap (CRUD Produk, Manajemen Stok, Kategori) di dalam sistem ini?

# Jelas sangat dibutuhkan

5. Skala V1 (Minimum Viable Product)
Apakah ada fitur spesifik lain yang wajib ada di V1 selain penjelajahan katalog produk?

# saya rasa sudah cukup untuk v1 






------------------------------------------------------------------------------------------------
## 1. Overview & Goal
Membangun platform e-commerce fashion mandiri (full checkout) dengan pengalaman pengguna (user experience) modern. Sistem ini memisahkan layer presentasi di Next.js dan layer backend high-performance menggunakan Go (Fiber) dengan penyimpan data PostgreSQL.

Key Success Metrics (MVP)
Waktu pemuatan halaman (Page Load Speed) katalog & detail produk di bawah 1.5 detik.

Keberhasilan aliran transaksi checkout & integrasi payment gateway.

Manajemen stok varian produk yang presisi tanpa ada isu overselling.

## 2. Target User & Roles
Customer (Pembeli):

Menjelajahi katalog produk fashion dengan pemfilteran yang cepat.

Mengelola keranjang belanja, wishlist, dan memberikan ulasan produk.

Melakukan pembayaran secara digital dan memantau status pesanan.

Admin / Store Manager:

Mengelola inventaris produk (Kategori, Varian Warna, Ukuran, Gambar).

Mengelola kupon promo, melihat ulasan, serta memperbarui status pengiriman/pesanan.

## 3. Scope & Feature Requirements (V1 Scope)
## 3.1 Consumer-Facing Features (Next.js Frontend)
A. Auth & User Management
Registrasi & Login pengguna (Email/Password, JWT based).

Pengelolaan Profil Pengguna (Alamat Pengiriman, Kontak, Password).

B. Catalog & Product Discovery
Home Page: Banner promo, Produk Terpopuler, Produk Terbaru.

Product Listing Page (PLP):

Grid tampilan produk.

Filter Utama: Rentang Harga (Min-Max Price), Kategori (termasuk Sub-kategori), Search Bar Nama Produk.

Sorting: Harga Terendah/Tertinggi, Produk Terbaru.

Product Detail Page (PDP):

Galeri Multi-Gambar per produk.

Pemilihan Varian: Ukuran (S, M, L, XL, dll.) & Warna (opsional).

Indikator Stok real-time berdasarkan varian yang dipilih.

Bagian Ulasan/Rating Produk dari pembeli lain.

Tombol Add to Cart & Add to Wishlist.

C. Engagement & Incentives
Wishlist: Menyimpan produk favorit ke dalam daftar Wishlist akun pengguna.

Product Reviews & Ratings: Pembeli dapat memberi rating (1–5 bintang) dan teks ulasan setelah transaksi selesai.

Promo / Kupon Diskon: Input kode voucher promo di halaman keranjang/checkout untuk mendapatkan potongan harga.

D. Cart & Checkout Flow
Cart (Keranjang Belanja):

Tambah/Edit/Hapus item per varian produk.

Kalkulasi otomatis subtotal harga.

Checkout:

### Type Courier 
Pemilihan Alamat Pengiriman.

Kalkulasi Ongkos Kirim (Integrasi API Ekspedisi, misal: RajaOngkir/Biteship).

Input Kode Kupon Promo.

Integrasi Payment Gateway (Misal: Midtrans / Xendit / Tripay) untuk QRIS, VA, atau Credit Card.

Order History & Tracking:

Riwayat transaksi lengkap dengan status (Pending Payment, Paid, Processing, Shipped, Completed, Cancelled).

### Pengambilan di Toko (Pickup in Store Flow):

Saat di halaman Checkout,

pengguna memilih Metode Pengiriman: Kurir (Shipping) atau Ambil di Toko (Pickup in Store).

Jika memilih Pickup in Store:Pengguna memilih lokasi cabang toko yang tersedia.Ongkos kirim (shipping_cost) dihitung Rp 0.

Setelah pembayaran berhasil (Paid), sistem generate Pickup Code (misal: PK-88219) atau QR Code.

Status Pesanan berubah menjadi: Ready for Pickup $\rightarrow$ Customer datang ke toko dan menunjukkan Pickup Code ke kasir $\rightarrow$ Kasir/Admin mengkonfirmasi via Admin Panel $\rightarrow$ Status berubah menjadi Completed.

## 3. 2.  Merchant & Administrative Features (Go Fiber + Next.js Admin Panel)
A. Product & Inventory Management (CRUD)
Manajemen Kategori & Sub-kategori.

CRUD Produk (Judul, Slug, Deskripsi, Harga Dasar).

CRUD Varian Produk (Kombinasi Ukuran, Warna opsional, SKU, Stok khusus per varian).

Upload & Pengelolaan Galeri Multi-Gambar Produk.

B. Order & Fulfillment Management
Melihat daftar pesanan masuk.

Mengubah status pesanan (Update Resi Pengiriman, Mark as Shipped).

C. Marketing & Reviews Management
CRUD Kode Kupon Promo (Jenis Potongan: Persentase atau Fixed Amount, Tanggal Kadaluarsa, Minimal Transaksi).

Moderasi Ulasan Produk.
Technical Stack Details
Frontend: Next.js (App Router), React, Tailwind CSS, State Management (Zustand/Redux Toolkit) untuk Keranjang Belanja.

Backend: Go (Fiber framework) untuk throughput tinggi, Gorm / Sqlx / Pgx driver.

Database: PostgreSQL dengan relasi terstruktur dan indexing pada slug, category_id, serta variant_id.

Third-Party Integrations:

Payment Gateway: Midtrans / Xendit (Webhook callback handler di Go Fiber).

Shipping Rate API: RajaOngkir / Biteship API.

Media Storage: Cloudinary / AWS S3 / Google Cloud Storage untuk asset gambar produk.

## 4. System Architecture & Technical Specifications

```text
  ┌─────────────────────────┐
  │   Next.js App Router    │  (Storefront & Admin UI)
  └────────────┬────────────┘
               │  REST API / JSON (JWT Auth)
  ┌────────────▼────────────┐
  │   Go Fiber Web Framework │  (Business Logic & API Endpoints)
  └────────────┬────────────┘
               │  SQL / pgx / GORM
  ┌────────────▼────────────┐
  │   PostgreSQL Database   │  (Relational Storage)
  └─────────────────────────┘
```
## 5. Non-Functional Requirements
Performance: Response time endpoint API Go Fiber rata-rata < 100ms.

Security:

Hash Password menggunakan bcrypt / Argon2.

JWT untuk authentication/authorization.

Webhook verification signature dari Payment Gateway.

Input sanitization & Prepared Statements SQL injection protection.

Scalability & Data Integrity:

Menggunakan Database Transaction (BEGIN...COMMIT/ROLLBACK) pada Go Fiber saat proses reduksi stok sewaktu checkout untuk mencegah race condition.

## 6. Next Steps
Dengan PRD V1 yang telah dikunci ini, tahap berikutnya adalah:

Membuat ERD (Entity-Relationship Diagram) PostgreSQL versi baru yang mengakomodasi seluruh fitur PRD V1 ini (termasuk wishlist, reviews, coupons, orders, dan shipping addresses).

Rancangan REST API Spec / Contract antara Go Fiber dan Next.js.