// Ders 13 — Goroutine, channel, WaitGroup, Mutex, select
//
// Çalıştır:            go run ./13-goroutine-channel
// Yarış denetleyicisi: go run -race ./13-goroutine-channel
//
// Goroutine, Java'daki Thread'in çok hafif hâlidir (birkaç KB yığın).
// Binlercesini açmak sorun değildir. Go runtime bunları az sayıda işletim
// sistemi thread'i üzerine dağıtmaktadır.
package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

func main() {
	temelGoroutine()
	mutexOrnegi()
	scatterGather()
	zamanAsimliCagri()
}

// [1] go anahtar kelimesi + WaitGroup ile hepsinin bitmesini beklemek.
func temelGoroutine() {
	var wg sync.WaitGroup
	for i := range 3 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("[1] goroutine", i, "çalıştı")
		}()
	}
	wg.Wait()
}

// [2] Paylaşılan veri + Mutex. Kilidi kaldırıp -race ile çalıştırın!
func mutexOrnegi() {
	var (
		mu     sync.Mutex
		sayac  = map[string]int{}
		wg     sync.WaitGroup
		uyeler = []string{"uye1", "uye2", "uye3"}
	)
	for i := range 3000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			sayac[uyeler[i%3]]++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println("[2] üye başına mesaj sayısı:", sayac)
}

// [3] Scatter-gather: projedeki sorgu akışının iskeleti.
// Lider sorguyu tüm shard'lara AYNI ANDA göndermekte, her shard kendi en iyi
// k sonucunu döndürmekte, lider de bunları birleştirip genel en iyi k sonucu bulmaktadır.
type Sonuc struct {
	Shard int
	ID    int
	Skor  float64
}

func shardAra(shard, k int) []Sonuc {
	time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond) // ağ + arama gecikmesi
	sonuclar := make([]Sonuc, k)
	for i := range sonuclar {
		sonuclar[i] = Sonuc{Shard: shard, ID: shard*1000 + i, Skor: rand.Float64()}
	}
	return sonuclar
}

func scatterGather() {
	const shardSayisi, k = 4, 3
	kanal := make(chan []Sonuc, shardSayisi) // tamponlu channel

	basla := time.Now()
	for s := range shardSayisi {
		go func() { kanal <- shardAra(s, k) }() // scatter
	}

	var hepsi []Sonuc
	for range shardSayisi {
		hepsi = append(hepsi, <-kanal...) // gather
	}
	sort.Slice(hepsi, func(i, j int) bool { return hepsi[i].Skor > hepsi[j].Skor })

	fmt.Printf("[3] %d shard paralel sorgulandı, süre %v (sırayla olsaydı ~%d ms)\n",
		shardSayisi, time.Since(basla).Round(time.Millisecond), shardSayisi*100)
	for _, s := range hepsi[:k] {
		fmt.Printf("[3]   shard=%d id=%d skor=%.3f\n", s.Shard, s.ID, s.Skor)
	}
}

// [4] select + time.After: yanıt vermeyen üyeyi beklemeye bir sınır koymak.
func zamanAsimliCagri() {
	cevap := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second) // çökmüş ya da çok yavaş bir üye
		cevap <- "geç kalan cevap"
	}()

	select {
	case c := <-cevap:
		fmt.Println("[4]", c)
	case <-time.After(300 * time.Millisecond):
		fmt.Println("[4] zaman aşımı! Bu üyeyi atla, yedek kopyadan oku.")
	}
}

// ALIŞTIRMA: scatterGather'da shard'lardan biri hiç cevap vermezse ne olur?
// zamanAsimliCagri'daki select kalıbıyla bunu düzeltin.
