package main

import "fmt"

// TODO(Level 1): lihat SOAL.md untuk kontrak lengkap tiap fungsi di bawah.
// Ganti setiap "panic" dengan implementasi yang benar.

func HitungSubtotal(qty int, hargaSatuan float64) float64 {
	panic("belum diimplementasikan")
}

func HitungTotalPesanan(qty []int, hargaSatuan []float64) float64 {
	panic("belum diimplementasikan")
}

func TerapkanPajak(total float64, tarifPajak float64) float64 {
	panic("belum diimplementasikan")
}

func HitungDiskon(total float64) float64 {
	panic("belum diimplementasikan")
}

func TotalSetelahDiskon(qty []int, hargaSatuan []float64, tarifPajak float64) float64 {
	panic("belum diimplementasikan")
}

func ValidasiPesanan(qty []int, hargaSatuan []float64) (bool, string) {
	panic("belum diimplementasikan")
}

func TentukanStatus(total float64) string {
	panic("belum diimplementasikan")
}

func RingkasanPesanan(qty []int, hargaSatuan []float64, tarifPajak float64) string {
	panic("belum diimplementasikan")
}

// TODO(Level 9): signature ini SUDAH benar (cari tahu sendiri kenapa
// bentuknya begini - lihat SOAL.md) - tinggal implementasikan isinya.
func Total(harga ...float64) float64 {
	panic("belum diimplementasikan")
}

// TODO(Level 10, bonus): signature ini SUDAH benar (cari tahu sendiri
// kenapa ada dua nilai balik - lihat SOAL.md) - tinggal implementasikan isinya.
func HitungOngkosKirim(beratKg float64, jarakKm float64) (float64, error) {
	panic("belum diimplementasikan")
}

func main() {
	fmt.Println("Sales Order Processor - pertemuan 2")
}
