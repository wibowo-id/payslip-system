# Payslip System 🧾

Sistem penggajian terintegrasi berbasis web dengan fitur absensi, lembur, penggantian biaya (reimbursement), dan slip gaji. Dibangun dengan bahasa Go menggunakan framework **Gin** dan ORM **GORM**, serta menggunakan basis data PostgreSQL/SQLite untuk keperluan pengujian.

---

## ✨ Fitur Utama

- ✅ **Manajemen Pengguna (User Auth)**
  - Registrasi dan login dengan JWT.
  - Middleware autentikasi untuk rute API.

- 🕑 **Absensi (Attendance)**
  - Check-in karyawan pada hari kerja.
  - Deteksi otomatis hari libur (weekend).

- 🕒 **Lembur (Overtime)**
  - Pengajuan lembur maksimal 3 jam/hari.
  - Tidak boleh ada pengajuan duplikat pada tanggal yang sama.

- 💸 **Reimbursement**
  - Pengajuan penggantian biaya oleh karyawan.

- 📄 **Slip Gaji (Payslip)**
  - Perhitungan total take-home pay.
  - Komponen: Gaji, lembur, reimbursement.

- 📅 **Payroll Period**
  - Penutupan dan pembukaan periode gaji.

---

## 📂 Struktur Proyek

payslip-system/
├── cmd/ # Entry point aplikasi (main.go)
├── config/ # File konfigurasi (.env, dsb)
├── internal/
│ ├── auth/ # Modul autentikasi
│ ├── attendance/ # Modul absensi
│ ├── overtime/ # Modul lembur
│ ├── reimbursement/ # Modul reimbursement
│ ├── payroll/ # Modul periode gaji
│ ├── payslip/ # Modul slip gaji
│ └── user/ # Manajemen pengguna
├── pkg/
│ ├── middleware/ # Middleware JWT dan mock
│ └── logger/ # Logging utilitas
├── tests/ # Integrasi end-to-end test
└── go.mod / go.sum # Modul Go


---

## 🚀 Cara Menjalankan

### 1. Clone repositori

```bash
git clone https://github.com/namauser/payslip-system.git
cd payslip-system

2. Jalankan Aplikasi
Dengan SQLite (default testing mode)

go run cmd/server/main.go

Dengan PostgreSQL

Set environment variable .env:

DB_DRIVER=postgres
DB_DSN=postgres://user:pass@localhost:5432/payslip_db?sslmode=disable
JWT_SECRET=secret

Lalu:

go run cmd/server/main.go

🧪 Testing
Unit Test (dengan SQLite in-memory)

go clean -testcache
go test ./internal/... -v

Hasil test akan mencakup seluruh modul: attendance, auth, overtime, reimbursement, payslip, payroll, user.
🔒 Autentikasi

    Menggunakan JWT dengan Authorization: Bearer <token>.

    Middleware AuthOnly() akan memvalidasi token dan inject user_id, role ke gin.Context.

🧰 Tools dan Teknologi

    Go (Golang) v1.20+

    Gin – Web framework

    GORM – ORM untuk DB

    SQLite / PostgreSQL – Database

    JWT – Autentikasi

    Testify – Unit testing

    dotenv – Konfigurasi environment

🛠️ Pengembangan & Testing

Untuk pengujian handler dengan middleware:

    Gunakan MockAuthMiddleware di test.

    Gunakan NowFunc agar waktu bisa di-mock saat test service.

👤 Kontributor

    Chandra Wibowo — @github

📄 Lisensi

MIT License © 2025