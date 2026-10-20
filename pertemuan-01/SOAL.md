# Pertemuan 1 — Perkenalan: Setup Repo & Program Pertama

Tugas ini bersifat **onboarding**: memastikan environment (Go, IDE, akses git organization) siap dan kalian terbiasa dengan alur git yang dipakai sepanjang semester. Kompleksitas *kode* sengaja ringan — tapi **semua 10 level dicek otomatis lewat `go test`**, termasuk bagian yang biasanya reflektif (penjelasan di README.md di root repo kalian). Ikuti persis nama fungsi dan nama heading README yang diminta — pengecekan otomatis mencocokkan teks/perilaku secara presisi, bukan menilai "kira-kira sudah benar".

Nilai mengikuti **level tertinggi yang lolos test secara berurutan** (kalau Level 3 gagal, Level 8 tidak dihitung meski lolos) — kerjakan sejauh kemampuan. Kalian **boleh mengerjakan tidak berurutan** (mis. Level 3 dulu sebelum Level 1/2) — `go test -v` tetap menunjukkan level mana saja yang benar-benar sudah lolos apa adanya. Tapi untuk nilai akhir, tetap usahakan level dasar (terutama Level 1–3) beres duluan, karena Level 4 secara khusus mencocokkan isi kode dengan README Level 1 — tidak akan lolos kalau Level 1 belum diisi.

## Cara Kerja Folder Ini

File `main.go` dan `main_test.go` di folder ini **sudah ada** (ikut ter-fetch bersama folder `pertemuan-01/`) — tidak perlu disalin dari mana pun.

```bash
cd pertemuan-01
go mod init pertemuan01
go test ./... -v
```

Semua level akan **gagal (FAIL)** di awal — itu normal, kalian belum mengedit apa-apa. Jalankan ulang tiap kali selesai mengedit `main.go` atau `README.md` (di root repo) untuk melihat level mana yang sudah lolos.

**Jangan edit `main_test.go`** — perubahan di sana tidak berpengaruh saat penilaian (dosen menimpa ulang file ini dengan versi asli). Dosen memantau progres dengan menjalankan `go test` yang sama terhadap kode yang kalian push — tidak ada CI/GitHub Actions, jadi tidak perlu setup apa pun selain langkah di atas.

Level 1, 2, 6, 9, dan 10 diisi di **`README.md` milik repo kalian sendiri** (bukan file di folder ini) — heading yang dicek sudah disiapkan di sana, tinggal isi.

---

## Level 1 — Identitas di README

Isi bagian `## Identitas` di README dengan Nama, NRP (angka asli kalian), dan Kelas.

**Dicek otomatis:** `Nama:` dan `Kelas:` terisi (bukan teks contoh), `NRP:` berupa angka.

## Level 2 — Penjelasan Commit vs Push

Isi `## Commit vs Push`: jelaskan dengan bahasamu sendiri (minimal ±40 karakter) apa bedanya `git commit` dan `git push`, dan beri satu contoh situasi nyata seseorang commit tapi lupa push — apa akibatnya bagi rekan satu tim?

**Dicek otomatis:** section terisi memadai dan menyebut kata "commit" dan "push".

## Level 3 — Fungsi `ResolveNama`

Di `pertemuan-01/main.go`, implementasikan:
```go
func ResolveNama(args []string, fallback string) string
```
Kembalikan `args[0]` kalau `args` tidak kosong, kalau tidak kembalikan `fallback`.

**Dicek otomatis:** `ResolveNama(nil, "fallback")` harus `"fallback"`; `ResolveNama([]string{"Budi"}, "fallback")` harus `"Budi"`.

## Level 4 — Konstanta `NRP`

Ganti nilai `const NRP = "0000000000"` di `main.go` dengan NRP kalian sungguhan — **harus sama persis** dengan yang ditulis di README bagian Identitas.

