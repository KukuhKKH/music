# App Music

Aplikasi monorepo musik yang terdiri dari backend berbasis Go dan frontend berbasis Nuxt 4.

## 📂 Struktur Proyek

- **backend/**: API Server menggunakan Go (Fiber Framework).
- **frontend/**: Client App menggunakan Nuxt 4, TailwindCSS, dan Shadcn UI.
- **.gitlab-ci/**: Konfigurasi CI/CD dan Docker.

## 🚀 Memulai (Local Development)

### Backend

1. Masuk ke direktori backend:
   ```bash
   cd backend
   ```
2. Salin konfigurasi:
   ```bash
   cp config/config.toml.example config/config.toml
   ```
   > Sesuaikan konfigurasi database dan port di `config/config.toml`.

3. Jalankan aplikasi:
   ```bash
   go run ./cmd/web
   ```
   > Gunakan flag `-migrate` dan `-seed` untuk inisialisasi database pertama kali:
   > `go run ./cmd/web -migrate -seed`

### Frontend

1. Masuk ke direktori frontend:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   pnpm install
   ```
3. Jalankan server development:
   ```bash
   pnpm dev
   ```

## 🐳 Docker Build

Proyek ini memiliki konfigurasi build Docker di dalam folder `.gitlab-ci`.

Untuk membuild image backend dan frontend:

```bash
docker-compose -f .gitlab-ci/build/docker-compose.yml build
```

## 📚 Dokumentasi Lebih Lanjut

- **Backend Docs**: Lihat [backend/README.md](backend/README.md)
- **Frontend Docs**: Lihat `package.json` di folder frontend untuk skrip yang tersedia.
