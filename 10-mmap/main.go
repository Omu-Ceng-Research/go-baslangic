//go:build linux

// Ders 10 — mmap ile dosyayı belleğe eşlemek, read() ile karşılaştırmak
//
// Çalıştır: go run ./10-mmap                 (varsayılan 100.000 vektör, ~150 MB)
//
//	go run ./10-mmap -adet 300000    (daha büyük deneme)
//
// Bu ders projenin kalbidir: vektör deposunu böyle okuyacaksınız.
package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"syscall"
	"time"
	"unsafe"
)

const boyut = 384 // her vektör 384 float32 = 1536 bayt

func main() {
	adet := flag.Int("adet", 100_000, "vektör sayısı")
	yol := flag.String("dosya", "ders10-vektorler.f32", "vektör dosyası")
	flag.Parse()

	if _, err := os.Stat(*yol); err != nil {
		fmt.Printf("%s oluşturuluyor (%d vektör, %.0f MB)...\n",
			*yol, *adet, float64(*adet*boyut*4)/(1<<20))
		dosyaOlustur(*yol, *adet)
	}

	sorgu := rastgeleVektor(rand.New(rand.NewSource(42)))

	// Yöntem A: tüm dosyayı read(2) ile heap'e kopyala, sonra ara.
	// Yöntem B: mmap(2) ile dosyayı adres alanına eşle, kopyasız ara.
	olc("A) os.ReadFile (read)", func() (int, float32) { return readIleAra(*yol, sorgu) })
	olc("B) mmap", func() (int, float32) { return mmapIleAra(*yol, sorgu) })

	fmt.Println()
	fmt.Println("Soğuk önbellek denemesi için (root gerekir):")
	fmt.Println("  sync; echo 3 | sudo tee /proc/sys/vm/drop_caches")
	fmt.Println("ardından programı tekrar çalıştırın ve 'major fault' sayısına bakın.")
}

// olc: süreyi ve sayfa hatalarını (page fault) getrusage(2) ile ölçer.
func olc(ad string, f func() (int, float32)) {
	var once, sonra syscall.Rusage
	syscall.Getrusage(syscall.RUSAGE_SELF, &once)
	basla := time.Now()

	id, skor := f()

	sure := time.Since(basla)
	syscall.Getrusage(syscall.RUSAGE_SELF, &sonra)
	fmt.Printf("%-24s en yakın id=%d skor=%.4f süre=%v minor fault=%d major fault=%d\n",
		ad, id, skor, sure.Round(time.Millisecond),
		sonra.Minflt-once.Minflt, sonra.Majflt-once.Majflt)
}

func readIleAra(yol string, sorgu []float32) (int, float32) {
	ham, err := os.ReadFile(yol) // tüm dosya heap'e kopyalanır
	if err != nil {
		panic(err)
	}
	return enYakin(baytlariFloataCevir(ham), sorgu)
}

func mmapIleAra(yol string, sorgu []float32) (int, float32) {
	f, err := os.Open(yol)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	bilgi, _ := f.Stat()

	// mmap(2): dosyayı salt okunur ve paylaşımlı olarak belleğe eşle.
	// Henüz hiçbir bayt okunmadı! Sayfalar ancak dokunulduğunda yüklenir (page fault).
	ham, err := syscall.Mmap(int(f.Fd()), 0, int(bilgi.Size()),
		syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	defer syscall.Munmap(ham)

	// madvise(2): çekirdeğe "baştan sona sırayla okuyacağım" ipucu ver.
	// Çekirdek buna göre ileriye dönük okuma (readahead) yapar.
	syscall.Madvise(ham, syscall.MADV_SEQUENTIAL)

	return enYakin(baytlariFloataCevir(ham), sorgu)
}

// []byte'ı kopyalamadan []float32 olarak görmek.
// unsafe: aynı bellek bölgesine farklı tipte bakıyoruz (C'deki (float*)ptr cast'i).
// Not: makinenin little-endian olduğunu varsayar (x86-64 ve ARM64 öyledir).
func baytlariFloataCevir(b []byte) []float32 {
	return unsafe.Slice((*float32)(unsafe.Pointer(&b[0])), len(b)/4)
}

// enYakin: brute-force arama. Her vektörle iç çarpım (dot product) hesaplar.
func enYakin(tum []float32, sorgu []float32) (int, float32) {
	enIyiID, enIyiSkor := -1, float32(math.Inf(-1))
	adet := len(tum) / boyut
	for i := range adet {
		v := tum[i*boyut : (i+1)*boyut]
		var skor float32
		for j := range boyut {
			skor += v[j] * sorgu[j]
		}
		if skor > enIyiSkor {
			enIyiID, enIyiSkor = i, skor
		}
	}
	return enIyiID, enIyiSkor
}

func dosyaOlustur(yol string, adet int) {
	f, err := os.Create(yol)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriterSize(f, 1<<20)
	r := rand.New(rand.NewSource(1))
	tampon := make([]byte, 4)
	for range adet {
		for _, x := range rastgeleVektor(r) {
			binary.LittleEndian.PutUint32(tampon, math.Float32bits(x))
			w.Write(tampon)
		}
	}
	if err := w.Flush(); err != nil {
		panic(err)
	}
}

// Birim uzunlukta rastgele vektör (iç çarpım = kosinüs benzerliği olsun diye).
func rastgeleVektor(r *rand.Rand) []float32 {
	v := make([]float32, boyut)
	var norm float64
	for i := range v {
		v[i] = float32(r.NormFloat64())
		norm += float64(v[i] * v[i])
	}
	n := float32(math.Sqrt(norm))
	for i := range v {
		v[i] /= n
	}
	return v
}

// SORULAR:
//  1) İki yöntem aynı id'yi buluyor mu? Bulmalı. Bulmuyorsa neden?
//  2) Hangi yöntemde minor fault daha çok? Neden? (İpucu: ReadFile'ın ayırdığı
//     yeni heap belleği de ilk dokunuşta page fault üretir.)
//  3) /usr/bin/time -v ./d10 ile "Maximum resident set size" değerine bakın.
//     Dosyadan büyük bir RAM'e sahip olmasaydınız hangi yöntem hâlâ çalışırdı?
//  4) strace -e trace=mmap,munmap,madvise,read ./d10 çıktısında kendi mmap
//     çağrınızı bulun. Boyut argümanı dosya boyutuyla aynı mı?
