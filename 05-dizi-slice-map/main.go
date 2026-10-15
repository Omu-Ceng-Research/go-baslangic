// Ders 05 — Dizi, slice, map
//
// Çalıştır: go run ./05-dizi-slice-map
package main

import (
	"fmt"
	"sort"
)

func main() {
	// 1) Dizi: boyutu tipin parçasıdır, sabittir. Nadiren doğrudan kullanılmaktadır.
	var vektor [4]float32 = [4]float32{0.1, 0.2, 0.3, 0.4}
	fmt.Println("dizi:", vektor, "uzunluk:", len(vektor))

	// 2) Slice: dinamik görünüm. Java'daki ArrayList, Ruby'deki Array gibi.
	// Arka planda (işaretçi, uzunluk, kapasite) üçlüsüdür.
	var s []int
	for i := range 10 {
		s = append(s, i)
		fmt.Printf("len=%2d cap=%2d\n", len(s), cap(s)) // kapasite nasıl büyüyor?
	}

	// Önceden yer ayırmak (make) gereksiz kopyalamayı önlemektedir.
	// Bir milyon vektör saklanırken bu önemli olacaktır.
	buyuk := make([]float32, 0, 1_000_000)
	fmt.Println("make ile: len =", len(buyuk), "cap =", cap(buyuk))

	// 3) Dilimleme: kopya DEĞİL, aynı belleğe bakan yeni bir görünüm.
	orijinal := []int{10, 20, 30, 40, 50}
	parca := orijinal[1:3] // [20 30]
	parca[0] = 999
	fmt.Println("orijinal (dikkat!):", orijinal) // 999 orijinali de değiştirdi

	// Gerçek kopya istiyorsanız:
	kopya := make([]int, len(orijinal))
	copy(kopya, orijinal)
	kopya[0] = -1
	fmt.Println("orijinal:", orijinal, "kopya:", kopya)

	// 4) Map: Java'daki HashMap, Ruby'deki Hash.
	// Örnek: hangi mesaj hangi üyelerde tutuluyor?
	konumlar := map[int][]string{
		500:  {"uye3", "uye4"},
		4501: {"uye3", "uye5", "uye6"},
	}
	konumlar[42] = []string{"uye1", "uye2"}

	if uyeler, ok := konumlar[500]; ok {
		fmt.Println("500 numaralı mesajın yerleri:", uyeler)
	}
	if _, ok := konumlar[7]; !ok {
		fmt.Println("7 numaralı mesaj yok")
	}
	delete(konumlar, 42)

	// DİKKAT: map üzerinde gezinme sırası her çalıştırmada RASTGELEDİR.
	// Sıra gerekiyorsa anahtarları sıralayın.
	anahtarlar := make([]int, 0, len(konumlar))
	for k := range konumlar {
		anahtarlar = append(anahtarlar, k)
	}
	sort.Ints(anahtarlar)
	for _, k := range anahtarlar {
		fmt.Println(k, "->", konumlar[k])
	}

	// 5) Yapı slice'ını özel bir kurala göre sıralama: top-k aramada kullanılacaktır.
	type sonuc struct {
		id   int
		skor float32
	}
	sonuclar := []sonuc{{1, 0.3}, {2, 0.9}, {3, 0.5}, {4, 0.7}}
	sort.Slice(sonuclar, func(i, j int) bool { return sonuclar[i].skor > sonuclar[j].skor })
	fmt.Println("en benzer 2:", sonuclar[:2])

	// 6) İki boyutlu slice
	matris := make([][]float32, 3)
	for i := range matris {
		matris[i] = make([]float32, 4)
		for j := range matris[i] {
			matris[i][j] = float32(i*4 + j)
		}
	}
	fmt.Println("matris:", matris)
}

// SORU: 3 milyon float32'yi [][]float32 yerine tek bir []float32 içinde tutmanın
// bellek ve önbellek (cache) açısından avantajı ne olabilir?
