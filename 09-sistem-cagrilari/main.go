//go:build linux

// Ders 09 — Sistem çağrılarını doğrudan çağırmak
//
// Çalıştır: go run ./09-sistem-cagrilari     (Linux veya WSL2)
//
// Ders 08'de os.Create, f.Write gibi üst seviye fonksiyonlar kullandık.
// Onların altında open(2), write(2), close(2) sistem çağrıları vardır.
// Burada aynı çağrıları "syscall" paketiyle doğrudan yapıyoruz.
// C'de yazacağınız kodun neredeyse birebir aynısı.
//
// Not: syscall paketi standart kütüphanededir. Projede onun güncel hâli olan
// golang.org/x/sys/unix paketini kullanacağız; fonksiyon isimleri aynıdır.
package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"
)

func main() {
	kimlik()
	hamDosyaIslemleri()
	errnoOrnegi()
	fsyncMaliyeti()
}

// [1] Process bilgileri: getpid(2), getppid(2), getuid(2), uname(2)
func kimlik() {
	fmt.Println("[1] pid:", syscall.Getpid(), "ppid:", syscall.Getppid(), "uid:", syscall.Getuid())

	var u syscall.Utsname
	if err := syscall.Uname(&u); err == nil {
		fmt.Println("[1] çekirdek sürümü:", cString(u.Release[:]))
	}
}

// [2] open / write / lseek / read / close — C'deki gibi dosya tanımlayıcı (fd) ile
func hamDosyaIslemleri() {
	yol := "ders09-ham.txt"

	// open(2): bayraklar ve izinler C ile aynıdır.
	fd, err := syscall.Open(yol, syscall.O_CREAT|syscall.O_WRONLY|syscall.O_TRUNC, 0o644)
	if err != nil {
		panic(err)
	}
	fmt.Println("[2] open() dönen fd:", fd) // 0,1,2 stdin/stdout/stderr; sonraki ilk boş fd

	n, err := syscall.Write(fd, []byte("çekirdeğe doğrudan yazıldı\n"))
	fmt.Println("[2] write() yazılan bayt:", n, "hata:", err)
	syscall.Close(fd)

	fd, err = syscall.Open(yol, syscall.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer syscall.Close(fd)

	// lseek(2): 7. bayta atla, oradan oku
	syscall.Seek(fd, 7, 0)
	tampon := make([]byte, 64)
	n, _ = syscall.Read(fd, tampon)
	fmt.Printf("[2] 7. bayttan itibaren read(): %q\n", tampon[:n])

	// fstat(2)
	var st syscall.Stat_t
	syscall.Fstat(fd, &st)
	fmt.Printf("[2] fstat: boyut=%d bayt, inode=%d, blok=%d\n", st.Size, st.Ino, st.Blocks)
}

// [3] Hatalar errno olarak gelir. C'deki perror() çıktısıyla karşılaştırın.
func errnoOrnegi() {
	_, err := syscall.Open("/root/gizli", syscall.O_RDONLY, 0)
	var errno syscall.Errno
	if errors.As(err, &errno) {
		fmt.Printf("[3] errno=%d (%s)\n", int(errno), errno.Error())
	}
	_, err = syscall.Open("/olmayan", syscall.O_RDONLY, 0)
	fmt.Println("[3] ENOENT mi?", err == syscall.ENOENT)
}

// [4] write(2) veriyi diske değil, çekirdeğin sayfa önbelleğine (page cache) yazar.
// Verinin gerçekten diske ulaştığından emin olmak için fsync(2) gerekir.
// Bunun bir bedeli vardır. Ölçelim:
func fsyncMaliyeti() {
	const adet = 200
	kayit := []byte("SET 1 kalıcı olmalı\n")

	olc := func(ad string, herSeferindeFsync bool) {
		fd, err := syscall.Open("ders09-"+ad+".log", syscall.O_CREAT|syscall.O_WRONLY|syscall.O_TRUNC, 0o644)
		if err != nil {
			panic(err)
		}
		defer syscall.Close(fd)

		basla := time.Now()
		for range adet {
			syscall.Write(fd, kayit)
			if herSeferindeFsync {
				syscall.Fsync(fd)
			}
		}
		if !herSeferindeFsync {
			syscall.Fsync(fd) // sadece en sonda bir kez
		}
		fmt.Printf("[4] %-22s %d kayıt: %v\n", ad+":", adet, time.Since(basla))
	}

	olc("fsync-yok-sonda-bir", false)
	olc("her-kayitta-fsync", true)

	fmt.Println("[4] Soru: \"OK\" cevabını hangi noktada vermek dürüst olur?")
	os.Remove("ders09-fsync-yok-sonda-bir.log")
	os.Remove("ders09-her-kayitta-fsync.log")
}

// C tarzı sıfırla biten diziyi Go string'ine çevirir.
// (Generic: Utsname alanları amd64'te int8, bazı ARM'larda uint8'dir.)
func cString[T int8 | uint8](b []T) string {
	s := make([]byte, 0, len(b))
	for _, c := range b {
		if c == 0 {
			break
		}
		s = append(s, byte(c))
	}
	return string(s)
}

// DENEYİN:
//   go build -o d09 ./09-sistem-cagrilari
//   strace -e trace=openat,write,read,lseek,fsync,close ./d09
// Kendi çağırdığınız syscall'ları çıktıda bulun. Go'nun syscall.Open'ı
// aslında hangi çağrıyı yapıyor: open mı, openat mı?
