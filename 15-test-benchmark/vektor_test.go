package vektor

import (
	"math"
	"math/rand"
	"testing"
)

// Tablo güdümlü test: Go'da en yaygın test kalıbı.
func TestIcCarpim(t *testing.T) {
	durumlar := []struct {
		ad   string
		a, b []float32
		bekl float32
	}{
		{"boş", []float32{}, []float32{}, 0},
		{"tek eleman", []float32{3}, []float32{4}, 12},
		{"dik vektörler", []float32{1, 0}, []float32{0, 1}, 0},
		{"5 eleman (4'e bölünmez)", []float32{1, 2, 3, 4, 5}, []float32{1, 1, 1, 1, 1}, 15},
	}
	for _, d := range durumlar {
		t.Run(d.ad, func(t *testing.T) {
			if got := IcCarpim(d.a, d.b); got != d.bekl {
				t.Errorf("IcCarpim = %v, beklenen %v", got, d.bekl)
			}
			if got := IcCarpim4(d.a, d.b); got != d.bekl {
				t.Errorf("IcCarpim4 = %v, beklenen %v", got, d.bekl)
			}
			if got := IcCarpim4SK(d.a, d.b); got != d.bekl {
				t.Errorf("IcCarpim4SK = %v, beklenen %v", got, d.bekl)
			}
		})
	}
}

func TestNormallestir(t *testing.T) {
	v := []float32{3, 4}
	Normallestir(v)
	if uzunluk := IcCarpim(v, v); math.Abs(float64(uzunluk-1)) > 1e-6 {
		t.Errorf("normalleştirilmiş vektörün uzunluğu 1 olmalı, %v bulundu", uzunluk)
	}

	sifir := []float32{0, 0}
	Normallestir(sifir) // panik yapmamalı, NaN üretmemeli
	if sifir[0] != 0 {
		t.Errorf("sıfır vektör değişmemeli, %v bulundu", sifir)
	}
}

// İki gerçeklemenin rastgele girdilerde aynı sonucu verdiğini kontrol et.
// float32 toplama sırası değişince sonuç birkaç bit farklı çıkabilir!
func TestGerceklemelerTutarli(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	a, b := rastgele(r, 384), rastgele(r, 384)
	x, y, z := IcCarpim(a, b), IcCarpim4(a, b), IcCarpim4SK(a, b)
	if math.Abs(float64(x-y)) > 1e-4 || math.Abs(float64(x-z)) > 1e-4 {
		t.Errorf("tutarsız: %v, %v, %v", x, y, z)
	}
}

func rastgele(r *rand.Rand, n int) []float32 {
	v := make([]float32, n)
	for i := range v {
		v[i] = float32(r.NormFloat64())
	}
	return v
}

// Benchmark: Go, b.N değerini süre anlamlı olana kadar kendisi büyütür.
func BenchmarkIcCarpim(b *testing.B) {
	r := rand.New(rand.NewSource(1))
	x, y := rastgele(r, 384), rastgele(r, 384)
	b.ResetTimer()
	for range b.N {
		IcCarpim(x, y)
	}
}

func BenchmarkIcCarpim4(b *testing.B) {
	r := rand.New(rand.NewSource(1))
	x, y := rastgele(r, 384), rastgele(r, 384)
	b.ResetTimer()
	for range b.N {
		IcCarpim4(x, y)
	}
}

func BenchmarkIcCarpim4SK(b *testing.B) {
	r := rand.New(rand.NewSource(1))
	x, y := rastgele(r, 384), rastgele(r, 384)
	b.ResetTimer()
	for range b.N {
		IcCarpim4SK(x, y)
	}
}

// SORULAR:
//  1) IcCarpim4, IcCarpim'den hızlı mı çıktı? Beklentiniz neydi?
//     "Optimizasyon" sanılan bir değişiklik neden yavaşlatmış olabilir?
//  2) IcCarpim4SK ile IcCarpim4 arasındaki tek fark nedir? Sonucu ne kadar değiştirdi?
//  3) Benchmark'ı -count 5 ile çalıştırın. Aynı fonksiyon için sonuçlar ne
//     kadar oynuyor? Tek bir ölçüme güvenmek neden tehlikeli?
//  4) Benchmark çıktısındaki ns/op değerini kullanarak 1 milyon vektör
//     üzerinde tek bir brute-force sorgunun kaç milisaniye süreceğini tahmin
//     edin. Ders 10'daki ölçümle karşılaştırın. Fark nereden geliyor?
