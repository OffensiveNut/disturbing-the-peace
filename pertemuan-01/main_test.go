// File ini disediakan dosen untuk mengecek progres level secara otomatis.
// JANGAN DIUBAH — perubahan pada file ini tidak akan dipakai saat penilaian
// (dosen menimpa ulang file ini sebelum menjalankan grading).
package main

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

const placeholder = "(tulis di sini)"

var (
	headingRe = regexp.MustCompile(`(?im)^##\s+(.+?)\s*$`)
	nrpLineRe = regexp.MustCompile(`NRP:\s*(\d{6,12})`)
)

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

// section mengembalikan isi di bawah heading "## <heading>" sampai heading
// berikutnya (atau akhir file), atau string kosong kalau heading tidak ada.
func section(readme, heading string) string {
	locs := headingRe.FindAllStringSubmatchIndex(readme, -1)
	headings := headingRe.FindAllStringSubmatch(readme, -1)
	for i, h := range headings {
		if strings.EqualFold(strings.TrimSpace(h[1]), heading) {
			start := locs[i][1]
			end := len(readme)
			if i+1 < len(locs) {
				end = locs[i+1][0]
			}
			return strings.TrimSpace(readme[start:end])
		}
	}
	return ""
}

// filled menganggap section "belum dikerjakan" kalau kosong, masih persis
// teks placeholder soal, atau lebih pendek dari minLen karakter.
func filled(s string, minLen int) bool {
	s = strings.TrimSpace(s)
	if s == "" || strings.EqualFold(s, placeholder) {
		return false
	}
	return len(s) >= minLen
}

