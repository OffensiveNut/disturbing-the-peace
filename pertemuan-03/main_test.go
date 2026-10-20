// File ini disediakan dosen untuk mengecek progres level secara otomatis.
// JANGAN DIUBAH — perubahan pada file ini tidak akan dipakai saat penilaian
// (dosen menimpa ulang file ini sebelum menjalankan grading).
package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
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

func TestLevels(t *testing.T) {
	t.Run("Level 1 - TambahTugas", func(t *testing.T) {
		toko := &TokoTugas{}
		var tugas1, tugas2 Task
		var err error

		aman(t, func() { tugas1, err = TambahTugas(toko, "Belajar Go") })
		if err != nil {
			t.Fatalf("TambahTugas(\"Belajar Go\") seharusnya tidak error, dapat: %v", err)
		}
		if tugas1.Judul != "Belajar Go" {
			t.Errorf("Judul = %q, ingin %q", tugas1.Judul, "Belajar Go")
		}
		if tugas1.ID == 0 {
			t.Errorf("ID tugas pertama seharusnya terisi (bukan 0), dapat %d", tugas1.ID)
		}
		if tugas1.Selesai {
			t.Errorf("tugas baru seharusnya belum selesai (Selesai=false)")
		}
		if len(toko.Daftar) != 1 {
			t.Fatalf("toko.Daftar seharusnya berisi 1 tugas, dapat %d", len(toko.Daftar))
		}

		aman(t, func() { tugas2, err = TambahTugas(toko, "Push ke GitHub") })
		if err != nil {
			t.Fatalf("TambahTugas(\"Push ke GitHub\") seharusnya tidak error, dapat: %v", err)
		}
		if tugas2.ID == tugas1.ID {
			t.Errorf("ID tugas kedua (%d) harus berbeda dari tugas pertama (%d)", tugas2.ID, tugas1.ID)
		}
		if len(toko.Daftar) != 2 {
			t.Errorf("toko.Daftar seharusnya berisi 2 tugas setelah 2x TambahTugas, dapat %d", len(toko.Daftar))
		}
	})

	t.Run("Level 2 - LihatTugas", func(t *testing.T) {
		toko := &TokoTugas{Daftar: []Task{
			{ID: 1, Judul: "Belajar Go"},
			{ID: 2, Judul: "Push ke GitHub"},
		}, NextID: 2}

		var tugas Task
		var err error
		aman(t, func() { tugas, err = LihatTugas(toko, 2) })
		if err != nil {
			t.Fatalf("LihatTugas(toko, 2) seharusnya tidak error, dapat: %v", err)
		}
		if tugas.ID != 2 || tugas.Judul != "Push ke GitHub" {
			t.Errorf("LihatTugas(toko, 2) = %+v, ingin ID=2 Judul=\"Push ke GitHub\"", tugas)
		}
	})

	t.Run("Level 3 - HapusTugas", func(t *testing.T) {
		toko := &TokoTugas{Daftar: []Task{
			{ID: 1, Judul: "Belajar Go"},
			{ID: 2, Judul: "Push ke GitHub"},
		}, NextID: 2}

		var err error
		aman(t, func() { err = HapusTugas(toko, 1) })
		if err != nil {
			t.Fatalf("HapusTugas(toko, 1) seharusnya tidak error, dapat: %v", err)
		}
		if len(toko.Daftar) != 1 {
			t.Fatalf("toko.Daftar seharusnya tersisa 1 tugas setelah hapus, dapat %d", len(toko.Daftar))
		}
		for _, sisa := range toko.Daftar {
			if sisa.ID == 1 {
				t.Errorf("tugas ID=1 masih ada di toko.Daftar setelah dihapus")
			}
		}
	})

	t.Run("Level 4 - LihatTugas (tidak ditemukan)", func(t *testing.T) {
		toko := &TokoTugas{Daftar: []Task{{ID: 1, Judul: "Belajar Go"}}, NextID: 1}

		var err error
		aman(t, func() { _, err = LihatTugas(toko, 999) })
		if err == nil {
			t.Fatal("LihatTugas dengan id yang tidak ada seharusnya mengembalikan error, dapat nil")
		}
		if !errors.Is(err, ErrTugasTidakDitemukan) {
			t.Errorf("error yang dikembalikan harus ErrTugasTidakDitemukan (lewat errors.Is), dapat: %v", err)
		}
	})

	t.Run("Level 5 - HapusTugas (tidak ditemukan)", func(t *testing.T) {
		toko := &TokoTugas{Daftar: []Task{{ID: 1, Judul: "Belajar Go"}}, NextID: 1}

		var err error
		aman(t, func() { err = HapusTugas(toko, 999) })
		if err == nil {
			t.Fatal("HapusTugas dengan id yang tidak ada seharusnya mengembalikan error, dapat nil")
		}
		if !errors.Is(err, ErrTugasTidakDitemukan) {
			t.Errorf("error yang dikembalikan harus ErrTugasTidakDitemukan (lewat errors.Is), dapat: %v", err)
		}
		if len(toko.Daftar) != 1 {
			t.Errorf("toko.Daftar tidak boleh berubah kalau penghapusan gagal, dapat %d tugas", len(toko.Daftar))
		}
	})

	t.Run("Level 6 - TambahTugas (input kosong)", func(t *testing.T) {
		toko := &TokoTugas{}

		var err error
		aman(t, func() { _, err = TambahTugas(toko, "") })
		if err == nil {
			t.Error("TambahTugas dengan judul kosong seharusnya mengembalikan error, dapat nil")
		} else if !errors.Is(err, ErrInputKosong) {
			t.Errorf("error yang dikembalikan harus ErrInputKosong (lewat errors.Is), dapat: %v", err)
		}

		aman(t, func() { _, err = TambahTugas(toko, "   ") })
		if err == nil {
			t.Error("TambahTugas dengan judul hanya spasi seharusnya mengembalikan error, dapat nil")
		} else if !errors.Is(err, ErrInputKosong) {
			t.Errorf("error yang dikembalikan harus ErrInputKosong (lewat errors.Is), dapat: %v", err)
		}

		if len(toko.Daftar) != 0 {
			t.Errorf("toko.Daftar tidak boleh bertambah kalau TambahTugas gagal, dapat %d tugas", len(toko.Daftar))
		}
	})

	t.Run("Level 7 - HapusTugasTercatat (defer)", func(t *testing.T) {
		toko := &TokoTugas{Daftar: []Task{{ID: 1, Judul: "Belajar Go"}}, NextID: 1}

		var err error
		aman(t, func() { err = HapusTugasTercatat(toko, 1) })
		if err != nil {
			t.Fatalf("HapusTugasTercatat(toko, 1) seharusnya tidak error, dapat: %v", err)
		}
		if len(toko.Log) != 1 {
			t.Fatalf("toko.Log seharusnya berisi 1 catatan setelah 1x panggilan, dapat %d", len(toko.Log))
		}
		if !strings.Contains(toko.Log[0], "1") || !strings.Contains(strings.ToLower(toko.Log[0]), "berhasil") {
			t.Errorf("catatan sukses harus memuat id dan kata \"berhasil\", dapat: %q", toko.Log[0])
		}

		// Kasus gagal: defer TETAP harus mencatat meski HapusTugas error.
		aman(t, func() { err = HapusTugasTercatat(toko, 999) })
		if err == nil {
			t.Fatal("HapusTugasTercatat(toko, 999) seharusnya mengembalikan error (id tidak ada)")
		}
		if len(toko.Log) != 2 {
			t.Fatalf("toko.Log seharusnya bertambah jadi 2 catatan walau operasinya gagal (defer wajib tetap jalan), dapat %d", len(toko.Log))
		}
		if !strings.Contains(strings.ToLower(toko.Log[1]), "gagal") {
			t.Errorf("catatan gagal harus memuat kata \"gagal\", dapat: %q", toko.Log[1])
		}
	})

	t.Run("Level 8 - AmankanPanggilan", func(t *testing.T) {
		var err error
		aman(t, func() { err = AmankanPanggilan(func() error { return nil }) })
		if err != nil {
			t.Errorf("AmankanPanggilan dengan fungsi yang sukses seharusnya mengembalikan nil, dapat: %v", err)
		}

		errAsli := ErrTugasTidakDitemukan
		aman(t, func() { err = AmankanPanggilan(func() error { return errAsli }) })
		if !errors.Is(err, errAsli) {
			t.Errorf("AmankanPanggilan seharusnya meneruskan error biasa apa adanya, dapat: %v", err)
		}

		aman(t, func() {
			err = AmankanPanggilan(func() error {
				panic("kesalahan tak terduga")
			})
		})
		if err == nil {
			t.Fatal("AmankanPanggilan dengan fungsi yang panic seharusnya mengembalikan error, dapat nil")
		}
		if !strings.Contains(err.Error(), "kesalahan tak terduga") {
			t.Errorf("error hasil recover harus memuat pesan panic aslinya, dapat: %q", err.Error())
		}
	})

	t.Run("Level 9 - AmankanHandler", func(t *testing.T) {
		normal := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		var w1 *httptest.ResponseRecorder
		aman(t, func() {
			w1 = httptest.NewRecorder()
			AmankanHandler(normal)(w1, httptest.NewRequest(http.MethodGet, "/tasks", nil))
		})
		if w1.Code != http.StatusOK || w1.Body.String() != "OK" {
			t.Errorf("handler normal yang dibungkus AmankanHandler harus tetap jalan seperti biasa, dapat status=%d body=%q", w1.Code, w1.Body.String())
		}

		panik := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("handler ini sengaja panic")
		})
		var w2 *httptest.ResponseRecorder
		aman(t, func() {
			w2 = httptest.NewRecorder()
			AmankanHandler(panik)(w2, httptest.NewRequest(http.MethodGet, "/tasks", nil))
		})
		if w2.Code != http.StatusInternalServerError {
			t.Errorf("handler yang panic, setelah dibungkus AmankanHandler, harus menjawab status 500, dapat %d", w2.Code)
		}
	})

	t.Run("Level 10 - BuatHandlerTugas (GET /tasks)", func(t *testing.T) {
		toko := &TokoTugas{Daftar: []Task{
			{ID: 1, Judul: "Belajar Go", Selesai: false},
			{ID: 2, Judul: "Push ke GitHub", Selesai: true},
			{ID: 3, Judul: "Baca dokumentasi net/http", Selesai: false},
		}, NextID: 3}

		var rekap map[string]int
		aman(t, func() { rekap = RekapStatus(toko) })
		if rekap["selesai"] != 1 {
			t.Errorf("RekapStatus(toko)[\"selesai\"] = %d, ingin 1", rekap["selesai"])
		}
		if rekap["belum selesai"] != 2 {
			t.Errorf("RekapStatus(toko)[\"belum selesai\"] = %d, ingin 2", rekap["belum selesai"])
		}

		var w *httptest.ResponseRecorder
		aman(t, func() {
			handler := BuatHandlerTugas(toko)
			w = httptest.NewRecorder()
			handler(w, httptest.NewRequest(http.MethodGet, "/tasks", nil))
		})

		body := w.Body.String()
		if !strings.Contains(body, "1. Belajar Go [belum selesai]") {
			t.Errorf("body respons harus memuat baris \"1. Belajar Go [belum selesai]\", dapat:\n%s", body)
		}
		if !strings.Contains(body, "2. Push ke GitHub [selesai]") {
			t.Errorf("body respons harus memuat baris \"2. Push ke GitHub [selesai]\", dapat:\n%s", body)
		}
		if !strings.Contains(body, "Ringkasan: 1 selesai, 2 belum selesai") {
			t.Errorf("body respons harus memuat baris \"Ringkasan: 1 selesai, 2 belum selesai\", dapat:\n%s", body)
		}

		// RekapStatus pada toko kosong: kedua kunci tetap ada, nilainya 0.
		tokoKosong := &TokoTugas{}
		var rekapKosong map[string]int
		aman(t, func() { rekapKosong = RekapStatus(tokoKosong) })
		if rekapKosong["selesai"] != 0 || rekapKosong["belum selesai"] != 0 {
			t.Errorf("RekapStatus pada toko kosong seharusnya {selesai:0, belum selesai:0}, dapat: %v", rekapKosong)
		}
	})
}
