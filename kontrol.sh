#!/usr/bin/env bash
# Tüm dersleri derler, testleri çalıştırır ve etkileşimsiz dersleri sırayla koşar.
# Kullanım: ./kontrol.sh
set -euo pipefail
cd "$(dirname "$0")"

if ! command -v go >/dev/null; then
    echo "Go bulunamadı. Kurulum için README.md'ye bakın." >&2
    exit 1
fi
echo "== $(go version)"

echo "== gofmt"
fark=$(gofmt -l .)
if [[ -n "$fark" ]]; then
    echo "Biçimlendirilmemiş dosyalar (gofmt -w . ile düzeltin):"
    echo "$fark"
fi

echo "== go vet + go build"
go vet ./...
go build ./...

echo "== go test"
go test ./...

# 12 (sinyaller) ve 14 (TCP sunucusu) etkileşimlidir, elle çalıştırın.
for ders in 01-merhaba 02-degiskenler-tipler 03-kontrol-donguler \
            04-fonksiyonlar-hatalar 05-dizi-slice-map 06-struct-interface \
            07-pointer-bellek 08-dosya-islemleri 09-sistem-cagrilari \
            10-mmap 11-process-exec 13-goroutine-channel; do
    echo
    echo "################ $ders"
    if ! go run "./$ders"; then
        echo "!! $ders başarısız oldu" >&2
        exit 1
    fi
done

echo
echo "Tamam. 12-sinyaller ve 14-tcp-sunucu derslerini README'deki gibi elle deneyin."
