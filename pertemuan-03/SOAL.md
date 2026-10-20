# Pertemuan 3 — Task Manager v0: Fungsi, Error, & Server HTTP Pertama

Ini adalah **awal proyek berkelanjutan Task Manager**, yang akan terus tumbuh sampai pertemuan 12. Hari ini kalian membangun fondasinya: penyimpanan tugas di memori, fungsi-fungsi dasar dengan penanganan error ala Go, sedikit `defer`/`panic`/`recover`, dan satu endpoint HTTP pertama.

Berbeda dari Pertemuan 2, konsep-konsep di sini (fungsi, multiple return value, error) **sudah dibahas di kelas hari ini** — soal ini menguji penerapannya, bukan menyembunyikan konsep dasarnya. Signature tiap fungsi sudah disiapkan persis di starter `main.go`; tugas kalian mengganti setiap `panic("belum diimplementasikan")` dengan logika yang benar.

Semua level dicek otomatis lewat `go test`. Nilai mengikuti level tertinggi yang lolos **berurutan** dari 1 — kerjakan sejauh kemampuan, boleh tidak berurutan (`go test -v` selalu menunjukkan status tiap level apa adanya, lepas dari urutan pengerjaan).

**Catatan format respons `/tasks`:** kita belum membahas JSON (itu topik Pertemuan 6) — untuk sekarang cukup teks biasa, format persis dijelaskan di Level 10.

---

## Level 1 — `TambahTugas`

```go
func TambahTugas(toko *TokoTugas, judul string) (Task, error)
```

Tambahkan satu tugas baru ke `toko.Daftar`, dengan `ID` yang **unik dan bertambah** setiap kali dipanggil (tugas pertama ID berapa saja asal konsisten, tugas kedua harus berbeda dari yang pertama, dst — manfaatkan `toko.NextID`). Tugas baru selalu mulai dengan `Selesai: false`. Kembalikan tugas yang baru dibuat beserta `nil` sebagai error (penanganan kasus judul kosong ada di Level 6 — untuk level ini anggap input selalu valid).

## Level 2 — `LihatTugas`

```go
func LihatTugas(toko *TokoTugas, id int) (Task, error)
```

Cari tugas dengan `ID` yang cocok di `toko.Daftar`. Kalau ketemu, kembalikan tugas itu beserta `nil`. (Kasus tidak ketemu ada di Level 4.)

## Level 3 — `HapusTugas`

```go
func HapusTugas(toko *TokoTugas, id int) error
```

Cari dan hapus tugas dengan `ID` yang cocok dari `toko.Daftar`, lalu kembalikan `nil`. (Kasus tidak ketemu ada di Level 5.)

## Level 4 — `LihatTugas` (tidak ditemukan)

Lengkapi `LihatTugas` di atas: kalau `id` yang dicari **tidak ada** di `toko.Daftar`, kembalikan `Task{}` beserta error `ErrTugasTidakDitemukan` (variabel ini sudah didefinisikan di starter — pakai langsung, jangan buat error baru).

Ini pertama kalinya kalian benar-benar memakai `multiple return value` untuk mengomunikasikan kegagalan: pemanggil **wajib** mengecek nilai error sebelum memakai hasilnya.

## Level 5 — `HapusTugas` (tidak ditemukan)

Lengkapi `HapusTugas`: kalau `id` yang mau dihapus **tidak ada**, kembalikan `ErrTugasTidakDitemukan` dan **jangan ubah** `toko.Daftar` sama sekali.

## Level 6 — `TambahTugas` (input kosong)

Lengkapi `TambahTugas`: kalau `judul` kosong atau hanya berisi spasi (gunakan `strings.TrimSpace` untuk mengecek), kembalikan `Task{}` beserta error `ErrInputKosong`, dan **jangan tambahkan apa pun** ke `toko.Daftar`.

## Level 7 — `HapusTugasTercatat` (defer)

```go
func HapusTugasTercatat(toko *TokoTugas, id int) error
```

Panggil `HapusTugas` yang sudah kalian buat di Level 3 untuk melakukan penghapusan sesungguhnya. Pakai **`defer`** untuk mencatat hasil operasi itu ke `toko.Log` (sebuah `[]string`) — **baik saat berhasil maupun saat gagal**, defer itu harus tetap jalan dan menambahkan satu baris catatan.

Format catatan bebas, asal **memuat angka ID yang dihapus**, dan memuat kata **"berhasil"** untuk kasus sukses atau kata **"gagal"** untuk kasus gagal (huruf besar/kecil bebas). Contoh: `"hapus id=3: berhasil"` atau `"hapus id=99: gagal (tugas tidak ditemukan)"`.