**Dicek otomatis:** `NRP` berformat angka, dan nilainya harus ditemukan di dalam section `## Identitas` README.

## Level 5 — Fungsi `CetakInfo` (versi awal)

Implementasikan:
```go
func CetakInfo(nama string) string
```
Kembalikan string yang memuat `nama`, `NRP`, dan hasil `runtime.Version()`, misalnya:
```
Nama: Budi
NRP: 5025201012
go1.23.0
```

**Dicek otomatis:** `CetakInfo("Budi")` harus memuat `"Budi"`, memuat `NRP`, dan memuat pola versi Go (`go1.<angka>`).

## Level 6 — Penjelasan Reproducibility

Isi `## Reproducibility`: jelaskan dengan contoh konkret, apa yang bisa terjadi kalau anggota tim menjalankan program ini dengan versi Go yang berbeda — kapan itu jadi masalah nyata, kapan tidak?

**Dicek otomatis:** section terisi memadai dan menyinggung kata "versi"/"version".

Sampai di sini, jalankan `go run pertemuan-01/main.go Budi` dan `go run pertemuan-01/main.go` (tanpa argumen) — pastikan keduanya jalan tanpa crash.

## Level 7 — Fungsi `Sapa`

Implementasikan:
```go
func Sapa(nama string) string
```
Kembalikan kalimat sapaan yang memuat nama, misalnya `"Halo, Budi! Selamat datang di kelas Backend."`.

**Dicek otomatis:** `Sapa("Budi")` harus memuat `"Budi"` dan lebih panjang dari sekadar nama itu sendiri.

## Level 8 — Merge Conflict Sungguhan

Sekarang alami sendiri merge conflict, bukan cuma tahu istilahnya:

1. Buat branch baru dari branch utama (contoh: `fitur-sapaan`). Di branch ini, ubah `CetakInfo` supaya hasilnya **juga** memuat `Sapa(nama)` (mis. tambahkan baris sapaan di akhir output). Commit di branch ini.
2. **Sebelum merge**, kembali ke branch utama dan ubah baris `return` yang sama di `CetakInfo` dengan cara berbeda (mis. ubah formatnya) — sengaja bikin perubahan ini bentrok dengan langkah 1. Commit langsung di branch utama.
3. Merge branch `fitur-sapaan` ke branch utama — git akan melaporkan **conflict**. Selesaikan secara manual: hasil akhir `CetakInfo` harus tetap memuat Nama/NRP/versi Go **dan** sapaan dari `Sapa()`. Commit hasil resolusinya dengan pesan `pertemuan-01: level 8 selesai`.

**Dicek otomatis:** riwayat git punya minimal satu merge commit sungguhan (`git log --merges`), tidak ada sisa conflict marker (`<<<<<<<`/`=======`/`>>>>>>>`) di `main.go`, dan `CetakInfo(...)` benar-benar memuat hasil `Sapa(...)` setelah merge.

## Level 9 — Catatan Merge Conflict

Isi `## Catatan Merge Conflict`: baris mana yang bentrok, kenapa (dua sisi mengubah baris yang sama), dan bagaimana kamu memutuskan hasil akhirnya.

**Dicek otomatis:** section terisi memadai dan menyinggung soal conflict/konflik/bentrok.

## Level 10 — `.gitignore` dan Refleksi

Tambahkan `.gitignore` di root repo yang mengabaikan minimal: hasil build (`bin/`, `*.exe`, atau `*.out`), `.idea/`, dan `.vscode/`. Lalu isi dua bagian penutup README:
- **`## Kenapa .gitignore Penting`**: satu skenario konkret masalah kolaborasi tim kalau file build/IDE ikut ter-commit.
- **`## Refleksi`**: 2–3 kalimat jujur soal bagian mana dari setup pertemuan ini yang paling membingungkan, dan bagaimana akhirnya kamu memahaminya.

**Dicek otomatis:** `.gitignore` memuat pola yang diminta, kedua section README terisi memadai.
