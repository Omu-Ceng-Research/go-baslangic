// Ders 04 — Fonksiyonlar ve hata yönetimi
//
// Çalıştır: go run ./04-fonksiyonlar-hatalar
//
// Go'da exception yoktur. Hatalar sıradan dönüş değerleridir.
// Bu, C'deki "return -1; errno'ya bak" geleneğinin daha düzenli hâlidir.
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Birden fazla değer döndürmek olağandır: (sonuç, hata)
func bol(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("sıfıra bölme")
	}
	return a / b, nil
}

// Paketin dışına açılan hata değeri (Java'daki özel exception sınıfının karşılığı)
var ErrBulunamadi = errors.New("bulunamadı")

func mesajGetir(id string) (string, error) {
	veriler := map[string]string{"1": "merhaba", "2": "dünya"}
	m, ok := veriler[id]
	if !ok {
		// %w ile hata "sarmalanmaktadır". Böylece üst katman kök nedeni görebilir.
		return "", fmt.Errorf("mesaj %s: %w", id, ErrBulunamadi)
	}
	return m, nil
}

// Komut ayrıştırma: geçen seneki SET/GET protokolünün Go hâli.
// "SET 42 merhaba dünya" -> ("SET", "42", "merhaba dünya")
func komutAyristir(satir string) (komut, id, govde string, err error) {
	parcalar := strings.SplitN(strings.TrimSpace(satir), " ", 3)
	if len(parcalar) < 2 {
		return "", "", "", fmt.Errorf("eksik argüman: %q", satir)
	}
	komut = strings.ToUpper(parcalar[0])
	id = parcalar[1]
	if len(parcalar) == 3 {
		govde = parcalar[2]
	}
	return komut, id, govde, nil // isimli dönüş değerleri
}

// Değişken sayıda argüman (Java'daki int... sayilar)
func topla(sayilar ...int) int {
	t := 0
	for _, s := range sayilar {
		t += s
	}
	return t
}

func main() {
	if s, err := bol(10, 2); err == nil {
		fmt.Println("10/2 =", s)
	}
	if _, err := bol(1, 0); err != nil {
		fmt.Println("hata:", err)
	}

	_, err := mesajGetir("99")
	fmt.Println("hata:", err)
	// errors.Is sarmalanmış hatanın içine bakmaktadır.
	fmt.Println("bulunamadı hatası mı?", errors.Is(err, ErrBulunamadi))

	k, id, g, err := komutAyristir("set 42 merhaba dünya")
	fmt.Printf("komut=%s id=%s gövde=%q err=%v\n", k, id, g, err)
	_, _, _, err = komutAyristir("GET")
	fmt.Println("hata:", err)

	fmt.Println("topla:", topla(1, 2, 3, 4))

	// Fonksiyonlar birer değerdir (Ruby'deki lambda, Java'daki Function<>)
	kare := func(x int) int { return x * x }
	fmt.Println("kare(7) =", kare(7))

	// Sistem çağrısı hataları da error olarak gelmektedir.
	// Altında hangi errno'nun olduğuna bakmak mümkündür.
	_, err = os.Open("/olmayan/dosya")
	fmt.Println("os.Open hatası:", err)
	fmt.Println("dosya yok hatası mı?", errors.Is(err, os.ErrNotExist))
}
