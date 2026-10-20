# Pertemuan 2 — Sales Order Processor: Fungsi & Alur Keputusan

Tugas ini melatih kemampuan **merancang dan menggunakan fungsi** di Go, dengan studi kasus memproses pesanan penjualan toko (lanjutan dari contoh Sales Order Processor di materi "Dasar Go"). Konteksnya berdiri sendiri — bukan bagian dari proyek Task Manager yang kalian kerjakan di pertemuan lain.

**Soal ini sengaja tidak menunjukkan kode/algoritma jadi.** Tiap level hanya menjelaskan *apa* yang harus dilakukan fungsi tersebut (kontrak: nama, input, output, aturan, contoh) — bagaimana cara menulisnya adalah bagian yang harus kalian pikirkan sendiri. Untuk dua level (9 dan 10), signature fungsinya sudah benar di starter `main.go`, tapi bentuknya belum pernah dibahas di kelas — cari tahu sendiri alasannya sebelum mengisi isinya (dokumentasi resmi Go di [go.dev](https://go.dev) adalah tempat yang bagus untuk mulai).

Semua level dicek otomatis lewat `go test`. Nilai mengikuti level tertinggi yang lolos **berurutan** dari 1 — kerjakan sejauh kemampuan, boleh tidak berurutan (lihat `go test -v` untuk status tiap level apa adanya).

---

## Level 1 — `HitungSubtotal`

Buat fungsi bernama `HitungSubtotal` yang menerima **jumlah barang** (bilangan bulat) dan **harga satuan** (bilangan desimal), lalu mengembalikan hasil kali keduanya sebagai bilangan desimal.

Contoh: jumlah 3, harga satuan 15000 → kembalikan 45000.

## Level 2 — `HitungTotalPesanan`

Buat fungsi bernama `HitungTotalPesanan` yang menerima **daftar jumlah barang** dan **daftar harga satuan** (mewakili beberapa item dalam satu pesanan — item ke-i punya jumlah dan harga di indeks ke-i pada masing-masing daftar), lalu mengembalikan total keseluruhan pesanan (jumlah dari subtotal tiap item).

Kalau banyaknya jumlah barang dan banyaknya harga satuan **tidak sama**, anggap pesanan tidak bisa diproses dan kembalikan 0.

Contoh: jumlah `[2, 5, 1]`, harga `[15000, 8000, 25000]` → kembalikan 95000 (= 2×15000 + 5×8000 + 1×25000).

## Level 3 — `TerapkanPajak`

Buat fungsi bernama `TerapkanPajak` yang menerima **total belanja** dan **tarif pajak** (mis. 0.11 untuk PPN 11%), lalu mengembalikan total **setelah** ditambah pajaknya.

Contoh: total 100000, tarif 0.11 → kembalikan 111000.

## Level 4 — `HitungDiskon`

Buat fungsi bernama `HitungDiskon` yang menerima **total belanja**, lalu mengembalikan **nominal diskonnya** (bukan total setelah diskon) berdasarkan aturan bertingkat:
- Total ≥ 1.000.000 → diskon 10% dari total.
- Total ≥ 500.000 (tapi < 1.000.000) → diskon 5% dari total.
- Selain itu → tidak ada diskon (0).

Contoh: total 1.200.000 → kembalikan 120.000. Total 600.000 → kembalikan 30.000. Total 200.000 → kembalikan 0.

## Level 5 — `TotalSetelahDiskon`

Buat fungsi bernama `TotalSetelahDiskon` yang menerima **daftar jumlah barang**, **daftar harga satuan**, dan **tarif pajak**, lalu menghitung alur lengkap satu pesanan dari awal sampai akhir: jumlahkan semua item jadi subtotal, kurangi dengan nominal diskon yang berlaku untuk subtotal itu, **baru** kenakan pajak pada hasil setelah diskon. Kembalikan total akhirnya.

Manfaatkan kembali fungsi-fungsi yang sudah kalian buat di level sebelumnya — jangan tulis ulang logikanya dari nol.

Contoh: jumlah `[10, 10]`, harga `[60000, 60000]`, tarif pajak 0.11 → subtotal 1.200.000 → diskon 10% (120.000) → setelah diskon 1.080.000 → setelah pajak 11% → kembalikan 1.198.800.

## Level 6 — `ValidasiPesanan`

Buat fungsi bernama `ValidasiPesanan` yang menerima **daftar jumlah barang** dan **daftar harga satuan**, lalu mengembalikan **dua nilai sekaligus**: sebuah penanda benar/salah (apakah pesanan valid), dan sebuah teks pesan.

Pesanan dianggap **tidak valid** kalau salah satu dari ini terjadi: banyaknya jumlah barang dan harga satuan tidak sama, ATAU ada jumlah barang yang bukan bilangan positif, ATAU ada harga satuan yang bukan bilangan positif.

- Kalau valid: penanda benar, teks pesan **kosong**.
- Kalau tidak valid: penanda salah, teks pesan **tidak boleh kosong** (isi bebas, jelaskan alasannya dengan kalimatmu sendiri — dicek otomatis hanya bahwa pesannya ada isinya, bukan kata-kata persisnya).

Ini adalah pertama kalinya kalian membuat fungsi yang mengembalikan lebih dari satu nilai sekaligus — kalau belum tahu caranya, cari tahu dulu.

## Level 7 — `TentukanStatus`

Buat fungsi bernama `TentukanStatus` yang menerima **total akhir pesanan**, lalu mengembalikan salah satu dari tiga teks berikut:
- `"Prioritas"` kalau total lebih dari 1.000.000.
- `"Reguler"` kalau total lebih dari 100.000 (tapi tidak lebih dari 1.000.000).
- `"Hemat"` untuk selain itu.

## Level 8 — `RingkasanPesanan`

Buat fungsi bernama `RingkasanPesanan` yang menerima **daftar jumlah barang**, **daftar harga satuan**, dan **tarif pajak**, lalu mengembalikan **satu teks ringkasan** (boleh multi-baris, format bebas — kalian yang desain) yang nantinya bisa langsung dicetak ke layar oleh `main()`.

Ringkasan itu **wajib memuat**, dalam bentuk angka yang bisa dibaca:
1. Subtotal pesanan (sebelum diskon dan pajak),
2. Nominal diskon yang didapat,
3. Total akhir (setelah diskon dan pajak), dan
4. Status pesanannya.

Manfaatkan kembali fungsi-fungsi yang sudah kalian buat sebelumnya untuk menghitung keempat hal itu — fungsi ini hanya bertugas merangkainya jadi satu teks.

## Level 9 — `Total` (bentuk fungsi baru)

Di starter `main.go`, perhatikan signature fungsi `Total` — bentuknya berbeda dari fungsi-fungsi sebelumnya. Fungsi ini harus bisa **dipanggil dengan berapa pun banyaknya angka**, termasuk **tanpa argumen sama sekali** (kembalikan 0 kalau begitu), dan mengembalikan jumlah dari semua angka yang diberikan.

Contoh: dipanggil tanpa argumen → 0. Dipanggil dengan satu angka (5) → 5. Dipanggil dengan lima angka (1, 2, 3, 4, 5) → 15.

Cari tahu: bentuk parameter seperti apa yang memungkinkan sebuah fungsi Go menerima jumlah argumen yang tidak tetap seperti ini? Ini tidak dibahas di kelas — dokumentasi resmi Go adalah tempat yang bagus untuk mulai mencari.

## Level 10 — `HitungOngkosKirim` (Bonus)

Di starter `main.go`, perhatikan juga signature fungsi `HitungOngkosKirim` — fungsi ini mengembalikan **dua nilai**: ongkos kirimnya, dan sebuah nilai bertipe khusus yang dipakai Go untuk menandai kegagalan.

Aturan:
- Kalau **berat ≤ 0** atau **jarak < 0**: input dianggap tidak valid — nilai kegagalan pada return kedua harus **terisi** (menandakan ada masalah); nilai ongkos pada kasus ini tidak dicek.
- Kalau input valid: ongkos dihitung dengan rumus `(berat × 2000) + (jarak × 3000)`, dan nilai kegagalan pada return kedua harus **kosong/tidak ada** (menandakan tidak ada masalah).

Cari tahu: tipe apa yang dipakai Go untuk merepresentasikan "mungkin ada kegagalan, mungkin tidak" sebagai nilai balik kedua seperti ini, dan bagaimana cara membuat nilai yang menandakan "ada kegagalan" dari sebuah pesan teks?
