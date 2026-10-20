package main

import (
	"errors"
	"fmt"
	"net/http"
)

// TODO(Level 1): lihat SOAL.md untuk kontrak lengkap tiap fungsi di bawah.
// Ganti setiap "panic" dengan implementasi yang benar.

// ErrTugasTidakDitemukan dikembalikan ketika ID tugas yang dicari/dihapus
// tidak ada di daftar.
var ErrTugasTidakDitemukan = errors.New("tugas tidak ditemukan")

// ErrInputKosong dikembalikan ketika judul tugas yang diberikan kosong
// (atau hanya berisi spasi).
var ErrInputKosong = errors.New("input tidak boleh kosong")

// Task merepresentasikan satu tugas.
type Task struct {
	ID      int
	Judul   string
	Selesai bool
}

// TokoTugas menyimpan seluruh tugas di memori (bukan database), penghitung
// ID berikutnya, serta catatan operasi (dipakai di Level 7).
type TokoTugas struct {
	Daftar []Task
	NextID int
	Log    []string
}

func TambahTugas(toko *TokoTugas, judul string) (Task, error) {
	panic("belum diimplementasikan")
}

func LihatTugas(toko *TokoTugas, id int) (Task, error) {
	panic("belum diimplementasikan")
}

func HapusTugas(toko *TokoTugas, id int) error {
	panic("belum diimplementasikan")
}

// HapusTugasTercatat memanggil HapusTugas, lalu memakai defer untuk
// MENCATAT hasilnya ke toko.Log -- baik saat berhasil maupun saat gagal.
func HapusTugasTercatat(toko *TokoTugas, id int) error {
	panic("belum diimplementasikan")
}

// AmankanPanggilan menjalankan fn. Kalau fn panic, AmankanPanggilan
// menangkapnya lewat recover dan mengembalikannya sebagai error biasa,
// alih-alih membiarkan panic itu merambat dan menghentikan program.
func AmankanPanggilan(fn func() error) (err error) {
	panic("belum diimplementasikan")
}

// AmankanHandler membungkus next: kalau next panic saat memproses satu
// request, server tetap hidup untuk request-request lain (request yang
// panic itu dijawab status 500).
func AmankanHandler(next http.HandlerFunc) http.HandlerFunc {
	panic("belum diimplementasikan")
}

// RekapStatus menghitung berapa tugas yang sudah selesai dan berapa yang
// belum, dikembalikan sebagai map dengan persis dua kunci: "selesai" dan
// "belum selesai".
func RekapStatus(toko *TokoTugas) map[string]int {
	panic("belum diimplementasikan")
}

// BuatHandlerTugas mengembalikan HandlerFunc yang menuliskan daftar tugas
// di toko sebagai teks biasa ke w (satu tugas per baris), diikuti satu
// baris ringkasan dari RekapStatus.
func BuatHandlerTugas(toko *TokoTugas) http.HandlerFunc {
	panic("belum diimplementasikan")
}

func main() {
	fmt.Println("Task Manager v0 - pertemuan 3")

	toko := &TokoTugas{}
	http.HandleFunc("/tasks", AmankanHandler(BuatHandlerTugas(toko)))

	fmt.Println("Server jalan di :8080")
	http.ListenAndServe(":8080", nil)
}
