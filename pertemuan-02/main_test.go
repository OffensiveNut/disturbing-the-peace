// File ini disediakan dosen untuk mengecek progres level secara otomatis.
// JANGAN DIUBAH — perubahan pada file ini tidak akan dipakai saat penilaian
// (dosen menimpa ulang file ini sebelum menjalankan grading).
package main

import (
	"math"
	"strings"
	"testing"
)

// aman menjalankan fn dan mengubah panic (mis. dari starter yang belum
// diimplementasikan) jadi kegagalan test yang rapi, bukan crash seluruh
// binary test.
func aman(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("fungsi panic: %v (kemungkinan belum diimplementasikan)", r)
		}
	}()
	fn()
}

func dekat(a, b float64) bool {
	return math.Abs(a-b) < 0.01
}

func TestLevels(t *testing.T) {
	t.Run("Level 1 - HitungSubtotal", func(t *testing.T) {
		var got float64
		aman(t, func() { got = HitungSubtotal(3, 15000) })
		if !dekat(got, 45000) {
			t.Errorf("HitungSubtotal(3, 15000) = %v, ingin 45000", got)
		}
		aman(t, func() { got = HitungSubtotal(0, 15000) })
		if !dekat(got, 0) {
			t.Errorf("HitungSubtotal(0, 15000) = %v, ingin 0", got)
		}
	})

	t.Run("Level 2 - HitungTotalPesanan", func(t *testing.T) {
		var got float64
		aman(t, func() {
			got = HitungTotalPesanan([]int{2, 5, 1}, []float64{15000, 8000, 25000})
		})
		if !dekat(got, 95000) {
			t.Errorf("HitungTotalPesanan(...) = %v, ingin 95000", got)
		}
		aman(t, func() {
			got = HitungTotalPesanan([]int{1, 2}, []float64{1000})
		})
		if !dekat(got, 0) {
			t.Errorf("HitungTotalPesanan dengan panjang tidak sama = %v, ingin 0", got)
		}
	})

	t.Run("Level 3 - TerapkanPajak", func(t *testing.T) {
		var got float64
		aman(t, func() { got = TerapkanPajak(100000, 0.11) })
		if !dekat(got, 111000) {
			t.Errorf("TerapkanPajak(100000, 0.11) = %v, ingin 111000", got)
		}
	})

	t.Run("Level 4 - HitungDiskon", func(t *testing.T) {
		kasus := []struct {
			total, ingin float64
		}{
			{1200000, 120000},
			{600000, 30000},
			{200000, 0},
			{500000, 25000},
			{1000000, 100000},
		}
		for _, k := range kasus {
			var got float64
			aman(t, func() { got = HitungDiskon(k.total) })
			if !dekat(got, k.ingin) {
				t.Errorf("HitungDiskon(%v) = %v, ingin %v", k.total, got, k.ingin)
			}
		}
	})

	t.Run("Level 5 - TotalSetelahDiskon", func(t *testing.T) {
		var got float64
		// subtotal 1.200.000 -> diskon 10% -> 1.080.000 -> pajak 11% -> 1.198.800
		aman(t, func() {
			got = TotalSetelahDiskon([]int{10, 10}, []float64{60000, 60000}, 0.11)
		})
		if !dekat(got, 1198800) {
			t.Errorf("TotalSetelahDiskon(...) = %v, ingin 1198800", got)
		}
	})

	t.Run("Level 6 - ValidasiPesanan", func(t *testing.T) {
		var ok bool
		var pesan string
		aman(t, func() { ok, pesan = ValidasiPesanan([]int{1, 2}, []float64{1000, 2000}) })
		if !ok {
			t.Errorf("pesanan valid seharusnya mengembalikan true, dapat false (pesan: %q)", pesan)
		}
		if strings.TrimSpace(pesan) != "" {
			t.Errorf("pesanan valid seharusnya pesan kosong, dapat %q", pesan)
		}

		kasusInvalid := []struct {
			nama  string
			qty   []int
			harga []float64
		}{
			{"panjang tidak sama", []int{1, 2}, []float64{1000}},
			{"qty nol/negatif", []int{1, -2}, []float64{1000, 2000}},
			{"harga nol/negatif", []int{1, 2}, []float64{1000, -2000}},
		}
		for _, k := range kasusInvalid {
			aman(t, func() { ok, pesan = ValidasiPesanan(k.qty, k.harga) })
			if ok {
				t.Errorf("kasus %q: seharusnya tidak valid, dapat true", k.nama)
			}
			if strings.TrimSpace(pesan) == "" {
				t.Errorf("kasus %q: pesan tidak valid tidak boleh kosong", k.nama)
			}
		}
	})

	t.Run("Level 7 - TentukanStatus", func(t *testing.T) {
		kasus := []struct {
			total float64
			ingin string
		}{
			{1500000, "Prioritas"},
			{1000001, "Prioritas"},
			{1000000, "Reguler"},
			{500000, "Reguler"},
			{100001, "Reguler"},
			{100000, "Hemat"},
			{5000, "Hemat"},
		}
		for _, k := range kasus {
			var got string
			aman(t, func() { got = TentukanStatus(k.total) })
			if got != k.ingin {
				t.Errorf("TentukanStatus(%v) = %q, ingin %q", k.total, got, k.ingin)
			}
		}
	})

	t.Run("Level 8 - RingkasanPesanan", func(t *testing.T) {
		var ringkasan string
		aman(t, func() {
			ringkasan = RingkasanPesanan([]int{10, 10}, []float64{60000, 60000}, 0.11)
		})
		low := strings.ToLower(ringkasan)

		// subtotal 1.200.000, diskon 120.000, total akhir 1.198.800, status Prioritas
		wajibMuncul := map[string][]string{
			"subtotal (1200000 / 1.200.000)":    {"1200000", "1.200.000", "1,200,000"},
			"diskon (120000 / 120.000)":         {"120000", "120.000", "120,000"},
			"total akhir (1198800 / 1.198.800)": {"1198800", "1.198.800", "1,198,800"},
			"status (Prioritas)":                {"prioritas"},
		}
		for label, alternatif := range wajibMuncul {
			cocok := false
			for _, a := range alternatif {
				if strings.Contains(low, strings.ToLower(a)) {
					cocok = true
					break
				}
			}
			if !cocok {
				t.Errorf("RingkasanPesanan harus memuat %s, tidak ditemukan di:\n%s", label, ringkasan)
			}
		}
	})

	t.Run("Level 9 - Total (variadic)", func(t *testing.T) {
		var got float64
		aman(t, func() { got = Total() })
		if !dekat(got, 0) {
			t.Errorf("Total() = %v, ingin 0", got)
		}
		aman(t, func() { got = Total(5) })
		if !dekat(got, 5) {
			t.Errorf("Total(5) = %v, ingin 5", got)
		}
		aman(t, func() { got = Total(1, 2, 3, 4, 5) })
		if !dekat(got, 15) {
			t.Errorf("Total(1,2,3,4,5) = %v, ingin 15", got)
		}
	})

	t.Run("Level 10 - HitungOngkosKirim (bonus)", func(t *testing.T) {
		var ongkos float64
		var err error
		aman(t, func() { ongkos, err = HitungOngkosKirim(5, 10) })
		if err != nil {
			t.Errorf("HitungOngkosKirim(5, 10) seharusnya tidak error, dapat: %v", err)
		}
		if !dekat(ongkos, 5*2000+10*3000) {
			t.Errorf("HitungOngkosKirim(5, 10) = %v, ingin %v", ongkos, 5*2000+10*3000)
		}

		aman(t, func() { _, err = HitungOngkosKirim(0, 10) })
		if err == nil {
			t.Error("HitungOngkosKirim(0, 10) seharusnya mengembalikan error (berat tidak valid)")
		}

		aman(t, func() { _, err = HitungOngkosKirim(5, -1) })
		if err == nil {
			t.Error("HitungOngkosKirim(5, -1) seharusnya mengembalikan error (jarak tidak valid)")
		}
	})
}
