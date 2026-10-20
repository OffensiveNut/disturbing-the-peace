# PROSES

## 1. Persona dan masalahnya
(2-3 kalimat dengan kata-katamu sendiri: siapa pengguna persona kamu, masalah apa yang diselesaikan fitur khasnya.)
persona E.

E. Asisten laboratorium
Lab dipakai praktikum atau servis sehingga tidak boleh dipinjam di jam tertentu
Blokir jadwal: slot yang diblokir menolak peminjaman baru dan membatalkan yang sudah ada

## 2. Rancangan struktur data
(Tulis struct dan slice/map apa yang kamu simpan, dan kenapa dipilih begitu. Rancang ini sebelum menulis kode, lalu catat kalau berubah di tengah jalan.)

``` go
type SumberDaya struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Jenis string `json:"jenis"`
	Stok  int    `json:"stok"`
}
```
menyesuaikan dengan soal Objek sumber daya:
``` JSON
{“id”: 1, “nama”: “Aula A”, “jenis”: “ruang”, “stok”: 1}
```
``` go
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
```
Objek peminjaman:
``` JSON
{
“id”: 1,
“sumber_daya_id”: 2,
“peminjam”: “Rina”,
“tanggal”: “2026-11-02”,
“slot_mulai”: 3,
“slot_selesai”: 4,
“jumlah”: 1,
“status”: “aktif”
}
```
``` go
type Server struct {
	sumberDaya     []SumberDaya
	peminjaman     []Peminjaman
	nextSumberDaya int
	nextPeminjaman int
}
```
## 3. Satu keputusan desain
(Satu keputusan yang kamu ambil, pilihan lain yang kamu pertimbangkan, dan alasanmu memilih yang ini. Contoh: kenapa cek bentrok per slot, bukan per rentang.)

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

### Bug 2
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

## 5. Change request
(Apa yang diminta, bagian kode mana yang kamu ubah, dan satu keputusan desain yang kamu ambil.)
