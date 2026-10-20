# PROSES

## 1. Persona dan masalahnya

Persona E adalah asisten laboratorium yang mengatur pemakaian lab. Lab kadang dipakai untuk praktikum atau servis, sehingga slot tertentu tidak boleh dipinjam. Fitur blokir jadwal menolak peminjaman baru pada slot yang diblokir dan membatalkan peminjaman aktif yang sudah ada pada slot tersebut.

## 2. Rancangan struktur data

Data disimpan di memori dalam slice. Saya memilih slice karena jumlah data kecil, urutan ID penting, dan tidak ada kebutuhan query kompleks. Setiap sumber daya dan peminjaman memiliki ID yang naik mulai dari 1.

```go
type SumberDaya struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Jenis string `json:"jenis"`
	Stok  int    `json:"stok"`
}

type Peminjaman struct {
	ID           int    `json:"id"`
	SumberDayaID int    `json:"sumber_daya_id"`
	Peminjam     string `json:"peminjam"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Jumlah       int    `json:"jumlah"`
	Status       string `json:"status"`
}

type Blokir struct {
	ID           int    `json:"id"`
	SumberDayaID int    `json:"sumber_daya_id"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Alasan       string `json:"alasan"`
}
```

Server menyimpan tiga slice dan counter ID. Satu `sync.Mutex` melindungi seluruh slice dan counter karena semua request membaca atau menulis data yang sama.

```go
type Server struct {
	mu             sync.Mutex
	sumberDaya     []SumberDaya
	peminjaman     []Peminjaman
	blokir         []Blokir
	nextSumberDaya int
	nextPeminjaman int
	nextBlokir     int
}
```

Ringkasan tidak disimpan sebagai counter terpisah. Data ringkasan dihitung ulang dari slice peminjaman dan sumber daya setiap kali `GET /ringkasan` dipanggil. Dengan begitu, data tetap konsisten setelah peminjaman dibuat, diubah, dihapus, atau dibatalkan oleh blokir.

## 3. Satu keputusan desain

Saya mengecek bentrok **per slot**, bukan hanya mengecek apakah dua rentang overlap. Untuk setiap slot pada rentang peminjaman baru, saya menjumlahkan `jumlah` semua peminjaman aktif yang memakai slot itu, lalu membandingkannya dengan `stok` sumber daya.

Pilihan lain adalah mengecek overlap rentang saja. Cara itu tidak cukup ketika stok lebih dari 1: dua peminjaman bisa overlap sebagian dan masih muat, atau justru melebihi kapasitas hanya di satu slot. Pengecekan per slot menangani stok banyak dan overlap sebagian dengan benar. Saat PUT, ID peminjaman yang sedang diedit dikecualikan agar peminjaman tidak dianggap bentrok dengan dirinya sendiri.

## 4. Bug atau kendala

### Kendala 1: Kontrak persona E tidak tersedia lengkap
- Gejala: tidak ada spesifikasi rinci di repo untuk path, field, dan status pembatalan pada fitur blokir.
- Cara menemukan: memeriksa dokumen dan kode yang tersedia, lalu mencocokkan dengan perilaku yang diminta pada deskripsi persona E.
- Penyebab dan solusi: saya menetapkan asumsi endpoint `POST /blokir`, field `alasan`, dan status peminjaman yang dibatalkan menjadi `dibatalkan`. Asumsi ini dipakai konsisten di kode, README, dan test. Jika kontrak dosen berbeda, path atau nama field perlu disesuaikan.

### Kendala 2: Ringkasan harus tetap konsisten setelah blokir dan hapus
- Gejala: ringkasan bisa menampilkan peminjaman yang sudah dibatalkan atau sudah dihapus jika hanya memakai counter yang tidak ikut diperbarui.
- Cara menemukan: menguji alur buat peminjaman, blokir slot, lalu minta ringkasan pada tanggal yang sama.
- Penyebab dan solusi: counter terpisah berisiko tidak sinkron dengan data utama. Solusinya, ringkasan dihitung saat request dari slice peminjaman yang masih ada, dan hanya menghitung status `aktif`. Peminjaman yang dibatalkan atau dihapus otomatis tidak ikut. `per_sumber_daya` juga diisi untuk semua sumber daya, termasuk yang jumlahnya nol, dan diurutkan berdasarkan ID.

### Kendala 3: Ekspektasi test bentrok sempat salah
- Gejala: test B5 sempat gagal karena mengharapkan 409 pada slot 4, tetapi server mengembalikan 201.
- Cara menemukan: menjalankan `go test -race -v ./...` dan membaca body respons.
- Penyebab dan solusi: pada stok 2, slot 4 baru terpakai 1 unit sehingga penambahan 1 unit masih muat. Ekspektasi test diperbaiki ke slot 3 yang memang sudah terpakai 2 unit. Ini kesalahan ekspektasi test, bukan kesalahan logika server.

## 5. Change request

Change request A meminta endpoint `GET /ringkasan?tanggal=YYYY-MM-DD` dengan respons:

```json
{
  "tanggal": "2026-11-02",
  "jumlah_peminjaman": 2,
  "total_unit_slot": 7,
  "per_sumber_daya": [
    {"sumber_daya_id": 1, "nama": "Aula A", "jumlah_peminjaman": 1},
    {"sumber_daya_id": 2, "nama": "Proyektor", "jumlah_peminjaman": 1},
    {"sumber_daya_id": 3, "nama": "Lab RPL", "jumlah_peminjaman": 0}
  ]
}
```

Bagian kode yang diubah adalah `getRingkasan` di `api.go`, ditambah struct `ringkasanResponse` dan `ringkasanSumberDaya`. Ringkasan dihitung dari data peminjaman yang masih ada dan hanya menghitung status `aktif`, yaitu status yang masih menahan slot. Peminjaman yang dibatalkan atau dihapus tidak dihitung. `total_unit_slot` dihitung dari `jumlah × (slot_selesai - slot_mulai + 1)` untuk setiap peminjaman aktif pada tanggal tersebut.

Keputusan desainnya: ringkasan dihitung saat request, bukan disimpan sebagai counter. Ini menghindari data ringkasan yang basi setelah create, PUT, DELETE, atau blokir. Semua sumber daya tetap muncul di `per_sumber_daya` meskipun jumlah peminjamannya nol, dan urutannya mengikuti ID. Jika tidak ada sumber daya, field tersebut tetap berupa array kosong `[]`. Tanggal yang kosong atau tidak valid menghasilkan 400; tanggal valid tanpa peminjaman menghasilkan 200 dengan angka nol.
