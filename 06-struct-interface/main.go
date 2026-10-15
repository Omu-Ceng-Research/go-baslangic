// Ders 06 — struct, metot, interface
//
// Çalıştır: go run ./06-struct-interface
//
// Go'da sınıf ve kalıtım yoktur. Bunların yerine:
//   - struct  : veri (C'deki struct)
//   - metot   : struct'a bağlanan fonksiyon
//   - interface: "bu metotlara sahip olan her şey" (Ruby'deki duck typing'in
//     derleme zamanında kontrol edilen hâli)
package main

import (
	"errors"
	"fmt"
	"strings"
)

// Uye, kümedeki bir düğümü temsil etmektedir.
type Uye struct {
	Host  string
	Port  int
	Canli bool
}

// Değer alıcılı (value receiver) metot: Uye'nin kopyası üzerinde çalışmaktadır.
func (u Uye) Adres() string {
	return fmt.Sprintf("%s:%d", u.Host, u.Port)
}

// İşaretçi alıcılı (pointer receiver) metot: asıl nesneyi değiştirmektedir.
func (u *Uye) Oldur() {
	u.Canli = false
}

// String() metodu olan her tip fmt ile yazdırılırken bu metodu kullanmaktadır
// (Java'daki toString, Ruby'deki to_s).
func (u Uye) String() string {
	durum := "canlı"
	if !u.Canli {
		durum = "ölü"
	}
	return u.Adres() + " (" + durum + ")"
}

// Depo, mesaj saklayabilen her şeyin sözleşmesidir.
// Java'dan farkı şudur: hiçbir tip "implements Depo" yazmamaktadır.
// Metotları taşıyan her tip otomatik olarak Depo sayılmaktadır.
type Depo interface {
	Kaydet(id string, veri []byte) error
	Oku(id string) ([]byte, error)
}

// BellekDepo: Depo'nun RAM'de çalışan gerçeklemesi.
type BellekDepo struct {
	veriler map[string][]byte
}

func YeniBellekDepo() *BellekDepo {
	return &BellekDepo{veriler: make(map[string][]byte)}
}

func (d *BellekDepo) Kaydet(id string, veri []byte) error {
	d.veriler[id] = veri
	return nil
}

func (d *BellekDepo) Oku(id string) ([]byte, error) {
	v, ok := d.veriler[id]
	if !ok {
		return nil, errors.New("yok: " + id)
	}
	return v, nil
}

// BuyukHarfDepo: başka bir Depo'yu sarmalayıp davranışını değiştirmektedir.
// Kalıtım yerine "gömme" (embedding) ve sarmalama kullanılmaktadır.
type BuyukHarfDepo struct {
	Depo // gömülü alan: Oku metodu otomatik olarak buraya devredilmektedir
}

func (b BuyukHarfDepo) Kaydet(id string, veri []byte) error {
	return b.Depo.Kaydet(id, []byte(strings.ToUpper(string(veri))))
}

// Bu fonksiyon hangi depoyla çalıştığını bilmemektedir. Projede diskteki depo,
// ağ üzerindeki uzak depo vb. hepsi aynı arayüzle kullanılabilir.
func dene(d Depo) {
	_ = d.Kaydet("1", []byte("merhaba"))
	v, err := d.Oku("1")
	fmt.Printf("%T -> %s (hata: %v)\n", d, v, err)
}

func main() {
	u := Uye{Host: "127.0.0.1", Port: 5555, Canli: true}
	fmt.Println(u)
	u.Oldur()
	fmt.Println(u)

	// Anonim struct ve struct karşılaştırma
	a := Uye{"127.0.0.1", 5556, true}
	b := Uye{"127.0.0.1", 5556, true}
	fmt.Println("a == b ?", a == b)

	dene(YeniBellekDepo())
	dene(BuyukHarfDepo{Depo: YeniBellekDepo()})

	// Tip sorgulama (Java'daki instanceof)
	var d Depo = YeniBellekDepo()
	if bd, ok := d.(*BellekDepo); ok {
		fmt.Println("bu bir BellekDepo, kayıt sayısı:", len(bd.veriler))
	}
}

// ALIŞTIRMA: Depo arayüzüne uyan bir DosyaDepo yazın (ders 08'den sonra).
// Her mesaj ayrı bir dosyada tutulmalı ve main içinde dene(DosyaDepo{...}) çalışmalıdır.
