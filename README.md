# Go'ya Hızlı Başlangıç — Sistem Programlama

Bu repo, C, Java ve Ruby bilen öğrencilerin **iki haftada** Go ile dönem ödevine hazır
hâle gelmesi için yazılmış 15 küçük ders programını ve bunların çalıştırılacağı hazır bir
Linux ortamını içermektedir. Her program tek başına çalıştırılabilmektedir.

Derslerin ikinci yarısı (09–15) doğrudan ödevin parçalarına hazırlık niteliğindedir: dosyayı `mmap`
ile okumak, `fsync` ile kalıcılık, ayrı bir process'i (model servisi) yönetmek, sinyaller,
eşzamanlı sorgu ve ölçüm.

> **Önerilen tempo:** 1. hafta 01–08 arası dersler, 2. hafta 09–15 arası.

## İçindekiler

1. [Neden Docker?](#1-neden-docker)
2. [Kurulum](#2-kurulum-bir-kez)
3. [Ortamı başlatmak](#3-ortamı-başlatmak)
4. [Konteynerin içinde sık kullanılan komutlar](#4-konteynerin-içinde-sık-kullanılan-komutlar)
5. [Her derste ne yapacaksınız?](#5-her-derste-ne-yapacaksınız)
6. [Dersler](#6-dersler)
7. [Ödev için modelli ortam](#7-ödev-için-modelli-ortam)
8. [Sorun çıkarsa](#8-sorun-çıkarsa)
9. [Ayrıntılar](#9-ayrıntılar)
10. [Klasörde ne var?](#10-klasörde-ne-var)

## 1. Neden Docker?

Derslerin bir kısmı (09–12) yalnızca Linux'ta bulunan sistem çağrılarını kullanmaktadır.
Sınıfta Windows, macOS ve Linux kullanan öğrenciler bulunmaktadır. Herkesin aynı sonucu
görmesi için tüm çalışmalar Docker içinde **aynı Linux ortamında** yapılacaktır. Go'nun
ya da başka bir aracın ayrıca kurulması gerekmemektedir.

## 2. Kurulum (bir kez)

| Sistem | Kurulum |
|---|---|
| **Windows 10/11** | [Docker Desktop](https://www.docker.com/products/docker-desktop/)'ı kurun. Kurulumda **WSL 2** seçeneğini işaretli bırakın. Kurulum sırasında bilgisayarın yeniden başlatılması gerekebilir. |
| **macOS** (Intel ya da Apple M1–M4) | [Docker Desktop](https://www.docker.com/products/docker-desktop/) ya da daha hafif bir alternatif olan [OrbStack](https://orbstack.dev). |
| **Linux** | `sudo apt install docker.io`, ardından `sudo usermod -aG docker $USER` ve oturumu kapatıp yeniden açın. |

Kontrol için bir terminalde şunu çalıştırın:

```bash
docker run --rm hello-world
```

`Hello from Docker!` yazısı görüntüleniyorsa kurulum tamamlanmıştır.

> [!IMPORTANT]
> Docker Desktop'ta **Settings → Resources** bölümünden Docker'a **en az 4 GB bellek**
> ayırın. Daha az bellekle model programları açılırken sessizce kapanmaktadır.

Sonra bu repoyu indirin: sağ üstteki yeşil **Code** düğmesinden **Download ZIP** ya da

```bash
git clone https://github.com/Omu-Ceng-Research/go-baslangic.git
```

## 3. Ortamı başlatmak

Docker Desktop açıkken, bu klasörün içinde:

| Sistem | Komut |
|---|---|
| Windows | `baslat.cmd` dosyasına **çift tıklayın** |
| macOS / Linux | `./baslat.sh` |

İlk çalıştırmada yaklaşık **900 MB** indirilmektedir; bu işlem birkaç dakika sürebilir.
Sonraki açılışlar birkaç saniye sürmektedir. Ortam açıldığında bir Linux terminali
görüntülenmektedir; komutlar bu terminale yazılmalıdır.

## 4. Konteynerin içinde sık kullanılan komutlar

```bash
ders 01                  # 01. dersi derle ve çalıştır
ders 10 -adet 300000     # derse argüman ver
ders 15                  # testleri ve ölçümleri (benchmark) çalıştır
cd /calisma              # derlenmiş programlar (d01, d10 ...) burada
strace -c ./d08          # programın yaptığı sistem çağrılarını say
man 2 mmap               # bir sistem çağrısının belgesini oku
exit                     # çık
```

- **Kodu kendi bilgisayarınızdaki editörle** (VS Code, GoLand, Notepad++ vb.) bu klasörde
  düzenleyin. Değişiklikler konteynerde anında görünmektedir; `ders NN` ile yeniden
  çalıştırmanız yeterlidir.
- **İkinci terminal:** Bazı derslerde iki terminal gerekmektedir (ör. ders 12'de `kill`, ders 14'te
  `nc`). Başlatma scriptini (`baslat.cmd` / `baslat.sh`) bir kez daha çalıştırın; script açık olan
  ortama ikinci bir terminal açmaktadır.
- `exit` ile çıktığınızda ortam kapanmaktadır. Kodunuz ve `/calisma` klasörü
  **silinmemektedir**; bir sonraki açılışta yerinde durmaktadır.

## 5. Her derste ne yapacaksınız?

1. Dersin `main.go` dosyasını baştan sona okuyun. Açıklamalar kodun içindedir ve C, Java,
   Ruby'deki karşılıkları belirtilmiştir.
2. `ders NN` ile çalıştırın, çıktıyı kodla karşılaştırın.
3. Dosyanın sonundaki **DENEYİN / SORULAR / ALIŞTIRMA** bölümlerini yapın.

Ölçüm yapan derslerde (08, 10, 15) sonuçlar makineden makineye ve çalıştırmadan
çalıştırmaya değişmektedir. Tek ölçüme güvenmeyin; birkaç kez tekrarlayın
(ör. `go test -bench . -count 5`).

## 6. Dersler

| # | Klasör | Konu | C / Java / Ruby'deki karşılığı | Ödevde nerede? |
|---|---|---|---|---|
| 01 | [`01-merhaba`](01-merhaba) | `package main`, `fmt`, `go build`, statik binary | `printf`, `main`, JVM'siz çalışma | — |
| 02 | [`02-degiskenler-tipler`](02-degiskenler-tipler) | `:=`, sıfır değerler, örtük dönüşüm yok, UTF-8 string | `int32_t`, `float`; `String` | Türkçe metinler |
| 03 | [`03-kontrol-donguler`](03-kontrol-donguler) | tek döngü `for`, `range`, `switch`, `defer` | `while`, `each`, `finally` | Her yerde |
| 04 | [`04-fonksiyonlar-hatalar`](04-fonksiyonlar-hatalar) | çoklu dönüş, `error`, `%w`, `errors.Is` | `errno`, exception, `raise` | Komut ayrıştırma |
| 05 | [`05-dizi-slice-map`](05-dizi-slice-map) | slice, `append`, `cap`, map, sıralama | `ArrayList`, `HashMap`, `Hash` | Top-k sonuçlar |
| 06 | [`06-struct-interface`](06-struct-interface) | struct, metot, örtük interface, gömme | `struct`, `implements`, duck typing | Depo arayüzü |
| 07 | [`07-pointer-bellek`](07-pointer-bellek) | işaretçi, aritmetik yok, GC, `sizeof` | `*p`, `malloc/free` | Bellek hesabı |
| 08 | [`08-dosya-islemleri`](08-dosya-islemleri) | `os`, `bufio`, buffered/unbuffered, ikili dosya, `Seek` | `fopen`, `BufferedWriter`, `File.open` | Vektör dosyası |
| 09 | [`09-sistem-cagrilari`](09-sistem-cagrilari) | `open/write/read/lseek/fstat/fsync`, errno | C'deki `unistd.h` birebir | **Kalıcılık** |
| 10 | [`10-mmap`](10-mmap) | `mmap`, `madvise`, `getrusage`, page fault | `mmap(2)` | **Vektör deposu** |
| 11 | [`11-process-exec`](11-process-exec) | `exec`, pipe, zaman aşımı, `SIGKILL`, gözetmen | `fork/exec/pipe/waitpid`, `Process.spawn` | **Model process'i** |
| 12 | [`12-sinyaller`](12-sinyaller) | `os/signal`, düzgün kapanış, SIGKILL'in farkı | `sigaction`, `Signal.trap` | Düğüm kapanışı |
| 13 | [`13-goroutine-channel`](13-goroutine-channel) | goroutine, channel, `Mutex`, `select`, `-race` | `pthread`, `Thread`, `Queue` | **Eşzamanlı sorgu** |
| 14 | [`14-tcp-sunucu`](14-tcp-sunucu) | `net.Listen`, bağlantı başına goroutine, SET/GET | `socket/bind/listen/accept` | İstemci arayüzü |
| 15 | [`15-test-benchmark`](15-test-benchmark) | `go test`, tablo testleri, benchmark | JUnit, RSpec | **Ölçüm tabloları** |

## 7. Ödev için modelli ortam

Ödevde metinleri vektöre çeviren küçük bir yapay zekâ modeli kullanılacaktır. Koordinatör
bu modeli kendi çocuk process'i olarak başlatmaktadır; bu yüzden model koordinatörle **aynı
imajda** olmalıdır. Model yalnızca **koordinatörü çalıştıracak makinede** gereklidir;
grupta bir kişinin bilgisayarı yeterlidir.

| Komut | İmaj | İçinde | Boyut | Kim kullanır? |
|---|---|---|---|---|
| `./baslat.sh` | `:latest` | Go, strace, gcc, man sayfaları | ~0,9 GB | Herkes; dersler ve depolama düğümleri |
| `./baslat.sh model` | `:model` | + `model-embedder` (multilingual-e5-small) | ~1,6 GB | Koordinatörü çalıştıran |
| `./baslat.sh llm` | `:llm` | + `model-uretici` (Qwen3.5-0.8B) | ~2,5 GB | Yalnızca cevap üretme bonusunu yapan grup |

Windows'ta komut satırından: `baslat.cmd model` ya da `baslat.cmd llm`. Katmanlar ortak
olduğundan büyük imajı indiren kullanıcı küçük imajları da indirmiş sayılmaktadır.

İçine girince kurulumu doğrulayın:

```bash
model-kontrol
```

En sonda `TAMAM` yazıyorsa bilgisayarınız modeli çalıştırmaya uygundur. **İlk hafta
herkesin bu kontrolü bir kez yapması önerilmektedir;** grupların koordinatör makinesi
buna göre seçilmelidir.

| Model | Görev | Bellek | Hız (1 cümle) | Lisans |
|---|---|---|---|---|
| [multilingual-e5-small](https://huggingface.co/intfloat/multilingual-e5-small) | Metni 384 boyutlu vektöre çevirir | ~1,2 GB | ~30–100 ms | MIT |
| [Qwen3.5-0.8B](https://huggingface.co/Qwen/Qwen3.5-0.8B) (Q8_0) | Bağlama göre cevap üretir (bonus) | ~1,0 GB | birkaç saniye | Apache 2.0 |

## 8. Sorun çıkarsa

| Belirti | Çözüm |
|---|---|
| "Docker kurulu değil" | 2. adımdaki kurulumu yapın. Kurulumdan sonra terminali kapatıp yeniden açın. |
| "Docker çalışmıyor" | Docker Desktop uygulamasını açın, motorun çalıştığını gösteren durumu bekleyin, sonra tekrar deneyin. |
| Windows'ta "WSL 2 installation is incomplete" ya da benzeri | Yönetici olarak açılmış PowerShell'de `wsl --update`, sonra bilgisayarı yeniden başlatın. |
| Windows'ta "Virtualization is not enabled" | BIOS/UEFI ayarlarında **Intel VT-x** ya da **AMD-V / SVM** seçeneğini açın. Bilgisayarınızın markasına göre arayın. |
| İndirme çok yavaş ya da yarıda kesiliyor | Scripti tekrar çalıştırın; tamamlanmış parçalar yeniden indirilmemektedir. Okul ağında herkes aynı anda indirdiğinde hız düşmektedir; indirmenin evde yapılması önerilmektedir. |
| `model-kontrol` çalışırken `Killed` yazıp kapanıyor | Docker'a ayrılan bellek yetersizdir. **Settings → Resources → Memory** değerini en az 4 GB yapın, tarayıcı gibi büyük programları kapatıp tekrar deneyin. |
| Ders 09–12 Windows'ta `go run` ile derlenmiyor | Normaldir; bu dersler yalnızca Linux içindir. Docker ortamının içinde `ders 09` ile çalıştırın. |
| "Ders bulunamadı" | Ortamı bu klasörün içinden başlattığınızdan emin olun. |

Bunlarla çözülmeyen bir sorunda **hata mesajının tamamını** (ekran görüntüsü değil, metin
olarak), işletim sisteminizi ve hangi komutu çalıştırdığınızı yazarak dersin iletişim
kanalından sorun.

## 9. Ayrıntılar

### VS Code ile (isteğe bağlı)

"Dev Containers" eklentisini kurun, bu klasörü açın ve sağ alttaki **Reopen in Container**
düğmesine basın. Editör, terminal ve Go eklentisi doğrudan konteynerin içinde çalışmaktadır.

### Docker yerine

- **Windows, WSL2 + Ubuntu:** `wsl --install` ile Ubuntu kurup Go'yu oraya kurabilirsiniz.
  Dosyaları `C:\` altında değil WSL'in kendi diskinde (`~/...`) tutun; aksi hâlde ölçüm
  dersleri yavaş ve yanıltıcı sonuç vermektedir.
- **Linux:** Go'yu <https://go.dev/dl> adresinden kurup `go run ./01-merhaba` gibi doğrudan
  çalıştırabilirsiniz. `strace` paketini de kurun. Tüm dersleri derleyip etkileşimsiz
  olanları sırayla çalıştırmak için: `./kontrol.sh`

### Neden kod ve çıktılar ayrı yerde?

Mac ve Windows'ta Docker'a bağlanan klasörler sanal bir paylaşım katmanından geçmektedir
ve yavaştır. Ders 08–10'daki ölçümler orada yapılırsa çekirdeğin page cache'i yerine bu
katman ölçülmüş olur. `ders` komutu bu yüzden programları `/calisma`'da (gerçek Linux
dosya sistemi) çalıştırmaktadır.

### strace ile gözlem

Sistem çağrısı derslerinde programı önce derleyin, sonra `strace` altında çalıştırın.
`go run` ile çalıştırıldığında derleyicinin çağrıları da görünmektedir. Docker'da
`ders NN` derlemeyi sizin yerinize yapmaktadır; binary `/calisma/dNN` olarak
saklanmaktadır.

```bash
go build -o d10 ./10-mmap
strace -c ./d10                                   # özet: hangi çağrı kaç kez
strace -f -e trace=mmap,madvise,read ./d10        # sadece belirli çağrılar
strace -f -e trace=clone,clone3,execve,wait4 ./d11
```

Go runtime birden fazla thread kullandığı için `-f` bayrağı çoğu zaman gereklidir.
Çıktıdaki `futex` ve `rt_sigaction` gürültüsünü `-e trace=...` ile filtreleyin.

### Hızlı karşılaştırma

| Konu | C | Java | Ruby | Go |
|---|---|---|---|---|
| Değişken | `int x = 5;` | `int x = 5;` | `x = 5` | `x := 5` |
| Döngü | `for (;;)`, `while` | `for`, `while` | `each`, `times` | sadece `for` |
| Dinamik dizi | elle `realloc` | `ArrayList<T>` | `Array` | `[]T` + `append` |
| Sözlük | — | `HashMap<K,V>` | `Hash` | `map[K]V` |
| Hata | dönüş değeri + `errno` | exception | exception | `error` dönüş değeri |
| Temizlik | `goto cleanup` | `finally` | `ensure` | `defer` |
| Arayüz | fonksiyon işaretçisi | `interface` + `implements` | duck typing | örtük `interface` |
| Eşzamanlılık | `pthread` | `Thread`, `Executor` | `Thread` (GVL) | `go f()` + channel |
| Bellek | `malloc/free` | GC | GC | GC (+ `mmap` ile GC dışı) |

### Bilinen sınırlar ve notlar

- Ders 10'daki soğuk önbellek deneyi (`drop_caches`) konteynerde çalışmamaktadır; çekirdek ayarı
  konteynerin dışındadır. WSL2 ya da gerçek Linux'ta `sudo` ile yapılabilir.
- Programlar çalıştıkları klasöre `ders08-cikti/`, `ders10-vektorler.f32` gibi dosyalar
  üretmektedir. Bu dosyalar `.gitignore` içinde yer almaktadır ve silinebilir.
- `nc -q1` bayrağı Debian/Ubuntu'daki `netcat-openbsd` içindir. Başka bir `nc` sürümünde
  `-q1` yerine `-N` gerekebilir.
- İmajın yeni sürümünü almak için:
  `docker pull ghcr.io/omu-ceng-research/sysprog-go:latest`

## 10. Klasörde ne var?

| Yol | Açıklama |
|---|---|
| `01-merhaba` … `15-test-benchmark` | Ders programları |
| `baslat.cmd`, `baslat.ps1`, `baslat.sh` | Ortamı başlatan scriptler |
| `docker/` | Ortamın tarifi (dokunmanız gerekmez) |
| `.devcontainer/` | VS Code "Dev Containers" ayarı |
| `kontrol.sh` | Bütün dersleri sırayla derleyip denetlemektedir |

Bu klasördeki dosyaları silmeyin ya da yeniden adlandırmayın; başlatma scriptleri
onları bu adlarla aramaktadır.

## Lisans

[MIT](LICENSE)
