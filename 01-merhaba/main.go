// Ders 01 — Merhaba Go
//
// Çalıştır:   go run ./01-merhaba
// Derle:      go build -o merhaba ./01-merhaba && ./merhaba
//
// C bilen için:    #include yerine import, main() yerine package main + func main().
// Java bilen için: sınıf ve JVM bulunmamaktadır. Çıktı tek bir native binary'dir.
// Ruby bilen için: kod yorumlanmamakta, önce derlenmektedir (go run bunu sizin yerinize yapmaktadır).
package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	fmt.Println("Merhaba, Go!")

	// Printf, C'deki printf ile aynı mantıkla çalışmaktadır. %v "ne olursa yaz" anlamına gelmektedir.
	fmt.Printf("Go sürümü: %s, işletim sistemi: %s, mimari: %s\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)

	// Komut satırı argümanları: C'deki argv, Java'daki String[] args, Ruby'deki ARGV.
	// os.Args[0] programın kendisidir.
	fmt.Println("Argümanlar:", os.Args)

	// Process kimliği: getpid(2) sistem çağrısının Go karşılığı.
	fmt.Println("PID:", os.Getpid())
}

// DENEYİN:
//   1) go build -o merhaba ./01-merhaba
//   2) ls -lh merhaba        -> boyutuna bakın
//   3) file merhaba          -> "statically linked" yazıyor mu?
//   4) ldd merhaba           -> Java'daki gibi bir çalışma zamanı (JVM) gerekiyor mu?
//   5) ./merhaba bir iki üç  -> argümanları görün
