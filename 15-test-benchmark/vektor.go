// Ders 15 — Birim testi ve benchmark
//
// Çalıştır:
//
//	go test ./15-test-benchmark                       (testler)
//	go test -v ./15-test-benchmark                    (ayrıntılı)
//	go test -bench . -benchmem ./15-test-benchmark    (benchmark)
//
// Go'da test için ayrı bir kütüphane (JUnit, RSpec) gerekmez. Dosya adı
// _test.go ile biten dosyalar "go test" ile çalışır.
// Projede ölçüm tablolarınızı büyük ölçüde benchmark'larla üreteceksiniz.
//
// Bu klasör bir "main" paketi değil, bir kütüphane paketidir.
package vektor

import "math"

// IcCarpim, iki vektörün iç çarpımını (dot product) hesaplar.
func IcCarpim(a, b []float32) float32 {
	var t float32
	for i := range a {
		t += a[i] * b[i]
	}
	return t
}

// IcCarpim4, aynı işi döngüyü 4'lü açarak (loop unrolling) yapar.
// Sezgi: işlemci 4 bağımsız toplamı aynı anda yürütebilir, daha hızlı olmalı.
// Gerçekten öyle mi? Tahmin etmeyin, benchmark söylesin.
func IcCarpim4(a, b []float32) float32 {
	var t0, t1, t2, t3 float32
	n := len(a) &^ 3 // 4'ün katına yuvarla
	for i := 0; i < n; i += 4 {
		t0 += a[i] * b[i]
		t1 += a[i+1] * b[i+1]
		t2 += a[i+2] * b[i+2]
		t3 += a[i+3] * b[i+3]
	}
	for i := n; i < len(a); i++ {
		t0 += a[i] * b[i]
	}
	return t0 + t1 + t2 + t3
}

// IcCarpim4SK, IcCarpim4'ün aynısıdır; tek farkı 4'lük alt dilimler almasıdır
// (SK: sınır kontrolü).
// Go her a[i] erişiminde "i dizinin içinde mi?" diye sınır kontrolü (bounds
// check) yapar. Uzunluğu 4 olan bir alt dilimde derleyici bu kontrolleri
// döngü dışına alabilir.
// Kontrolleri görmek için:
//
//	go build -gcflags=-d=ssa/check_bce/debug=1 ./15-test-benchmark
func IcCarpim4SK(a, b []float32) float32 {
	b = b[:len(a)] // derleyiciye "b en az a kadar uzun" bilgisini ver
	var t0, t1, t2, t3 float32
	n := len(a) &^ 3
	for i := 0; i < n; i += 4 {
		aa, bb := a[i:i+4:i+4], b[i:i+4:i+4]
		t0 += aa[0] * bb[0]
		t1 += aa[1] * bb[1]
		t2 += aa[2] * bb[2]
		t3 += aa[3] * bb[3]
	}
	for i := n; i < len(a); i++ {
		t0 += a[i] * b[i]
	}
	return t0 + t1 + t2 + t3
}

// Normallestir, vektörü birim uzunluğa getirir (yerinde değiştirir).
// Sıfır vektör olduğu gibi bırakılır.
func Normallestir(v []float32) {
	var kare float64
	for _, x := range v {
		kare += float64(x) * float64(x)
	}
	if kare == 0 {
		return
	}
	n := float32(math.Sqrt(kare))
	for i := range v {
		v[i] /= n
	}
}
