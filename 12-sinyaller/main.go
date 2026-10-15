//go:build unix

// Ders 12 — Sinyaller ve düzgün kapanış (graceful shutdown)
//
// Çalıştır: go run ./12-sinyaller
// Sonra:
//   - Ctrl+C'ye basın                         (SIGINT)
//   - başka terminalden: kill -TERM <pid>     (SIGTERM)
//   - başka terminalden: kill -USR1 <pid>     (SIGUSR1: istatistik bas)
//   - başka terminalden: kill -9 <pid>        (SIGKILL: yakalanamaz! ne oldu?)
//
// Bir sunucu kapanırken yarım kalan yazmaları bitirmeli, dosyaları kapatmalıdır.
// Ancak SIGKILL ve elektrik kesintisi bu şansı vermemektedir. Projede WAL ve
// fsync kullanılmasının sebebi budur.
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("PID:", os.Getpid())
	fmt.Println("Çalışıyorum. Ctrl+C, kill -TERM, kill -USR1 veya kill -9 deneyin.")

	f, err := os.Create("ders12-gunluk.log")
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)

	// Sinyaller bir channel'a teslim edilmektedir (channel'lar ders 13'te ele alınacaktır).
	sinyaller := make(chan os.Signal, 1)
	signal.Notify(sinyaller, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	tik := time.NewTicker(500 * time.Millisecond)
	defer tik.Stop()
	yazilan := 0

	for {
		select {
		case <-tik.C:
			yazilan++
			fmt.Fprintf(w, "kayıt %d\n", yazilan) // tamponda beklemektedir, henüz diske yazılmamıştır!

		case s := <-sinyaller:
			switch s {
			case syscall.SIGUSR1:
				fmt.Printf("SIGUSR1: şu ana kadar %d kayıt, tamponda %d bayt bekliyor\n",
					yazilan, w.Buffered())
			default:
				fmt.Printf("\n%v alındı, düzgün kapanıyorum...\n", s)
				w.Flush()
				f.Sync()
				f.Close()
				fmt.Printf("%d kayıt diske yazıldı. Hoşça kalın.\n", yazilan)
				return
			}
		}
	}
}

// SORULAR:
//  1) kill -TERM ile kapattıktan sonra ders12-gunluk.log dosyasında kaç satır var?
//  2) kill -9 ile kapattıktan sonra kaç satır var? Neden farklı?
//  3) Ruby'deki Signal.trap("TERM") { ... } ve C'deki sigaction() ile karşılaştırın.
