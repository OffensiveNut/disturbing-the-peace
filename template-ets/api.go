package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

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

type Server struct {
	sumberDaya     []SumberDaya
	peminjaman     []Peminjaman
	nextSumberDaya int
	nextPeminjaman int
}

func NewHandler() http.Handler {
	s := &Server{}
	mux := http.NewServeMux()

	mux.HandleFunc("POST /sumber-daya", s.createSumberDaya)
	mux.HandleFunc("GET /sumber-daya", s.listSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", s.getSumberDaya)

	mux.HandleFunc("POST /peminjaman", s.createPeminjaman)
	mux.HandleFunc("GET /peminjaman", s.listPeminjaman)
	mux.HandleFunc("GET /peminjaman/{id}", s.getPeminjaman)

	return mux
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, pesan string) {
	writeJSON(w, status, map[string]string{"error": pesan})
}

func pathID(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("id"))
}

func trim(s string) string {
	return strings.TrimSpace(s)
}

func tanggalValid(s string) bool {
	if len(s) != 10 {
		return false
	}
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}

func (s *Server) findSumberDaya(id int) (SumberDaya, bool) {
	for _, sd := range s.sumberDaya {
		if sd.ID == id {
			return sd, true
		}
	}
	return SumberDaya{}, false
}

func (s *Server) findPeminjaman(id int) (Peminjaman, bool) {
	for _, p := range s.peminjaman {
		if p.ID == id {
			return p, true
		}
	}
	return Peminjaman{}, false
}

type sumberDayaRequest struct {
	Nama  string `json:"nama"`
	Jenis string `json:"jenis"`
	Stok  int    `json:"stok"`
}

func (s *Server) createSumberDaya(w http.ResponseWriter, r *http.Request) {
	var in sumberDayaRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "body JSON tidak valid")
		return
	}

	in.Nama = trim(in.Nama)
	in.Jenis = trim(in.Jenis)

	if in.Nama == "" || utf8.RuneCountInString(in.Nama) > 50 {
		writeError(w, http.StatusBadRequest, "nama wajib diisi dan maksimal 50 karakter")
		return
	}
	if in.Jenis != "ruang" && in.Jenis != "alat" {
		writeError(w, http.StatusBadRequest, "jenis harus ruang atau alat")
		return
	}
	if in.Stok < 1 {
		writeError(w, http.StatusBadRequest, "stok minimal 1")
		return
	}

	s.nextSumberDaya++
	sd := SumberDaya{
		ID:    s.nextSumberDaya,
		Nama:  in.Nama,
		Jenis: in.Jenis,
		Stok:  in.Stok,
	}
	s.sumberDaya = append(s.sumberDaya, sd)

	writeJSON(w, http.StatusCreated, sd)
}

func (s *Server) listSumberDaya(w http.ResponseWriter, r *http.Request) {
	hasil := make([]SumberDaya, len(s.sumberDaya))
	copy(hasil, s.sumberDaya)

	writeJSON(w, http.StatusOK, hasil)
}

func (s *Server) getSumberDaya(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id harus angka")
		return
	}

	sd, ok := s.findSumberDaya(id)
	if !ok {
		writeError(w, http.StatusNotFound, "sumber daya tidak ditemukan")
		return
	}
	writeJSON(w, http.StatusOK, sd)
}

type peminjamanRequest struct {
	SumberDayaID int    `json:"sumber_daya_id"`
	Peminjam     string `json:"peminjam"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Jumlah       int    `json:"jumlah"`
}

func (s *Server) createPeminjaman(w http.ResponseWriter, r *http.Request) {
	var in peminjamanRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "body JSON tidak valid")
		return
	}

	in.Peminjam = trim(in.Peminjam)
	in.Tanggal = trim(in.Tanggal)

	if in.Peminjam == "" || utf8.RuneCountInString(in.Peminjam) > 50 {
		writeError(w, http.StatusBadRequest, "peminjam wajib diisi dan maksimal 50 karakter")
		return
	}
	if !tanggalValid(in.Tanggal) {
		writeError(w, http.StatusBadRequest, "tanggal harus valid dengan format YYYY-MM-DD")
		return
	}
	if in.SlotMulai < 1 || in.SlotMulai > 8 ||
		in.SlotSelesai < 1 || in.SlotSelesai > 8 ||
		in.SlotMulai > in.SlotSelesai {
		writeError(w, http.StatusBadRequest, "slot harus 1-8 dan slot_mulai <= slot_selesai")
		return
	}
	if in.Jumlah < 1 {
		writeError(w, http.StatusBadRequest, "jumlah minimal 1")
		return
	}

	sd, ok := s.findSumberDaya(in.SumberDayaID)
	if !ok {
		writeError(w, http.StatusBadRequest, "sumber daya tidak ditemukan")
		return
	}
	if in.Jumlah > sd.Stok {
		writeError(w, http.StatusBadRequest, "jumlah melebihi stok sumber daya")
		return
	}

	s.nextPeminjaman++
	p := Peminjaman{
		ID:           s.nextPeminjaman,
		SumberDayaID: in.SumberDayaID,
		Peminjam:     in.Peminjam,
		Tanggal:      in.Tanggal,
		SlotMulai:    in.SlotMulai,
		SlotSelesai:  in.SlotSelesai,
		Jumlah:       in.Jumlah,
		Status:       "aktif",
	}
	s.peminjaman = append(s.peminjaman, p)

	writeJSON(w, http.StatusCreated, p)
}

func (s *Server) listPeminjaman(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var filterSD int
	adaFilterSD := false
	if q.Has("sumber_daya_id") {
		n, err := strconv.Atoi(q.Get("sumber_daya_id"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "sumber_daya_id harus angka")
			return
		}
		filterSD = n
		adaFilterSD = true
	}

	tanggal := trim(q.Get("tanggal"))
	status := trim(q.Get("status"))

	hasil := make([]Peminjaman, 0)
	for _, p := range s.peminjaman {
		if adaFilterSD && p.SumberDayaID != filterSD {
			continue
		}
		if tanggal != "" && p.Tanggal != tanggal {
			continue
		}
		if status != "" && p.Status != status {
			continue
		}
		hasil = append(hasil, p)
	}

	writeJSON(w, http.StatusOK, hasil)
}

func (s *Server) getPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id harus angka")
		return
	}

	p, ok := s.findPeminjaman(id)
	if !ok {
		writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}
	writeJSON(w, http.StatusOK, p)
}