func TestLevels(t *testing.T) {
	readme := readFile(t, "../README.md")

	t.Run("Level 1 - Identitas di README", func(t *testing.T) {
		ident := section(readme, "Identitas")
		if ident == "" {
			t.Fatal("section '## Identitas' tidak ditemukan di README.md")
		}
		if !nrpLineRe.MatchString(ident) {
			t.Error("NRP belum diisi dengan format angka yang benar, mis. 'NRP: 5025201012'")
		}
		namaRe := regexp.MustCompile(`(?i)Nama:\s*(.+)`)
		if m := namaRe.FindStringSubmatch(ident); m == nil || strings.TrimSpace(m[1]) == "" ||
			strings.Contains(strings.ToLower(m[1]), "nama lengkap") {
			t.Error("Nama belum diisi dengan nama asli")
		}
		kelasRe := regexp.MustCompile(`(?i)Kelas:\s*(.+)`)
		if m := kelasRe.FindStringSubmatch(ident); m == nil || strings.TrimSpace(m[1]) == "" {
			t.Error("Kelas belum diisi")
		}
	})

	t.Run("Level 2 - Penjelasan Commit vs Push", func(t *testing.T) {
		s := section(readme, "Commit vs Push")
		if !filled(s, 40) {
			t.Fatalf("section '## Commit vs Push' belum diisi memadai (minimal 40 karakter, saat ini %d)", len(strings.TrimSpace(s)))
		}
		low := strings.ToLower(s)
		if !strings.Contains(low, "commit") || !strings.Contains(low, "push") {
			t.Error("penjelasan harus menyebut kata 'commit' dan 'push'")
		}
	})

	t.Run("Level 3 - ResolveNama berfungsi benar", func(t *testing.T) {
		if got := ResolveNama(nil, "fallback"); got != "fallback" {
			t.Errorf(`ResolveNama(nil, "fallback") = %q, ingin "fallback"`, got)
		}
		if got := ResolveNama([]string{"Budi"}, "fallback"); got != "Budi" {
			t.Errorf(`ResolveNama([]string{"Budi"}, "fallback") = %q, ingin "Budi"`, got)
		}
	})

	t.Run("Level 4 - NRP di kode cocok dengan README", func(t *testing.T) {
		if !regexp.MustCompile(`^\d{6,12}$`).MatchString(NRP) {
			t.Fatal("konstanta NRP di main.go belum diganti dari nilai default")
		}
		ident := section(readme, "Identitas")
		if !strings.Contains(ident, NRP) {
			t.Errorf("NRP di main.go (%q) tidak sama dengan NRP di README.md bagian Identitas", NRP)
		}
	})

	t.Run("Level 5 - CetakInfo mencetak versi Go", func(t *testing.T) {
		out := CetakInfo("Budi")
		if !strings.Contains(out, "Budi") {
			t.Error("CetakInfo harus memuat nama yang diberikan")
		}
		if !strings.Contains(out, NRP) {
			t.Error("CetakInfo harus memuat NRP")
		}
		if !regexp.MustCompile(`go1\.\d+`).MatchString(out) {
			t.Error("CetakInfo harus memuat hasil runtime.Version(), mis. 'go1.23.0'")
		}
	})

	t.Run("Level 6 - Penjelasan Reproducibility", func(t *testing.T) {
		s := section(readme, "Reproducibility")
		if !filled(s, 40) {
			t.Fatal("section '## Reproducibility' belum diisi memadai (minimal 40 karakter)")
		}
		low := strings.ToLower(s)
		if !strings.Contains(low, "versi") && !strings.Contains(low, "version") {
			t.Error("penjelasan harus menyinggung soal versi Go")
		}
	})

	t.Run("Level 7 - Fungsi Sapa", func(t *testing.T) {
		got := Sapa("Budi")
		if !strings.Contains(got, "Budi") {
			t.Fatalf(`Sapa("Budi") = %q, harus memuat nama yang diberikan`, got)
		}
		if len(got) <= len("Budi")+5 {
			t.Error("Sapa harus mengembalikan kalimat sapaan, bukan sekadar nama")
		}
	})

	t.Run("Level 8 - Merge conflict sungguhan & terintegrasi", func(t *testing.T) {
		out, err := exec.Command("git", "log", "--merges", "--oneline").Output()
		if err != nil {
			t.Fatalf("gagal menjalankan git log: %v", err)
		}
		if strings.TrimSpace(string(out)) == "" {
			t.Fatal("tidak ditemukan merge commit di riwayat git - branch belum di-merge lewat proses merge sungguhan")
		}
		mainGo := readFile(t, "main.go")
		for _, marker := range []string{"<<<<<<<", "=======", ">>>>>>>"} {
			if strings.Contains(mainGo, marker) {
				t.Errorf("masih ada sisa conflict marker %q di main.go", marker)
			}
		}
		if info := CetakInfo("Budi"); !strings.Contains(info, Sapa("Budi")) {
			t.Error("hasil CetakInfo harus menyertakan sapaan dari fungsi Sapa setelah merge")
		}
	})

	t.Run("Level 9 - Catatan Merge Conflict", func(t *testing.T) {
		s := section(readme, "Catatan Merge Conflict")
		if !filled(s, 40) {
			t.Fatal("section '## Catatan Merge Conflict' belum diisi memadai (minimal 40 karakter)")
		}
		low := strings.ToLower(s)
		keywords := []string{"conflict", "konflik", "bentrok", "tabrak"}
		found := false
		for _, k := range keywords {
			if strings.Contains(low, k) {
				found = true
				break
			}
		}
		if !found {
			t.Error("penjelasan harus menyinggung soal conflict/konflik/bentrok antar perubahan")
		}
	})

	t.Run("Level 10 - .gitignore dan refleksi", func(t *testing.T) {
		gi := readFile(t, "../.gitignore")
		if strings.TrimSpace(gi) == "" {
			t.Fatal(".gitignore tidak ditemukan atau kosong")
		}
		lowGi := strings.ToLower(gi)
		for _, n := range []string{".idea", ".vscode"} {
			if !strings.Contains(lowGi, n) {
				t.Errorf(".gitignore belum mengabaikan %q", n)
			}
		}
		if !regexp.MustCompile(`(?i)bin/|\.exe|\.out`).MatchString(gi) {
			t.Error(".gitignore belum mengabaikan hasil build (mis. bin/, *.exe, atau *.out)")
		}
		if !filled(section(readme, "Kenapa .gitignore Penting"), 40) {
			t.Error("section '## Kenapa .gitignore Penting' belum diisi memadai")
		}
		if !filled(section(readme, "Refleksi"), 40) {
			t.Error("section '## Refleksi' belum diisi memadai")
		}
	})
}
