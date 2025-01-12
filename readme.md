# NEXSHOP Services API

API ini adalah aplikasi e-commerce menggunakan Golang dan PostgreSQL.

## Persiapan Sebelum Memulai

1. **Pastikan Database Sudah Tersedia**
   - Pastikan Anda sudah membuat database PostgreSQL dengan nama **`nexshop_db`**.

2. **Jalankan Query Yang Tersedia**
    - Buka file `schema.sql` pada aplikasi ini.
    - Copy dan Jalankan pada Database Client seperti Dbeaver Dan lainnya.

3. **Konfigurasi Aplikasi**
   - Buka file `.env-sample-config` pada aplikasi ini.
   - Sesuaikan pengaturan koneksi database dan konfigurasi lainnya sesuai dengan lingkungan Anda.

## Langkah-Langkah Menjalankan Aplikasi

### 1. Mendownload seluruh kebutuhan aplikasi
    - Setelah setup dan konfigurasi  database  selesai , download package yang dibutuhkan pada aplikasi dengan perintah :
    ```bash
    go mod tidy
    ```
    - Aplikasi akan mendownload seluruh package yang dibutuhkan.

### 2. Menjalankan aplikasi
    - Setelah selesai mendownload package, jalankan aplikasi denga perintah :
    ```bash
    go run ./cmd/main.go
    ```
    - Aplikasi akan berjalan pada http://localhost:8080 jika anda tidak mengatur portnya pada .env file

## Struktur Direktori
- **cmd/**: Entry point aplikasi
- **internal/**: Bisnis logic untuk aplikasi
- **pkg/**: Menyimpan konfigurasi yang berkaitan dengan layanan pihak ketiga seperti postgres, jwt dan lainnya.
- **internal/handlers/**: Lapisan yang menangani permintaan dari pengguna, baik dari aplikasi mobile maupun web.
- **migration/**: Menyimpan file migration SQL dan kode untuk mengatur dan memperbarui struktur database (coming soon).
- **internal/models/**: Berisi struktur data (constructs di Golang) yang memudahkan dalam membuat kontrak untuk request dan response.
- **internal/respositories/**: Lapisan yang berfungsi khusus untuk berinteraksi dengan database, termasuk operasi pencatatan dan pengambilan data.
- **internal/routes/**: Menyimpan definisi endpoint utama yang mengarahkan 
- **.env**: File konfigurasi utama untuk mengatur koneksi database dan parameter aplikasi lainnya.

## Teknologi yang Digunakan

- **Golang**: Backend aplikasi utama.
- **PostgreSQL**: Database untuk menyimpan data pengguna.