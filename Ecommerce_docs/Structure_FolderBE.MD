Penjelasan Peran Setiap Layer
1. cmd/api/main.go
Tempat eksekusi pertama. Berisi proses memuat konfigurasi, menghubungkan database, melakukan Dependency Injection (menghubungkan Repository → UseCase → Handler), mendaftarkan routes Go Fiber, dan menjalankan HTTP server.

2. internal/domain/
Berisi Entity/Model Structs yang merepresentasikan tabel database dan Interfaces Contract. Layer ini murni tidak memiliki ketergantungan (dependency) pada framework Fiber maupun GORM.

Contoh Interface:

Go
type ProductRepository interface {
    FindAll(ctx context.Context, filter ProductFilter) ([]Product, error)
    FindBySlug(ctx context.Context, slug string) (*Product, error)
}
3. internal/repository/
Layer yang berinteraksi langsung dengan database PostgreSQL menggunakan GORM atau SQL query. Mengimplementasikan interface yang ditentukan di domain.

4. internal/usecase/
Tempat seluruh Aturan Bisnis (Business Logic) berjalan. Contoh: mengecek stok varian sebelum checkout, memverifikasi kupon diskon, dan menghitung total pembayaran pesanan.

5. internal/handler/http/
Presentation Layer yang menerima HTTP Request dari Go Fiber, memvalidasi payload JSON input, memanggil usecase, dan mengembalikan format standard response JSON ke Next.js.

6. internal/platform/
Tempat menyimpan integrasi pihak ketiga (3rd party SDKs/APIs) seperti Payment Gateway (Midtrans), API Ongkos Kirim (Biteship/RajaOngkir), dan Cloud Storage untuk media upload.

package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func Success(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCodeBerikut adalah struktur folder Go Fiber menggunakan **Clean Architecture / Layered Architecture** yang modular, scalable, dan siap digunakan untuk project backend e-commerce.

```text
ecommerce-backend/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point aplikasi
├── internal/
│   ├── config/
│   │   ├── config.go               # Struct & parser env configuration
│   │   └── database.go             # Koneksi Database (GORM, PostgreSQL, Redis)
│   │
│   ├── domain/                     # Entities & Core Business Data Structures
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── order.go
│   │   ├── cart.go
│   │   └── payment.go
│   │
│   ├── delivery/                   # Interface / Adapter layer (HTTP / Fiber handlers)
│   │   └── http/
│   │       ├── router/
│   │       │   └── router.go       # Setup REST API endpoints
│   │       ├── middleware/
│   │       │   ├── auth.go         # JWT Auth Middleware
│   │       │   └── error.go        # Custom Error Handler Fiber
│   │       └── handler/
│   │           ├── auth_handler.go
│   │           ├── product_handler.go
│   │           ├── cart_handler.go
│   │           └── order_handler.go
│   │
│   ├── usecase/                    # Business Logic Layer
│   │   ├── auth_usecase.go
│   │   ├── product_usecase.go
│   │   ├── cart_usecase.go
│   │   └── order_usecase.go
│   │
│   ├── repository/                 # Data Access Layer (Database / Cache / External Service)
│   │   ├── user_repository.go
│   │   ├── product_repository.go
│   │   ├── order_repository.go
│   │   └── payment_gateway/        # External Gateway Integration (e.g., Midtrans)
│   │       └── midtrans.go
│   │
│   └── pkg/                        # Utility & Reusable Helper Functions
│       ├── dto/                    # Request & Response Payload Structs
│       │   ├── product_dto.go
│       │   └── order_dto.go
│       ├── utils/                  # Password hashing, JWT generator, string utils
│       └── validator/              # Custom request validation logic
│
├── db/
│   └── migrations/                 # File-file SQL migration (.sql / golang-migrate)
│
├── .env.example                    # Template environment variables
├── .gitignore
├── go.mod
├── go.sum
└── Makefile                        # Shortcut command (run, build, migrate)