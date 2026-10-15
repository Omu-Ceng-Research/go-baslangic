//go:build linux

// Ders 11 — Process başlatma, pipe ile konuşma, zaman aşımı, çöken çocuğu yeniden başlatma
//
// Çalıştır: go run ./11-process-exec
//
// Projede embedding modeli AYRI bir process olarak çalışacak. Düğümünüz onu
// başlatacak, onunla pipe üzerinden konuşacak, çökerse yeniden başlatacak.
// Bu ders, bunun küçük bir provası.
//
// Program kendini "çocuk" modunda tekrar çalıştırır; harici bir araç gerekmez.
package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func main() {
	// Çocuk modları: ebeveyn bizi os.Args[1] ile bu modlardan birinde başlatır.
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "cocuk-servis":
			cocukServis()
			return
		case "cocuk-uyu":
			time.Sleep(time.Hour)
			return
		}
	}

	basitKomut()
	cikisKodu()
	pipeIleKonus()
	zamanAsimi()
	gozetmen()
}

// [1] En basit hâl: komutu çalıştır, çıktısını al. Altında fork/clone + execve vardır.
func basitKomut() {
	cikti, err := exec.Command("uname", "-sr").Output()
	fmt.Printf("[1] uname çıktısı: %s (hata: %v)\n", strings.TrimSpace(string(cikti)), err)
}

// [2] Çıkış kodu: shell'deki $? değerinin karşılığı.
func cikisKodu() {
	err := exec.Command("sh", "-c", "exit 3").Run()
	var cikis *exec.ExitError
	if errors.As(err, &cikis) {
		fmt.Println("[2] çocuğun çıkış kodu:", cikis.ExitCode())
	}
}

// [3] Uzun yaşayan çocukla stdin/stdout pipe'ları üzerinden satır satır konuşmak.
func pipeIleKonus() {
	kendim, _ := os.Executable()
	cmd := exec.Command(kendim, "cocuk-servis")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGKILL, // ebeveyn ölürse çekirdek çocuğu da öldürsün (prctl)
	}
	stdin, _ := cmd.StdinPipe()   // pipe(2) + dup2(2)
	stdout, _ := cmd.StdoutPipe() // pipe(2) + dup2(2)
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	fmt.Println("[3] çocuk başladı, pid:", cmd.Process.Pid)

	okuyucu := bufio.NewReader(stdout)
	for _, metin := range []string{"merhaba", "sistem programlama", "ğüşiöç"} {
		fmt.Fprintln(stdin, metin)
		cevap, _ := okuyucu.ReadString('\n')
		fmt.Printf("[3] %q -> %s", metin, cevap)
	}

	stdin.Close()     // çocuk EOF görür ve kendiliğinden çıkar
	err := cmd.Wait() // wait4(2): zombi kalmasın
	fmt.Println("[3] çocuk çıktı, hata:", err)
}

// [4] Zaman aşımı: yanıt vermeyen çocuğu öldürmek.
func zamanAsimi() {
	kendim, _ := os.Executable()
	ctx, iptal := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer iptal()

	basla := time.Now()
	err := exec.CommandContext(ctx, kendim, "cocuk-uyu").Run() // süre dolunca SIGKILL
	fmt.Printf("[4] %v sonra: %v\n", time.Since(basla).Round(time.Millisecond), err)

	var cikis *exec.ExitError
	if errors.As(err, &cikis) {
		durum := cikis.Sys().(syscall.WaitStatus)
		fmt.Println("[4] sinyalle mi öldü?", durum.Signaled(), "hangi sinyal:", durum.Signal())
	}
}

// [5] Gözetmen (supervisor): çocuk çökerse yeniden başlat.
// Çocuğu her seferinde bilerek öldürüyoruz.
func gozetmen() {
	kendim, _ := os.Executable()
	for deneme := 1; deneme <= 3; deneme++ {
		cmd := exec.Command(kendim, "cocuk-servis")
		stdin, _ := cmd.StdinPipe()
		stdout, _ := cmd.StdoutPipe()
		cmd.Start()

		fmt.Fprintln(stdin, "deneme")
		cevap, _ := bufio.NewReader(stdout).ReadString('\n')
		fmt.Printf("[5] deneme %d, pid %d, cevap: %s", deneme, cmd.Process.Pid, cevap)

		cmd.Process.Signal(syscall.SIGKILL) // kill -9 <pid>
		err := cmd.Wait()
		fmt.Printf("[5] çocuk öldü (%v), yeniden başlatılıyor...\n", err)
	}
}

// cocukServis: stdin'den satır okur, sahte bir "embedding" döndürür.
// Gerçek projede burada küçük bir dil modeli olacak.
func cocukServis() {
	okuyucu := bufio.NewReader(os.Stdin)
	for {
		satir, err := okuyucu.ReadString('\n')
		if err == io.EOF {
			return
		}
		satir = strings.TrimSpace(satir)
		// Sahte embedding: uzunluk, harf sayısı, boşluk sayısı
		fmt.Printf("[%d, %d, %d] (pid %d)\n",
			len(satir), len([]rune(satir)), strings.Count(satir, " "), os.Getpid())
	}
}

// DENEYİN:
//   go build -o d11 ./11-process-exec
//   strace -f -e trace=clone,clone3,execve,pipe2,wait4,kill ./d11 2>&1 | less
// -f bayrağı neden gerekli? Go'nun process başlatırken fork mu, clone mu
// kullandığını çıktıdan bulun.
