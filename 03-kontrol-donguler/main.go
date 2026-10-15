// Ders 03 — if, for, switch, defer
//
// Çalıştır: go run ./03-kontrol-donguler
package main

import "fmt"

func main() {
	// Go'da TEK döngü anahtar kelimesi vardır: for. while ve do-while yoktur.

	// 1) C/Java tarzı for
	for i := 0; i < 3; i++ {
		fmt.Println("klasik for, i =", i)
	}

	// 2) while gibi kullanım
	sayac := 3
	for sayac > 0 {
		fmt.Println("while gibi, sayac =", sayac)
		sayac--
	}

	// 3) range: Ruby'deki each, Java'daki for-each
	uyeler := []string{"5555", "5556", "5557"}
	for indeks, port := range uyeler {
		fmt.Printf("üye %d -> port %s\n", indeks, port)
	}

	// 4) Sayı üzerinde range (Go 1.22+): Ruby'deki 3.times
	for i := range 3 {
		fmt.Println("range 3, i =", i)
	}

	// 5) Sonsuz döngü + break/continue. Sunucularda sıkça kullanılmaktadır.
	n := 0
	for {
		n++
		if n%2 == 0 {
			continue
		}
		if n > 7 {
			break
		}
		fmt.Println("tek sayı:", n)
	}

	// 6) if'in içinde kısa tanımlama: değişken sadece if bloğunda yaşamaktadır.
	if kalan := 10 % 3; kalan != 0 {
		fmt.Println("10 % 3 =", kalan)
	}

	// 7) switch: break yazmaya gerek yoktur, case'ler birbirine düşmemektedir.
	komut := "GET"
	switch komut {
	case "SET":
		fmt.Println("yazma komutu")
	case "GET":
		fmt.Println("okuma komutu")
	default:
		fmt.Println("bilinmeyen komut")
	}

	// İfadesiz switch: uzun if-else zincirlerinin yerine geçmektedir.
	gecikmeMs := 120
	switch {
	case gecikmeMs < 50:
		fmt.Println("hızlı")
	case gecikmeMs < 200:
		fmt.Println("kabul edilebilir")
	default:
		fmt.Println("yavaş")
	}

	// 8) defer: fonksiyon bittiğinde LIFO sırasıyla çalışmaktadır.
	// Java'daki finally'nin, C'deki goto cleanup kalıbının karşılığı.
	// Dosya ve soket kapatmada sürekli kullanılacaktır.
	defer fmt.Println("defer 1 (en son çalışır)")
	defer fmt.Println("defer 2")
	fmt.Println("main bitiyor...")
}
