// Ders 07 — İşaretçiler ve bellek
//
// Çalıştır: go run ./07-pointer-bellek
//
// C'deki gibi & ve * vardır ama:
//   - işaretçi aritmetiği YOKTUR (p++ yazamazsınız),
//   - free() YOKTUR, çöp toplayıcı (GC) vardır,
//   - yerel değişkenin adresini döndürmek GÜVENLİDİR (derleyici onu heap'e taşır).
package main

import (
	"fmt"
	"runtime"
	"unsafe"
)

type Vektor struct {
	ID    uint64
	Deger [4]float32
}

// Değerle geçirme: struct kopyalanır, değişiklik dışarıya yansımaz.
func normallestirKopya(v Vektor) {
	v.Deger[0] = 0
}

// İşaretçiyle geçirme: asıl struct değişir.
func normallestir(v *Vektor) {
	v.Deger[0] = 0 // (*v).Deger[0] yazmaya gerek yok, Go otomatik çözer
}

// C'de bu fonksiyon tanımsız davranış (dangling pointer) olurdu. Go'da güvenli.
func yeniVektor(id uint64) *Vektor {
	v := Vektor{ID: id}
	return &v
}

func main() {
	x := 42
	p := &x
	*p = 100
	fmt.Println("x =", x, "p =", p)

	v := Vektor{ID: 1, Deger: [4]float32{9, 9, 9, 9}}
	normallestirKopya(v)
	fmt.Println("kopya ile:   ", v.Deger)
	normallestir(&v)
	fmt.Println("işaretçi ile:", v.Deger)

	w := yeniVektor(7)
	fmt.Println("yeniVektor:", w.ID)

	// Bellek düzeni: bir struct kaç bayt? (C'deki sizeof)
	fmt.Println("sizeof(Vektor) =", unsafe.Sizeof(Vektor{}), "bayt")
	fmt.Println("sizeof(float32) =", unsafe.Sizeof(float32(0)), "bayt")

	// 384 boyutlu 1 milyon vektör ne kadar yer kaplar?
	const boyut, adet = 384, 1_000_000
	fmt.Printf("384 x 1M float32 = %.2f GB\n", float64(boyut*adet*4)/(1<<30))

	// Çöp toplayıcıyı gözlemleyelim.
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("başlangıç heap: %d KB\n", m.HeapAlloc/1024)

	buyuk := make([]float32, 50_000_000) // ~200 MB
	buyuk[0] = 1
	runtime.ReadMemStats(&m)
	fmt.Printf("200 MB ayırdıktan sonra heap: %d MB\n", m.HeapAlloc/(1<<20))

	buyuk = nil // referansı bıraktık
	runtime.GC()
	runtime.ReadMemStats(&m)
	fmt.Printf("GC sonrası heap: %d KB, GC sayısı: %d\n", m.HeapAlloc/1024, m.NumGC)
}

// SORU: 1.5 GB'lık vektörleri heap'te tutmak yerine dosyayı mmap ile belleğe
// eşlersek GC bu belleği tarar mı? (İpucu: ders 10)