Fungsi ini harus tetap mengembalikan error yang sama seperti yang dikembalikan `HapusTugas`.

## Level 8 — `AmankanPanggilan`

```go
func AmankanPanggilan(fn func() error) (err error)
```

Jalankan `fn`. Kalau `fn` selesai normal (return nil atau return error biasa), teruskan apa adanya. Tapi kalau `fn` **panic**, `AmankanPanggilan` harus menangkapnya lewat `recover()` dan mengembalikannya sebagai **error biasa** — bukan membiarkan panic itu merambat ke pemanggil. Bungkus pesan panic ke dalam error dengan `fmt.Errorf("operasi gagal: %v", r)` (`r` adalah nilai yang didapat dari `recover()`), supaya pesan asli panic-nya tetap terlihat di teks errornya.

## Level 9 — `AmankanHandler`

```go
func AmankanHandler(next http.HandlerFunc) http.HandlerFunc
```

Bungkus `next` (sebuah HTTP handler) dengan pola `defer` + `recover` yang sama seperti materi hari ini: kalau `next` panic saat memproses satu request, jangan biarkan program (server) ikut berhenti — tangkap panic-nya, lalu balas request itu dengan `w.WriteHeader(http.StatusInternalServerError)` (status 500). Kalau `next` **tidak** panic, `AmankanHandler` harus meneruskan perilaku `next` apa adanya (jangan ubah apa pun yang sudah ditulis `next` ke `w`).

Ini jaring pengaman "server tidak boleh mati gara-gara satu request error" yang dibahas di materi.

## Level 10 — `RekapStatus` + `BuatHandlerTugas` (endpoint `GET /tasks`)

```go
func RekapStatus(toko *TokoTugas) map[string]int
```

Hitung berapa tugas di `toko.Daftar` yang sudah `Selesai == true` dan berapa yang `Selesai == false`, lalu kembalikan sebagai **map** dengan persis dua kunci: `"selesai"` dan `"belum selesai"`, isinya jumlah masing-masing (dalam bentuk `int`). Kalau `toko.Daftar` kosong, kedua nilainya `0` (bukan kunci yang hilang — map-nya tetap punya dua kunci itu, nilainya saja nol).

Contoh: 3 tugas, 1 selesai dan 2 belum selesai → `map[string]int{"selesai": 1, "belum selesai": 2}`.

```go
func BuatHandlerTugas(toko *TokoTugas) http.HandlerFunc
```

Kembalikan sebuah `http.HandlerFunc` yang menuliskan **setiap tugas di `toko.Daftar`** ke `w`, satu tugas per baris, dengan format **persis**:

```
<ID>. <Judul> [<status>]
```

di mana `<status>` persis salah satu dari `selesai` (kalau `Selesai == true`) atau `belum selesai` (kalau `Selesai == false`). Setelah semua baris tugas, tambahkan **satu baris ringkasan** (pisahkan dengan baris kosong), dibangun dari `RekapStatus(toko)`, dengan format persis:

```
Ringkasan: <jumlah selesai> selesai, <jumlah belum selesai> belum selesai
```

Contoh lengkap, untuk tiga tugas (ID 1 "Belajar Go" belum selesai, ID 2 "Push ke GitHub" sudah selesai, ID 3 "Baca dokumentasi net/http" belum selesai), responsnya:

```
1. Belajar Go [belum selesai]
2. Push ke GitHub [selesai]
3. Baca dokumentasi net/http [belum selesai]

Ringkasan: 1 selesai, 2 belum selesai
```

Urutan baris tugas mengikuti urutan di `toko.Daftar` (tidak perlu diurutkan ulang) — tapi urutan **ringkasannya** ("selesai" disebut duluan, baru "belum selesai") harus persis seperti contoh, jadi ambil kedua nilainya langsung dari map lewat key masing-masing (`rekap["selesai"]`, `rekap["belum selesai"]`), **jangan** `range` map itu untuk membangun teksnya — urutan iterasi map di Go tidak terjamin konsisten.

Fungsi `main()` di starter sudah memasang `BuatHandlerTugas` (dibungkus `AmankanHandler`) ke path `/tasks` dan menyalakan server di `:8080` — setelah Level 9 & 10 selesai, jalankan `go run main.go` lalu buka `http://localhost:8080/tasks` di browser atau `curl` untuk melihatnya hidup.
