// Ders 02 — Değişkenler, tipler, sabitler, string'ler
//
// Çalıştır: go run ./02-degiskenler-tipler
package main

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

// Sabitler. iota, C'deki enum gibi 0'dan artan değer üretir.
const (
	Lider = iota // 0
	Uye          // 1
	Olu          // 2
)

func main() {
	// 1) Tanımlama biçimleri
	var a int = 10 // uzun hâl (Java'ya benzer ama tip sonda)
	var b = 20     // tip çıkarımı
	c := 30        // en sık kullanılan kısa hâl (sadece fonksiyon içinde)
	fmt.Println("a, b, c:", a, b, c)

	// 2) Sıfır değerler: Go'da "ilklendirilmemiş" değişken yoktur.
	// C'deki çöp değer ya da Java'daki null hatası yaşanmaz.
	var i int
	var f float64
	var s string
	var ok bool
	fmt.Printf("sıfır değerler -> int:%d float:%g string:%q bool:%t\n", i, f, s, ok)

	// 3) Boyutu belli tipler. Projede vektörleri float32 olarak tutacağız.
	var x int32 = 7
	var y float32 = 0.5
	// Go örtük tip dönüşümü YAPMAZ. Aşağıdaki satır derlenmez:
	//   z := x * y
	z := float32(x) * y
	fmt.Println("float32(x) * y =", z)

	// 4) String <-> sayı dönüşümü (Java'daki Integer.parseInt, Ruby'deki to_i)
	n, err := strconv.Atoi("42")
	fmt.Println("Atoi:", n, err)
	_, err = strconv.Atoi("kırk iki")
	fmt.Println("Hatalı Atoi:", err)

	// 5) String'ler UTF-8 BAYT dizisidir. Türkçe karakterlerde dikkat!
	kelime := "ağaç"
	fmt.Println("len(\"ağaç\") =", len(kelime), "bayt")              // 6
	fmt.Println("karakter sayısı =", utf8.RuneCountInString(kelime)) // 4
	for idx, r := range kelime {
		fmt.Printf("  bayt %d: %c (U+%04X)\n", idx, r, r)
	}

	// 6) Sabitler
	fmt.Println("Lider, Uye, Olu =", Lider, Uye, Olu)
}

// SORU: C'de char* ile "ağaç" yazdırıp strlen() alırsanız kaç çıkar? Neden?
