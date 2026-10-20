package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const (
	statusAktif      = "aktif"
	statusDibatalkan = "dibatalkan"
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

type Blokir struct {
	ID           int    `json:"id"`
	SumberDayaID int    `json:"sumber_daya_id"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Alasan       string `json:"alasan"`
}

type Server struct {
	mu             sync.Mutex
	sumberDaya     []SumberDaya
	peminjaman     []Peminjaman
	blokir         []Blokir
	nextSumberDaya int
	nextPeminjaman int
	nextBlokir     int
}

type ringkasanSumberDaya struct {
	SumberDayaID     int    `json:"sumber_daya_id"`
	Nama             string `json:"nama"`
	JumlahPeminjaman int    `json:"jumlah_peminjaman"`
}

type ringkasanResponse struct {
	Tanggal          string                `json:"tanggal"`
	JumlahPeminjaman int                   `json:"jumlah_peminjaman"`
	TotalUnitSlot    int                   `json:"total_unit_slot"`
	PerSumberDaya    []ringkasanSumberDaya `json:"per_sumber_daya"`
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
	mux.HandleFunc("PUT /peminjaman/{id}", s.updatePeminjaman)
	mux.HandleFunc("DELETE /peminjaman/{id}", s.deletePeminjaman)

	mux.HandleFunc("POST /blokir", s.createBlokir)
	mux.HandleFunc("GET /blokir", s.listBlokir)
	mux.HandleFunc("GET /blokir/{id}", s.getBlokir)
	mux.HandleFunc("DELETE /blokir/{id}", s.deleteBlokir)

	mux.HandleFunc("GET /ringkasan", s.getRingkasan)
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

func decodeBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, "body JSON tidak valid")
		return false
	}
	return true
}

func idFromPath(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id harus angka")
		return 0, false
	}
	return id, true
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

func slotValid(mulai, selesai int) bool {
	return mulai >= 1 && mulai <= 8 && selesai >= 1 && selesai <= 8 && mulai <= selesai
}

func menahanSlot(status string) bool {
	return status == statusAktif
}

func overlap(aMulai, aSelesai, bMulai, bSelesai int) bool {
	return aMulai <= bSelesai && bMulai <= aSelesai
}

func (s *Server) indexSumberDaya(id int) int {
	for i, sd := range s.sumberDaya {
		if sd.ID == id {
			return i
		}
	}
	return -1
}

func (s *Server) indexPeminjaman(id int) int {
	for i, p := range s.peminjaman {
		if p.ID == id {
			return i
		}
	}
	return -1
}

func (s *Server) indexBlokir(id int) int {
	for i, b := range s.blokir {
		if b.ID == id {
			return i
		}
	}
	return -1
}

func (s *Server) bentrok(sd SumberDaya, tanggal string, mulai, selesai, jumlah, kecualiID int) bool {
	for slot := mulai; slot <= selesai; slot++ {
		terpakai := 0
		for _, p := range s.peminjaman {
			if p.ID != kecualiID &&
				p.SumberDayaID == sd.ID &&
				p.Tanggal == tanggal &&
				menahanSlot(p.Status) &&
				p.SlotMulai <= slot && slot <= p.SlotSelesai {
				terpakai += p.Jumlah
			}
		}
		if terpakai+jumlah > sd.Stok {
			return true
		}
	}
	return false
}

func (s *Server) kenaBlokir(id int, tanggal string, mulai, selesai int) bool {
	for _, b := range s.blokir {
		if b.SumberDayaID == id && b.Tanggal == tanggal &&
			overlap(b.SlotMulai, b.SlotSelesai, mulai, selesai) {
			return true
		}
	}
	return false
}

type sumberDayaRequest struct {
	Nama  string `json:"nama"`
	Jenis string `json:"jenis"`
	Stok  int    `json:"stok"`
}

func (s *Server) createSumberDaya(w http.ResponseWriter, r *http.Request) {
	var in sumberDayaRequest
	if !decodeBody(w, r, &in) {
		return
	}

	in.Nama = trim(in.Nama)
	in.Jenis = trim(in.Jenis)

	switch {
	case in.Nama == "" || utf8.RuneCountInString(in.Nama) > 50:
		writeError(w, http.StatusBadRequest, "nama wajib diisi dan maksimal 50 karakter")
		return
	case in.Jenis != "ruang" && in.Jenis != "alat":
		writeError(w, http.StatusBadRequest, "jenis harus ruang atau alat")
		return
	case in.Stok < 1:
		writeError(w, http.StatusBadRequest, "stok minimal 1")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextSumberDaya++
	sd := SumberDaya{ID: s.nextSumberDaya, Nama: in.Nama, Jenis: in.Jenis, Stok: in.Stok}
	s.sumberDaya = append(s.sumberDaya, sd)

	writeJSON(w, http.StatusCreated, sd)
}

func (s *Server) listSumberDaya(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	hasil := append([]SumberDaya{}, s.sumberDaya...)
	s.mu.Unlock()

	writeJSON(w, http.StatusOK, hasil)
}

func (s *Server) getSumberDaya(w http.ResponseWriter, r *http.Request) {
	id, ok := idFromPath(w, r)
	if !ok {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexSumberDaya(id)
	if idx == -1 {
		writeError(w, http.StatusNotFound, "sumber daya tidak ditemukan")
		return
	}
	writeJSON(w, http.StatusOK, s.sumberDaya[idx])
}

type peminjamanRequest struct {
	SumberDayaID int    `json:"sumber_daya_id"`
	Peminjam     string `json:"peminjam"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Jumlah       int    `json:"jumlah"`
}

func validasiPeminjaman(in *peminjamanRequest) string {
	in.Peminjam = trim(in.Peminjam)
	in.Tanggal = trim(in.Tanggal)

	switch {
	case in.Peminjam == "" || utf8.RuneCountInString(in.Peminjam) > 50:
		return "peminjam wajib diisi dan maksimal 50 karakter"
	case !tanggalValid(in.Tanggal):
		return "tanggal harus valid dengan format YYYY-MM-DD"
	case !slotValid(in.SlotMulai, in.SlotSelesai):
		return "slot harus 1-8 dan slot_mulai <= slot_selesai"
	case in.Jumlah < 1:
		return "jumlah minimal 1"
	}
	return ""
}

func (s *Server) cekAturanPeminjaman(in peminjamanRequest, kecualiID int) (int, string) {
	idx := s.indexSumberDaya(in.SumberDayaID)
	if idx == -1 {
		return http.StatusBadRequest, "sumber daya tidak ditemukan"
	}

	sd := s.sumberDaya[idx]
	if in.Jumlah > sd.Stok {
		return http.StatusBadRequest, "jumlah melebihi stok sumber daya"
	}
	if s.kenaBlokir(in.SumberDayaID, in.Tanggal, in.SlotMulai, in.SlotSelesai) {
		return http.StatusConflict, "slot sedang diblokir"
	}
	if s.bentrok(sd, in.Tanggal, in.SlotMulai, in.SlotSelesai, in.Jumlah, kecualiID) {
		return http.StatusConflict, "kapasitas slot tidak cukup"
	}
	return 0, ""
}

func (s *Server) createPeminjaman(w http.ResponseWriter, r *http.Request) {
	var in peminjamanRequest
	if !decodeBody(w, r, &in) {
		return
	}
	if pesan := validasiPeminjaman(&in); pesan != "" {
		writeError(w, http.StatusBadRequest, pesan)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if status, pesan := s.cekAturanPeminjaman(in, 0); status != 0 {
		writeError(w, status, pesan)
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
		Status:       statusAktif,
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

	s.mu.Lock()
	hasil := []Peminjaman{}
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
	s.mu.Unlock()

	writeJSON(w, http.StatusOK, hasil)
}

func (s *Server) getPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, ok := idFromPath(w, r)
	if !ok {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexPeminjaman(id)
	if idx == -1 {
		writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}
	writeJSON(w, http.StatusOK, s.peminjaman[idx])
}

func (s *Server) updatePeminjaman(w http.ResponseWriter, r *http.Request) {
	id, ok := idFromPath(w, r)
	if !ok {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexPeminjaman(id)
	if idx == -1 {
		writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}

	lama := s.peminjaman[idx]

	var in peminjamanRequest
	if !decodeBody(w, r, &in) {
		return
	}
	if pesan := validasiPeminjaman(&in); pesan != "" {
		writeError(w, http.StatusBadRequest, pesan)
		return
	}
	if status, pesan := s.cekAturanPeminjaman(in, lama.ID); status != 0 {
		writeError(w, status, pesan)
		return
	}

	lama.SumberDayaID = in.SumberDayaID
	lama.Peminjam = in.Peminjam
	lama.Tanggal = in.Tanggal
	lama.SlotMulai = in.SlotMulai
	lama.SlotSelesai = in.SlotSelesai
	lama.Jumlah = in.Jumlah
	s.peminjaman[idx] = lama

	writeJSON(w, http.StatusOK, lama)
}

func (s *Server) deletePeminjaman(w http.ResponseWriter, r *http.Request) {
	id, ok := idFromPath(w, r)
	if !ok {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexPeminjaman(id)
	if idx == -1 {
		writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}
	s.peminjaman = append(s.peminjaman[:idx], s.peminjaman[idx+1:]...)

	w.WriteHeader(http.StatusNoContent)
}

type blokirRequest struct {
	SumberDayaID int    `json:"sumber_daya_id"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Alasan       string `json:"alasan"`
}

func validasiBlokir(in *blokirRequest) string {
	in.Tanggal = trim(in.Tanggal)
	in.Alasan = trim(in.Alasan)

	switch {
	case !tanggalValid(in.Tanggal):
		return "tanggal harus valid dengan format YYYY-MM-DD"
	case !slotValid(in.SlotMulai, in.SlotSelesai):
		return "slot harus 1-8 dan slot_mulai <= slot_selesai"
	case in.Alasan == "":
		return "alasan wajib diisi"
	}
	return ""
}

func (s *Server) createBlokir(w http.ResponseWriter, r *http.Request) {
	var in blokirRequest
	if !decodeBody(w, r, &in) {
		return
	}
	if pesan := validasiBlokir(&in); pesan != "" {
		writeError(w, http.StatusBadRequest, pesan)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.indexSumberDaya(in.SumberDayaID) == -1 {
		writeError(w, http.StatusBadRequest, "sumber daya tidak ditemukan")
		return
	}

	s.nextBlokir++
	b := Blokir{
		ID:           s.nextBlokir,
		SumberDayaID: in.SumberDayaID,
		Tanggal:      in.Tanggal,
		SlotMulai:    in.SlotMulai,
		SlotSelesai:  in.SlotSelesai,
		Alasan:       in.Alasan,
	}
	s.blokir = append(s.blokir, b)

	for i := range s.peminjaman {
		p := &s.peminjaman[i]
		if p.SumberDayaID == b.SumberDayaID &&
			p.Tanggal == b.Tanggal &&
			menahanSlot(p.Status) &&
			overlap(p.SlotMulai, p.SlotSelesai, b.SlotMulai, b.SlotSelesai) {
			p.Status = statusDibatalkan
		}
	}

	writeJSON(w, http.StatusCreated, b)
}

func (s *Server) listBlokir(w http.ResponseWriter, r *http.Request) {
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

	s.mu.Lock()
	hasil := []Blokir{}
	for _, b := range s.blokir {
		if adaFilterSD && b.SumberDayaID != filterSD {
			continue
		}
		if tanggal != "" && b.Tanggal != tanggal {
			continue
		}
		hasil = append(hasil, b)
	}
	s.mu.Unlock()

	writeJSON(w, http.StatusOK, hasil)
}

func (s *Server) getBlokir(w http.ResponseWriter, r *http.Request) {
	id, ok := idFromPath(w, r)
	if !ok {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexBlokir(id)
	if idx == -1 {
		writeError(w, http.StatusNotFound, "blokir tidak ditemukan")
		return
	}
	writeJSON(w, http.StatusOK, s.blokir[idx])
}

func (s *Server) deleteBlokir(w http.ResponseWriter, r *http.Request) {
	id, ok := idFromPath(w, r)
	if !ok {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexBlokir(id)
	if idx == -1 {
		writeError(w, http.StatusNotFound, "blokir tidak ditemukan")
		return
	}
	s.blokir = append(s.blokir[:idx], s.blokir[idx+1:]...)

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) getRingkasan(w http.ResponseWriter, r *http.Request) {
	tanggal := trim(r.URL.Query().Get("tanggal"))
	if !tanggalValid(tanggal) {
		writeError(w, http.StatusBadRequest, "tanggal wajib diisi dengan format YYYY-MM-DD")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	perSD := map[int]int{}
	hasil := ringkasanResponse{
		Tanggal:       tanggal,
		PerSumberDaya: []ringkasanSumberDaya{},
	}

	for _, p := range s.peminjaman {
		if p.Tanggal == tanggal && menahanSlot(p.Status) {
			hasil.JumlahPeminjaman++
			hasil.TotalUnitSlot += p.Jumlah * (p.SlotSelesai - p.SlotMulai + 1)
			perSD[p.SumberDayaID]++
		}
	}

	for _, sd := range s.sumberDaya {
		hasil.PerSumberDaya = append(hasil.PerSumberDaya, ringkasanSumberDaya{
			SumberDayaID:     sd.ID,
			Nama:             sd.Nama,
			JumlahPeminjaman: perSD[sd.ID],
		})
	}

	writeJSON(w, http.StatusOK, hasil)
}
