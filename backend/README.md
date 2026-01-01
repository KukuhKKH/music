# App Music (Backend)

Deskripsi singkat aplikasi backend untuk proyek "app-music".

Aplikasi ini dibangun menggunakan Go + Fiber dengan modul-modul untuk autentikasi dan manajemen track. Menyertakan dokumentasi Swagger, migrasi/seed database, dan konfigurasi via file TOML.

Fitur utama
- Web framework: Fiber
- Dependency injection: Uber Fx
- ORM: GORM (MySQL)
- Logging: zerolog
- Autentikasi: JWT
- Dokumentasi API: Swagger

Prasyarat
- Go 1.24 atau lebih baru
- MySQL (atau database yang kompatibel dengan DSN MySQL)
- (Opsional) Sertifikat TLS jika ingin menjalankan dengan TLS

Persiapan konfigurasi
1. Salin contoh konfigurasi:
   - `config/config.toml.example` -> `config/config.toml`
2. Edit `config/config.toml` dan atur bagian berikut:
   - `app.port` (contoh: ":8080")
   - `db.mysql.dsn` (format: `<username>:<password>@tcp(<host>:<port>)/<database>?charset=utf8mb4&parseTime=True&loc=Local`)
   - `middleware.jwt.secret` untuk kunci JWT
   - `app.tls.cert-file` dan `app.tls.key-file` jika `app.tls.enable = true`

Contoh DSN MySQL (contoh umum):
`user:password@tcp(127.0.0.1:3306)/app_music?charset=utf8mb4&parseTime=True&loc=Local`

Build & Jalankan
Di direktori `backend`:

- Mengunduh dependensi (jika perlu):
```
go mod download
```

- Menjalankan langsung (development):
```
go run ./cmd/web
```

- Membangun executable:
```
go build -o bin/app ./cmd/web
```
  - Windows: executable akan berada di `bin\app.exe`.

- Menjalankan executable:
```
# Linux / macOS
./bin/app
# Windows (PowerShell)
.\bin\app.exe
```

Flags penting saat menjalankan
- `-migrate` : Jalankan migrasi database (AutoMigrate GORM)
- `-seed` : Isi data awal (seeder)

Contoh: jalankan executable dengan migrasi dan seeding
```
./bin/app -migrate -seed
```

Endpoint penting
- GET /ping — health check (mengembalikan "Pong! 👋")
- Swagger UI — `/swagger/index.html`

Database, Migrasi, dan Seeder
- Database dikonfigurasi melalui `config/config.toml` pada bagian `[db.mysql]` (field `dsn`).
- Migrasi model dijalankan dengan flag `-migrate`.
- Seeder disediakan (contoh: user seeder) dan dijalankan dengan flag `-seed`.

Dokumentasi API
- Dokumentasi Swagger berada di endpoint `/swagger/*`.
- Komentar swagger ditulis di `cmd/web/main.go`.

Pengembangan
- Struktur modul dan folder mengikuti pola modular:
  - `app/module` untuk fitur (auth, track)
  - `internal/bootstrap` untuk inisialisasi logger, server, dan database
  - `utils` berisi helper dan utilitas

Perintah tambahan
- Jalankan seluruh test (jika ada):
```
go test ./...
```

Kontribusi
- Silakan fork dan kirimkan pull request untuk perbaikan/fitur baru.

License
- Lihat file `LICENSE` untuk informasi lisensi.

