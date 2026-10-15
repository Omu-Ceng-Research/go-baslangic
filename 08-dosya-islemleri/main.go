// Ders 08 — Dosya işlemleri: metin, satır satır okuma, buffered/unbuffered, ikili veri
//
// Çalıştır: go run ./08-dosya-islemleri
// Çıktılar ./ders08-cikti/ klasörüne yazılır.
package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"time"
)

const klasor = "ders08-cikti"

func main() {
	if err := os.MkdirAll(klasor, 0o755); err != nil {
		panic(err)
	}

	metinDosyasi()
	satirSatirOku()
	bufferedKarsilastir()
	ikiliVektorYazOku()
}

// [1] En basit yol: tüm dosyayı tek seferde yaz/oku.
func metinDosyasi() {
	yol := filepath.Join(klasor, "merhaba.txt")
	if err := os.WriteFile(yol, []byte("SET 1 merhaba\nSET 2 dünya\nGET 1\n"), 0o644); err != nil {
		panic(err)
	}
	icerik, err := os.ReadFile(yol)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[1] %s içeriği (%d bayt):\n%s", yol, len(icerik), icerik)
}

// [2] Büyük dosyayı belleğe almadan satır satır okumak (Ruby'deki File.foreach).
func satirSatirOku() {
	f, err := os.Open(filepath.Join(klasor, "merhaba.txt"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	tarayici := bufio.NewScanner(f)
	no := 0
	for tarayici.Scan() {
		no++
		fmt.Printf("[2] satır %d: %s\n", no, tarayici.Text())
	}
	if err := tarayici.Err(); err != nil {
		panic(err)
	}
}

// [3] Buffered ve unbuffered yazma. Aynı veriyi iki şekilde yazıyoruz.
// Fark nereden geliyor? -> write(2) sistem çağrısı sayısından.
func bufferedKarsilastir() {
	const kayitSayisi = 100_000
	kayit := []byte("SET 123456 bu bir deneme mesajıdır\n")

	// Unbuffered: her f.Write doğrudan bir write(2) sistem çağrısıdır.
	basla := time.Now()
	f, err := os.Create(filepath.Join(klasor, "unbuffered.log"))
	if err != nil {
		panic(err)
	}
	for range kayitSayisi {
		f.Write(kayit)
	}
	f.Close()
	sureUnbuf := time.Since(basla)

	// Buffered: bufio 4 KB'lık bir tampon doldurur, dolunca tek write(2) yapar.
	basla = time.Now()
	f, err = os.Create(filepath.Join(klasor, "buffered.log"))
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)
	for range kayitSayisi {
		w.Write(kayit)
	}
	w.Flush() // UNUTMAYIN: Flush yoksa tampondaki son veri diske gitmez!
	f.Close()
	sureBuf := time.Since(basla)

	fmt.Printf("[3] %d kayıt -> unbuffered: %v, buffered: %v (%.1fx)\n",
		kayitSayisi, sureUnbuf, sureBuf, float64(sureUnbuf)/float64(sureBuf))
}

// [4] İkili (binary) veri: vektörleri metin olarak değil, ham float32 olarak yazmak.
// Dosya biçimi: [adet uint32][boyut uint32][adet*boyut adet float32], little-endian
func ikiliVektorYazOku() {
	const adet, boyut = 1000, 8
	yol := filepath.Join(klasor, "vektorler.bin")

	f, err := os.Create(yol)
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)
	binary.Write(w, binary.LittleEndian, uint32(adet))
	binary.Write(w, binary.LittleEndian, uint32(boyut))
	for i := range adet {
		for j := range boyut {
			binary.Write(w, binary.LittleEndian, float32(math.Sin(float64(i*boyut+j))))
		}
	}
	w.Flush()
	f.Close()

	bilgi, _ := os.Stat(yol)
	fmt.Printf("[4] %s: %d bayt (beklenen: 8 + %d*%d*4 = %d)\n",
		yol, bilgi.Size(), adet, boyut, 8+adet*boyut*4)

	// Geri okuma: 500. vektöre doğrudan atlamak için Seek kullanıyoruz.
	// Metin dosyasında bunu yapamazdık: satırların uzunluğu farklıdır.
	f, err = os.Open(yol)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	const hedef = 500
	ofset := int64(8 + hedef*boyut*4)
	if _, err := f.Seek(ofset, io.SeekStart); err != nil { // lseek(2)
		panic(err)
	}
	vektor := make([]float32, boyut)
	if err := binary.Read(f, binary.LittleEndian, vektor); err != nil {
		panic(err)
	}
	fmt.Printf("[4] %d. vektör: %.3f\n", hedef, vektor)
}

// DENEYİN (Linux/WSL2):
//   go build -o d08 ./08-dosya-islemleri
//   strace -c -e trace=write ./d08
// Kaç write çağrısı görüyorsunuz? 100.000 civarı mı, yoksa çok daha az mı?
// Hangi kısım hangisinden geliyor?
