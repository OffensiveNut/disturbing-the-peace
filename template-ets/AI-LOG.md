# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

## Implementasi awal B1-B3 (sumber daya dan peminjaman)
- Prompt: berikan contoh implementasi untuk B1-B3 beserta tasklist, berdasarkan dokumen Requirement ETS Backend PinjamKampus.
- Dipakai: `createSumberDaya`, `listSumberDaya`, `getSumberDaya`, `createPeminjaman`, `listPeminjaman`, `getPeminjaman`, `writeJSON`, `decodeBody`, dan `idFromPath` di `api.go`.
- Diubah: menyesuaikan nama field JSON, aturan validasi, kode status, format error `{"error":"..."}`, dan ID yang dimulai dari 1. Struktur awal juga diubah agar memakai slice di memori tanpa database.

## Implementasi B4-B7 dan fitur persona E
- Prompt: bagaimana implementasi B4-B7 dan fitur persona E berupa blokir jadwal.
- Dipakai: `validasiPeminjaman`, `cekAturanPeminjaman`, `bentrok`, `kenaBlokir`, `updatePeminjaman`, `deletePeminjaman`, `createBlokir`, `listBlokir`, `getBlokir`, dan `deleteBlokir` di `api.go`.
- Diubah: cek bentrok diubah menjadi per slot dengan menjumlahkan `jumlah` peminjaman aktif, bukan hanya mengecek overlap rentang. PUT mengecualikan ID peminjaman yang sedang diedit. Blokir baru langsung membatalkan peminjaman aktif yang overlap. Semua akses data dilindungi `sync.Mutex`.

## Change request A: GET /ringkasan
- Prompt: tambahkan endpoint `GET /ringkasan?tanggal=YYYY-MM-DD` sesuai change request A, dengan `jumlah_peminjaman`, `total_unit_slot`, dan `per_sumber_daya`.
- Dipakai: `getRingkasan`, `ringkasanResponse`, `ringkasanSumberDaya`, dan `menahanSlot` di `api.go`.
- Diubah: 

## Penyederhanaan api.go
- Prompt: bagaimana cara menysederhanakan seluruh `api.go`, hapus komentar, dan kurangi jumlah baris tanpa mengubah perilaku.
- Dipakai: handler dan helper di `api.go`.
- Diubah: validasi dipusatkan di `validasiPeminjaman` dan `validasiBlokir`; pengecekan aturan peminjaman dipusatkan di `cekAturanPeminjaman`; parsing ID, penulisan JSON, dan pengecekan overlap dijadikan helper. 

## Pengujian otomatis
- Prompt: jalankan semua tes sendiri dan berikan hasilnya.
- Dipakai: `api_test.go` yang menguji B1, B2-B3, B5, B6, B7, persona E, dan change request.
- Diubah: test awal untuk kasus bentrok B5 memakai slot 4 pada stok 2; ternyata slot itu baru terpakai 1 unit sehingga tidak bentrok. Ekspektasi test diperbaiki menjadi slot 3 yang memang terpakai 2 unit. Hasil akhir: 7/7 test lolos dengan `go test -race -v ./...`.
