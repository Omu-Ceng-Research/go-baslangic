// Ders 14 — Metin tabanlı TCP sunucusu (geçen seneki SET/GET protokolü)
//
// Çalıştır:  go run ./14-tcp-sunucu            (varsayılan port 6666)
// Bağlan:    nc 127.0.0.1 6666
//
//	SET 1 merhaba dünya
//	GET 1
//	STATS
//	QUIT
//
// Shell ile yük testi:
//
//	for i in $(seq 1 1000); do echo "SET $i mesaj-$i"; done | nc -q1 127.0.0.1 6666 | sort | uniq -c
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"sync/atomic"
)

type Sunucu struct {
	mu       sync.RWMutex // çok okuyucu, tek yazıcı
	veriler  map[string]string
	baglanti atomic.Int64
}

func main() {
	port := flag.Int("port", 6666, "dinlenecek TCP portu")
	flag.Parse()

	// net.Listen: socket(2) + bind(2) + listen(2)
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("dinleniyor: %s", ln.Addr())

	s := &Sunucu{veriler: make(map[string]string)}
	for {
		conn, err := ln.Accept() // accept4(2)
		if err != nil {
			log.Print(err)
			continue
		}
		// Her bağlantı için bir goroutine. Java'daki "new Thread" kalıbının
		// karşılığı ama çok daha ucuz. Arka planda Go, epoll(7) kullanır.
		go s.istemciyeHizmet(conn)
	}
}

func (s *Sunucu) istemciyeHizmet(conn net.Conn) {
	defer conn.Close()
	n := s.baglanti.Add(1)
	defer s.baglanti.Add(-1)
	log.Printf("yeni istemci %s (aktif: %d)", conn.RemoteAddr(), n)

	okuyucu := bufio.NewScanner(conn)
	yazici := bufio.NewWriter(conn)
	defer yazici.Flush()

	for okuyucu.Scan() {
		cevap, devam := s.isle(okuyucu.Text())
		yazici.WriteString(cevap + "\n")
		yazici.Flush()
		if !devam {
			return
		}
	}
}

func (s *Sunucu) isle(satir string) (cevap string, devam bool) {
	parcalar := strings.SplitN(strings.TrimSpace(satir), " ", 3)
	switch strings.ToUpper(parcalar[0]) {
	case "SET":
		if len(parcalar) < 3 {
			return "ERROR kullanım: SET <id> <mesaj>", true
		}
		s.mu.Lock()
		s.veriler[parcalar[1]] = parcalar[2]
		s.mu.Unlock()
		return "OK", true
	case "GET":
		if len(parcalar) < 2 {
			return "ERROR kullanım: GET <id>", true
		}
		s.mu.RLock()
		m, ok := s.veriler[parcalar[1]]
		s.mu.RUnlock()
		if !ok {
			return "NOT_FOUND", true
		}
		return m, true
	case "STATS":
		s.mu.RLock()
		defer s.mu.RUnlock()
		return fmt.Sprintf("mesaj=%d aktif_baglanti=%d", len(s.veriler), s.baglanti.Load()), true
	case "QUIT":
		return "BYE", false
	case "":
		return "ERROR boş komut", true
	default:
		return "ERROR bilinmeyen komut: " + parcalar[0], true
	}
}

// ALIŞTIRMALAR:
//  1) Veriler RAM'de; sunucu kapanınca kayboluyor. Ders 08/09'daki bilgilerle
//     her SET'i bir log dosyasına ekleyin, açılışta log'u okuyup map'i doldurun.
//  2) "OK" cevabını fsync'ten önce mi, sonra mı vermelisiniz?
//  3) strace -f -e trace=network,epoll_wait,epoll_ctl ./d14 ile bir istemci
//     bağlanıp koptuğunda hangi çağrıların yapıldığını izleyin.
